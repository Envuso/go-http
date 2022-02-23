package Utilility

import (
	"net/http"
	"strings"
)

var HttpMethodLookup = map[string]bool{
	http.MethodGet:    true,
	http.MethodHead:   true,
	http.MethodPost:   true,
	http.MethodPut:    true,
	http.MethodPatch:  true,
	http.MethodDelete: true,
}

func IsValidHttpMethod(method string) bool {
	method = strings.ToUpper(method)

	_, ok := HttpMethodLookup[method]

	return ok
}
