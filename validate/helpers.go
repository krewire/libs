package validate

import (
	"reflect"
	"strings"
)

func lenOf(v reflect.Value) (int, bool) {
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return v.Len(), true
	}
	return 0, false
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

func interfaceValue(v reflect.Value) any {
	if v.IsValid() && v.CanInterface() {
		return v.Interface()
	}
	return nil
}

func hasRule(tag, want string) bool {
	for _, r := range strings.Split(tag, ",") {
		if strings.TrimSpace(r) == want {
			return true
		}
	}
	return false
}

func parseRule(rule string) (name, arg string) {
	if i := strings.IndexByte(rule, '='); i >= 0 {
		return rule[:i], rule[i+1:]
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
