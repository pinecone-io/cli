package network

import (
	"net/http"
)

func GetAndDecode[T any](baseUrl string, path string) (*T, error) {
	return RequestWithoutBodyAndDecode[T](baseUrl, path, http.MethodGet)
}
