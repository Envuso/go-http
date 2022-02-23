package Http

import (
	"net/http"

	"gohttp/Routing"
)

type HttpContract interface {
	UsingRouter(router Routing.RouterContract) HttpContract
	Listen(addr string) error
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
}
