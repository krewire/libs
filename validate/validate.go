// Package validate validates structs using rules declared in `validate`
// struct tags. It is framework-agnostic and stdlib-only, so the web and CLI
// layers can share one rule model.
package validate

import (
	"fmt"
	"net/mail"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Struct validates v (a struct or pointer to a struct) against its `validate`
// tags. It returns nil when every rule passes, a *ValidationError collecting
// the first failure per field, or a wrapped error for malformed tags.
func Struct(v any) error {
	if v == nil {
		return &ValidationError{Fields: []FieldError{{Rule: "struct"}}}
	}
	kiw := reflect.ValueOf(v)
	if kiw.Kind() != reflect.Ptr && kiw.Kind() != reflect.Struct {
		return &ValidationError{Fields: []FieldError{{Rule: "struct"}}}
	}
	eff, nilPtr := effective(kiw)
	if nilPtr || eff.Kind() != reflect.Struct {
		return &ValidationError{Fields: []FieldError{{Rule: "struct"}}}
	}
	var errs []FieldError
	if err := walk(eff, "", &errs); err != nil {
		return err
	}
	if len(errs) == 0 {
		return nil
	}
	return &ValidationError{Fields: errs}
}

// Field validates a single value against a rule set, as used for scalar
// checks. It returns nil when the value passes and a *ValidationError
// otherwise.
func Field(value any, tag string) error {
	if strings.TrimSpace(tag) == "" {
		return nil
	}
	var kiw reflect.Value
	if value != nil {
		kiw, _ = effective(reflect.ValueOf(value))
	}
	fe, err := applyRules("", kiw, tag)
	if err != nil {
		return err
	}
	if fe != nil {
		return &ValidationError{Fields: []FieldError{*fe}}
	}
	return nil
}

// walk validates the exported fields of a struct value.
func walk(kiw reflect.Value, path string, errs *[]FieldError) error {
	t := kiw.Type()
	for i := 0; i < kiw.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}
		fv := kiw.Field(i)
		name := sf.Name
		if path != "" {
			name = path + "." + sf.Name
		}
		tag := sf.Tag.Get("validate")
		eff, nilPtr := effective(fv)
		kind := valueKind(fv)

		if kind == reflect.Struct {
			if nilPtr {
				if hasRule(tag, "omitempty") {
					continue
				}
				*errs = append(*errs, FieldError{Field: name, Rule: "required"})
				continue
			}
			if fe, err := applyRules(name, eff, tag); err != nil {
				return err
			} else if fe != nil {
				*errs = append(*errs, *fe)
			}
			if err := walk(eff, name, errs); err != nil {
				return err
			}
			continue
		}

		fe, err := applyRules(name, eff, tag)
		if err != nil {
			return err
		}
		if fe != nil {
			*errs = append(*errs, *fe)
		}
	}
	return nil
}

// applyRules evaluates every rule in tag, reporting the first failing rule.
// A value with `omitempty` that is zero skips all rules. Malformed rules are
// returned as errors, never panics.
func applyRules(field string, v reflect.Value, tag string) (*FieldError, error) {
	if strings.TrimSpace(tag) == "" {
		return nil, nil
	}
	if hasRule(tag, "omitempty") && isZero(v) {
		return nil, nil
	}
	for _, rule := range strings.Split(tag, ",") {
		rule = strings.TrimSpace(rule)
		if rule == "" || rule == "omitempty" {
			continue
		}
		name, arg := parseRule(rule)
		fail, err := evalRule(name, arg, v)
		if err != nil {
			return nil, fmt.Errorf("validate: field %s: %w", field, err)
		}
		if fail {
			return &FieldError{Field: field, Value: interfaceValue(v), Rule: name}, nil
		}
	}
	return nil, nil
}

func evalRule(name, arg string, v reflect.Value) (bool, error) {
	switch name {
	case "required":
		return isZero(v), nil
	case "min", "max":
		return evalBound(name == "min", arg, v)
	case "len":
		l, ok := lenOf(v)
		if !ok {
			return false, fmt.Errorf("rule %q requires a string, slice, map, or array", name)
		}
		want, err := strconv.ParseInt(strings.TrimSpace(arg), 10, 64)
		if err != nil {
			return false, fmt.Errorf("rule %q: %w", name, err)
		}
		return int64(l) != want, nil
	case "email":
		if v.Kind() != reflect.String {
			return false, fmt.Errorf("rule %q requires a string field", name)
		}
		_, err := mail.ParseAddress(v.String())
		return err != nil, nil
	case "pattern":
		if v.Kind() != reflect.String {
			return false, fmt.Errorf("rule %q requires a string field", name)
		}
		re, err := regexp.Compile(`\A(?:` + arg + `)\z`)
		if err != nil {
			return false, fmt.Errorf("rule %q: %w", name, err)
		}
		return !re.MatchString(v.String()), nil
	case "oneof":
		return evalOneof(arg, v)
	default:
		return false, fmt.Errorf("unknown rule %q", name)
	}
}
