package Routing

import (
	"testing"

	"github.com/envuso/go-http/HttpContext"
	"github.com/envuso/go-http/Routing"
	"github.com/envuso/go-http/Routing/Route"
)

type TestMiddleware struct{}

func (mw *TestMiddleware) Handle(ctx *HttpContext.RequestContext) {
	ctx.Params().Set("message", "hello there")
}

func TestDefiningMiddleware(t *testing.T) {
	router := Routing.NewRouterHandler()
	router.AddMiddleware("hi", new(TestMiddleware))

	if !router.Middlewares().Has("hi") {
		t.Errorf("Test middleware not defined...")
	}

	mw := router.Middlewares().Get("hi")
	if mw == nil {
		t.Errorf("Test middleware not defined...")
	}
}

func TestCallingMiddlewareHandleFunc(t *testing.T) {
	router := Routing.NewRouterHandler()
	router.AddMiddleware("hi", new(TestMiddleware))

	mw := router.Middlewares().Get("hi")
	if mw == nil {
		t.Errorf("Test middleware not defined...")
	}
	ctx := HttpContext.NewBasicRequestContext()
	// nextFunc := func(ctx *HttpContext.RequestContext) {
	// 	param := ctx.Params().Get("message")
	// 	if param != "hello there" {
	// 		t.Errorf("CTX param('message') is not 'hello there'")
	// 	}
	// }

	mw.Handle(ctx)
}

func TestDefiningMiddlewareOnRoute(t *testing.T) {
	router := Routing.NewRouterHandler()
	router.AddMiddleware("hi", new(TestMiddleware))

	route := router.Get("/yeet", func() {}).Middleware("hi")

	if !route.HasMiddleware("hi") {
		t.Errorf("'hi' middleware not found on route.")
	}

}

func TestDefiningMiddlewareOnGroup(t *testing.T) {
	router := Routing.NewRouterHandler()
	router.AddMiddleware("hi", new(TestMiddleware))
	router.AddMiddleware("hi-two", new(TestMiddleware))

	group := router.Group(func(stack *Route.RouteGroupStack) {
		stack.Get("/yeet", func() {}).Middleware("hi-two")
	}).Middleware("hi")

	router.Build()

	if !group.HasMiddleware("hi") {
		t.Errorf("'hi' middleware not defined on route group")
	}

	route := router.FindForMethod("GET", "/yeet")

	if !route.HasMiddleware("hi") {
		t.Errorf("'hi' middleware not defined on route")
	}
	if !route.HasMiddleware("hi-two") {
		t.Errorf("'hi-two' middleware not defined on route")
	}

}
