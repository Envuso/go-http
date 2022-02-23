package Routing

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"gohttp/Http"
	"gohttp/HttpContext"
	"gohttp/Routing"
)

func createResRouter() *Routing.RouterService {
	resRouter := Routing.NewRouterHandler()
	resRouter.Get("/error", func() error {
		return errors.New("This is an error response.")
	})

	resRouter.Build()

	return resRouter
}

var responseRouter = createResRouter()

type TestMiddlewareWithAfter struct{}

func (mw *TestMiddlewareWithAfter) Handle(ctx *HttpContext.RequestContext) {
	ctx.Params().Set("message", "hello there")
}
func (mw *TestMiddlewareWithAfter) HandleAfter(ctx *HttpContext.RequestContext) {
	log.Printf("TestMiddlewareWithAfter... after request")
}

func buildAndDispatch(router *Routing.RouterService, req *http.Request) *httptest.ResponseRecorder {
	httpInst := Http.NewHttp()
	httpInst.UsingRouter(router)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(httpInst.ServeHTTP)

	handler.ServeHTTP(rr, req)

	return rr
}

func buildAndDispatchReq(router *Routing.RouterService, method string, path string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}

	// req.Header.Set("content-type", "application/json")
	// req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "text/plain")
	req.Header.Set("accept", "text/plain")

	return buildAndDispatch(router, req)
}
func buildAndDispatchReqWithContentType(router *Routing.RouterService, method string, path string, contentType string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("content-type", contentType)
	req.Header.Set("accept", contentType)

	return buildAndDispatch(router, req)
}

func TestHttpHandling(t *testing.T) {
	router := Routing.NewRouterHandler()

	router.Get("/test", func() {
		log.Printf("hello from /test")
	})
	router.Build()

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := buildAndDispatch(router, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

}

func TestMiddlewaresProcessing(t *testing.T) {
	router := Routing.NewRouterHandler()
	router.AddMiddleware("hi", new(TestMiddleware))
	router.AddMiddleware("hi1", new(TestMiddleware))
	router.AddMiddleware("hi2", new(TestMiddleware))
	router.AddMiddleware("hi3", new(TestMiddlewareWithAfter))

	router.Get("/test", func() {
		print("hello from /test")
	}).Middleware("hi", "hi1", "hi2", "hi3")

	router.Build()

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := buildAndDispatch(router, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

}

type HttpControllerTest struct{}

func (c *HttpControllerTest) Hello() {
	log.Printf("Hello from HttpControllerTest.Hello")
}

type HelloRes struct {
	Message string `json:"message"`
}

func (c *HttpControllerTest) HelloWithJsonOne() HelloRes {
	return HelloRes{Message: "oh hai"}
}

func (c *HttpControllerTest) HelloWithContext(ctx *HttpContext.RequestContext) {
	if ctx == nil {
		panic("Context is not set.")
	}
}
func (c *HttpControllerTest) HelloWithRequest(req *HttpContext.Request) {
	if req == nil {
		panic("Request is not set.")
	}
}
func (c *HttpControllerTest) HelloWithResponse(res *HttpContext.Response) {
	if res == nil {
		panic("Response is not set.")
	}
}

func TestHttpWithControllerRoute(t *testing.T) {
	router := Routing.NewRouterHandler()
	router.AddMiddleware("hi", new(TestMiddleware))

	// TODO: Look into this shit more, kinda cool :D
	// reflect2.TypeByName()
	// n := reflect.TypeOf(HttpControllerTest{}).PkgPath()
	// tt := reflect2.TypeByPackageName("gohttp/Tests/Routing", "HttpControllerTest")
	// print(n, tt)

	route := router.Get("/test", HttpControllerTest{}, "Hello")
	route.Middleware("hi")

	router.Build()

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := buildAndDispatch(router, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

}

func TestHttpWithInjectedContext(t *testing.T) {
	router := Routing.NewRouterHandler()
	router.AddMiddleware("hi", new(TestMiddleware))
	router.Get("/ctx", HttpControllerTest{}, "HelloWithContext")
	router.Get("/req", HttpControllerTest{}, "HelloWithRequest")
	router.Get("/res", HttpControllerTest{}, "HelloWithResponse")
	router.Build()

	types := []string{"ctx", "req", "res"}
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

func TestRouteReturnsJson(t *testing.T) {
	router := Routing.NewRouterHandler()
	router.AddMiddleware("hi", new(TestMiddleware))
	router.Get("/test", HttpControllerTest{}, "HelloWithJsonOne")
	router.Build()

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := buildAndDispatch(router, req)

	expected := `{"message":"oh hai"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	if ctype := rr.Header().Get("content-type"); ctype != HttpContext.HEADER_CONTENT_TYPE_JSON {
		t.Errorf("handler return invalid content-type; got %v want %v", ctype, HttpContext.HEADER_CONTENT_TYPE_JSON)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

}

func TestErrorResponse(t *testing.T) {
	rr := buildAndDispatchReqWithContentType(
		responseRouter, "GET", "/error", HttpContext.HEADER_CONTENT_TYPE_JSON,
	)
	body := rr.Body.String()
	if body != `{"message":"This is an error response."}` {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
	if ctype := rr.Header().Get("content-type"); ctype != HttpContext.HEADER_CONTENT_TYPE_JSON {
		t.Errorf("handler return invalid content-type; got %v want %v", ctype, HttpContext.HEADER_CONTENT_TYPE_JSON)
	}
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	rr = buildAndDispatchReqWithContentType(
		responseRouter, "GET", "/error", HttpContext.HEADER_CONTENT_TYPE_TEXT_PLAIN,
	)
	body = rr.Body.String()
	if body != "Woops. Something went wrong!\nThis is an error response." {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
	if ctype := rr.Header().Get("content-type"); ctype != HttpContext.HEADER_CONTENT_TYPE_TEXT_PLAIN {
		t.Errorf("handler return invalid content-type; got %v want %v", ctype, HttpContext.HEADER_CONTENT_TYPE_TEXT_PLAIN)
	}
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

}
