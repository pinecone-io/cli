package presenters

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

func ColorizeBool(b bool) string {
	if b {
		return style.StatusGreen("true")
	}
	return style.StatusRed("false")
}

const nonePlaceholder = "<none>"

// DisplayOrNone renders pointers/interfaces as "<none>" when nil, keeps empty strings,
// and renders nil maps/slices as "{}" or "[]".
func DisplayOrNone(val any) string {
	if val == nil {
		return nonePlaceholder
	}

	v := reflect.ValueOf(val)
	for v.IsValid() {
		switch v.Kind() {
		case reflect.Ptr, reflect.Interface:
			if v.IsNil() {
				return nonePlaceholder
			}
			v = v.Elem()
		case reflect.String:
			return v.String()
		case reflect.Map:
			if v.IsNil() {
				return "{}"
			}
			return marshalJSONOrDefault(v.Interface(), "{}")
		case reflect.Slice:
			if v.IsNil() {
				return "[]"
			}
			return marshalJSONOrDefault(v.Interface(), "[]")
		default:
			return fmt.Sprint(v.Interface())
		}
	}

	return nonePlaceholder
}

func marshalJSONOrDefault(val any, def string) string {
	data, err := json.Marshal(val)
	if err != nil {
		return def
	}
	return string(data)
}
