package Routing

import (
	"net/http"
	"strings"
	"testing"

	"gohttp/HttpContext"
	"gohttp/Routing"
	"gohttp/Routing/Route"
)

type UserRequest struct {
	HttpContext.RequestBody

	Data struct {
		User struct {
			Name string `json:"name"`
		} `json:"user"`
	} `json:"data"`
}
type HttpControllerDITest struct {
	Route.Controller
}

func (c *HttpControllerDITest) HelloWithContext(ctx *HttpContext.RequestContext) {
	if ctx == nil {
		panic("Context is not set.")
	}
}
func (c *HttpControllerDITest) HelloWithControllerContext() {
	if c.Context == nil {
		panic("Context is not set.")
	}
}
func (c *HttpControllerDITest) HelloWithRequest(req *HttpContext.Request) {
	if req == nil {
		panic("Request is not set.")
	}
}
func (c *HttpControllerDITest) HelloWithResponse(res *HttpContext.Response) {
	if res == nil {
		panic("Response is not set.")
	}
}
func (c *HttpControllerDITest) HelloWithData(data UserRequest) {
	user := data.Get("data.user.name")

	if !user.Exists() || user.String() != "sam" {
		panic("Response is not set.")
	}
}

func TestControllerCtxInjection(t *testing.T) {
	router := Routing.NewRouterHandler()
	router.Get("/ctx", HttpControllerDITest{}, "HelloWithContext")
	router.Get("/ctrl-ctx", HttpControllerDITest{}, "HelloWithControllerContext")
	router.Get("/req", HttpControllerDITest{}, "HelloWithRequest")
	router.Get("/res", HttpControllerDITest{}, "HelloWithResponse")
	router.Build()

	types := []string{
		"ctx",
		"req",
		"res",
		"ctrl-ctx",
	}
	for _, typee := range types {
		req, err := http.NewRequest("GET", "/"+typee, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := buildAndDispatch(router, req)
		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
		}
	}

}

func TestControllerStructDataInjection(t *testing.T) {
	router := Routing.NewRouterHandler()
	router.Post("/data", HttpControllerDITest{}, "HelloWithData")
	router.Build()

	req, err := http.NewRequest("POST", "/data", strings.NewReader(`{"data":{"user":{"name":"sam"}}}`))
	if err != nil {
		t.Fatal(err)
	}
	rr := buildAndDispatch(router, req)
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

}
