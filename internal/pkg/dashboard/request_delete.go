package dashboard

import (
	"net/http"
)

func DeleteAndDecode[T any](path string) (*T, error) {
	return RequestWithoutBodyAndDecode[T](path, http.MethodDelete)
}
