package presenters

import (
	"reflect"

	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

func ColorizeBool(b bool) string {
	if b {
		return style.StatusGreen("true")
	}
	return style.StatusRed("false")
}

func DisplayOrNone(val any) any {
	if val == nil {
		return "<none>"
	}

	v := reflect.ValueOf(val)
	for v.IsValid() {
		switch v.Kind() {
		case reflect.Ptr, reflect.Interface:
			if v.IsNil() {
				return "<none>"
			}
			v = v.Elem()
		default:
			return v.Interface()
		}
	}

	return "<none>"
}
