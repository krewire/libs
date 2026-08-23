package validate

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func evalBound(isMin bool, arg string, v reflect.Value) (bool, error) {
	want, err := strconv.ParseFloat(strings.TrimSpace(arg), 64)
	if err != nil {
		return false, fmt.Errorf("invalid bound %q", arg)
	}
	var got float64
	switch v.Kind() {
	case reflect.String:
		got = float64(len(v.String()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		got = float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		got = float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		got = v.Float()
	default:
		return false, fmt.Errorf("min/max require a number or string, got %s", v.Kind())
	}
	if isMin {
		return got < want, nil
	}
	return got > want, nil
}
