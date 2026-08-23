package validate

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func evalOneof(arg string, v reflect.Value) (bool, error) {
	items := strings.Fields(arg)
	var s string
	switch v.Kind() {
	case reflect.String:
		s = v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s = strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s = strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		s = strconv.FormatFloat(v.Float(), 'g', -1, 64)
	default:
		return false, fmt.Errorf("oneof requires a string or comparable number, got %s", v.Kind())
	}
	for _, it := range items {
		if s == it {
			return false, nil
		}
	}
	return true, nil
}
