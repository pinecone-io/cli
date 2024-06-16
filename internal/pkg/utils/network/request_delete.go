package network

import (
	"net/http"
)

func DeleteAndDecode[T any](baseUrl string, path string) (*T, error) {
	return RequestWithoutBodyAndDecode[T](baseUrl, path, http.MethodDelete)
}
