package network

import (
	"net/http"
)

func GetAndDecode[T any](baseUrl string, path string, useApiKey bool) (*T, error) {
	return RequestWithoutBodyAndDecode[T](baseUrl, path, http.MethodGet, useApiKey)
}
