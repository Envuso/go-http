package Http

import (
	"net/http"

	"github.com/envuso/go-http/Routing"
)

type HttpContract interface {
	UsingRouter(router Routing.RouterContract) HttpContract
	Listen(addr string) error
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
}
