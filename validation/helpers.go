package validation

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func interfaceValue(v reflect.Value) any {
	if v.IsValid() && v.CanInterface() {
		return v.Interface()
	}
	return nil
}

func splitRules(tag string) []string {
	var rules []string
	var cur strings.Builder
	var inQuote byte
	var braces, brackets, parens int

	for i := 0; i < len(tag); i++ {
		c := tag[i]
		switch {
		case inQuote != 0:
			if c == inQuote {
				inQuote = 0
			}
			cur.WriteByte(c)
		case c == '"' || c == '\'':
			inQuote = c
			cur.WriteByte(c)
		case c == '{':
			braces++
			cur.WriteByte(c)
		case c == '}':
			if braces > 0 {
				braces--
			}
			cur.WriteByte(c)
		case c == '[':
			brackets++
			cur.WriteByte(c)
		case c == ']':
			if brackets > 0 {
				brackets--
			}
			cur.WriteByte(c)
		case c == '(':
			parens++
			cur.WriteByte(c)
		case c == ')':
			if parens > 0 {
				parens--
			}
			cur.WriteByte(c)
		case c == ',' && braces == 0 && brackets == 0 && parens == 0:
			rules = append(rules, cur.String())
			cur.Reset()
		default:
			cur.WriteByte(c)
		}
	}
	if cur.Len() > 0 {
		rules = append(rules, cur.String())
	}
	return rules
}

func hasRule(tag, want string) bool {
	for _, r := range splitRules(tag) {
		name, _ := parseRule(r)
		if name == want {
			return true
		}
	}
	return false
}

func parseRule(rule string) (name, arg string) {
	rule = strings.TrimSpace(rule)
	if i := strings.IndexAny(rule, "=:"); i >= 0 {
		name = strings.TrimSpace(rule[:i])
		arg = strings.TrimSpace(rule[i+1:])
		if (strings.HasPrefix(arg, "'") && strings.HasSuffix(arg, "'")) ||
			(strings.HasPrefix(arg, "\"") && strings.HasSuffix(arg, "\"")) {
			if len(arg) >= 2 {
				arg = arg[1 : len(arg)-1]
			}
		}
		return name, arg
	}
	return rule, ""
}

func effective(v reflect.Value) (reflect.Value, bool) {
	nilPtr := false
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			nilPtr = true
			break
		}
		v = v.Elem()
	}
	return v, nilPtr
}

func valueKind(v reflect.Value) reflect.Kind {
	t := v.Type()
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind()
}

func isZero(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Map:
		return v.Len() == 0
	}
	return v.IsZero()
}

func sortMapKeys(keys []reflect.Value) {
	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprint(keys[i].Interface()) < fmt.Sprint(keys[j].Interface())
	})
}
