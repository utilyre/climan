package httpparser

import "net/http"

func isValidMethod(method string) bool {
	return method == http.MethodGet ||
		method == http.MethodHead ||
		method == http.MethodPost ||
		method == http.MethodPut ||
		method == http.MethodDelete ||
		method == http.MethodConnect ||
		method == http.MethodOptions ||
		method == http.MethodTrace ||
		method == http.MethodPatch
}
