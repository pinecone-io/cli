package presenters

import (
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

func DisplayOrNone(val interface{}) string {
	v := reflect.ValueOf(val)

	if !v.IsValid() || (v.Kind() == reflect.Ptr && v.IsNil()) {
		return "<none>"
	}

	if v.Kind() == reflect.Ptr {
		return fmt.Sprintf("%v", v.Elem().Interface())
	}

	return fmt.Sprintf("%v", val)
}
