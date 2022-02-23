package Routing

import (
	"log"
	"net/http"
	"testing"

	"gohttp/Routing"
	"gohttp/Routing/Route"
)

func TestRegistration(t *testing.T) {
	router := Routing.NewRouterHandler()

	router.Get("/test", func() string {
		return "hi!"
	})

	registrar := router.GetRouteRegistrationsForMethod(http.MethodGet)

	match := registrar.Find("/test")

	if match == nil {
		t.Errorf("Could not find the only registered route.")
	}
}

func TestRegistrationWithParam(t *testing.T) {
	router := Routing.NewRouterHandler()

	router.Get("/test/{username}", func() string {
		return "hi!"
	})
	registrar := router.GetRouteRegistrationsForMethod(http.MethodGet)
	match := registrar.Find("/test/samuel")

	if match == nil {
		t.Errorf("Could not find the only registered route.")
	}

	if !match.Params.Has("username") {
		t.Errorf(":username param does not exist")
	}
	if !match.Params.Has("username", "samuel") {
		t.Errorf(":username param does not match provided string.")
	}
}

func TestDeepRegistrationWithParam(t *testing.T) {
	router := Routing.NewRouterHandler()
	router.Get("/user/{username}/view/{type}", func() string {
		return "hi!"
	}).Name("hello!")
	registrar := router.GetRouteRegistrationsForMethod(http.MethodGet)
	match := registrar.Find("/user/samuel/view/profile")

	if match == nil {
		t.Errorf("Could not find the only registered route.")
	}

	if !match.Params.Has("username", "samuel") {
		t.Errorf(":username param does not match provided string.")
	}
	if !match.Params.Has("type", "profile") {
		t.Errorf(":type param does not match provided string.")
	}
}

func TestBasicRouteGrouping(t *testing.T) {
	router := Routing.NewRouterHandler()

	router.Group(func(stack *Route.RouteGroupStack) {
		stack.Get("/user", func() {})
	}).Name("user.")

	router.Build()

	route := router.FindForMethod(http.MethodGet, "/user")

	if route == nil {
		t.Errorf("Could not find the only registered route.")
	}

	if route.Path != "/user" {
		t.Errorf("Could not find route /user")
	}
}

type ControllerTest struct{}

func (c *ControllerTest) Hello() {
	log.Printf("Hello from ControllerTest.Hello")
}

func TestRouteFromController(t *testing.T) {
	router := Routing.NewRouterHandler()

	router.Get("/hello", new(ControllerTest), "Hello")

	router.Build()

	route := router.FindForMethod(http.MethodGet, "/hello")

	if route == nil {
		t.Errorf("Could not find the only registered route.")
	}

	if route.Path != "/hello" {
		t.Errorf("Could not find route /hello")
	}

	// route.Handler.(func(*ControllerTest))(new(ControllerTest))
}
