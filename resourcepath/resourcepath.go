package resourcepath

import (
	"fmt"
	"reflect"
	"strings"
)

func UnmarshalString(path string, s any) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	st := reflect.TypeOf(v)
	pattern := ""
OUTER:
	for i := range st.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("path")
		if tag == "" {
			continue
		}
		ps := strings.Split(tag, ",")
		for _, p := range ps {
			if Match(p, path) {
				pattern = p
				break OUTER
			}
		}
	}
	if pattern == "" {
		return fmt.Errorf("no matching patterns for path: %s", path)
	}
	if err := Sscan(path, pattern, v); err != nil {
		return err
	}
	return nil
}
