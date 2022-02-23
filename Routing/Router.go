package Routing

import (
	"log"
	"net/http"
	"reflect"

	"gohttp/Routing/Route"
	"gohttp/Routing/Route/Middleware"
	"gohttp/Utilility"
)

type RouterContract interface {
	AddMiddleware(name string, handler Middleware.Middleware)
	Middlewares() *Middleware.MiddlewareList
	AddMiddlewares(middlewares Middleware.MiddlewareList)
	Get(path string, handlerArgs ...interface{}) *Route.Route
	Post(path string, handlerArgs ...interface{}) *Route.Route
	Group(stackHandler RouteGroupStackFunc) *Route.RouteGroup
	ControllerGroup(controller interface{}, stackHandler RouteGroupStackFunc) *Route.RouteGroup
	Add(httpMethod, path string, handlerArgs ...interface{}) *Route.Route
	AddRoute(route *Route.Route) *Route.Route
	GetRouteRegistrations() map[string]*RouteRegistrar
	GetRouteRegistrationsForMethod(httpMethod string) *RouteRegistrar
	FindForMethod(httpMethod, path string) *Route.RouteMatch
	Build()
}

var httpMethods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
}

type RouteGroupStackFunc = func(stack *Route.RouteGroupStack)

type RouterService struct {
	routes map[string]*RouteRegistrar
	groups []*Route.RouteGroup
}

func NewRouterHandler() *RouterService {
	router := &RouterService{
		routes: make(map[string]*RouteRegistrar),
		groups: []*Route.RouteGroup{},
	}

	for _, httpMethod := range httpMethods {
		router.routes[httpMethod] = NewRouteRegistrar(httpMethod)
	}

	return router
}

var Router = NewRouterHandler()

func (router *RouterService) AddMiddleware(name string, handler Middleware.Middleware) {
	Middleware.Middlewares.Set(name, handler)

}

func (router *RouterService) Middlewares() *Middleware.MiddlewareList {
	return Middleware.Middlewares
}

func (router *RouterService) AddMiddlewares(middlewares Middleware.MiddlewareList) {
	for name, middleware := range middlewares.Data() {
		router.AddMiddleware(name, middleware)
	}
}

func (router *RouterService) Get(path string, handlerArgs ...interface{}) *Route.Route {
	return router.Add(http.MethodGet, path, handlerArgs...)
}

func (router *RouterService) Post(path string, handlerArgs ...interface{}) *Route.Route {
	return router.Add(http.MethodPost, path, handlerArgs...)
}

func (router *RouterService) Group(stackHandler RouteGroupStackFunc) *Route.RouteGroup {
	group := Route.CreateRouteGroup()

	router.groups = append(router.groups, group)

	stackHandler(group.Stack())

	return group
}

func (router *RouterService) ControllerGroup(controller interface{}, stackHandler RouteGroupStackFunc) *Route.RouteGroup {
	controllerType := reflect.TypeOf(controller)

	if !Route.IsValidController(controllerType) {
		log.Printf("ControllerGroup error. First arg must be a controller struct... you passed: %v", controller)
		return nil
	}

	group := Route.CreateRouteControllerGroup(controllerType)

	router.groups = append(router.groups, group)

	stackHandler(group.Stack())

	return group
}

func (router *RouterService) Add(httpMethod, path string, handlerArgs ...interface{}) *Route.Route {
	route := Route.CreateDynamicRoute(httpMethod, path, handlerArgs...)

	return router.routes[route.HttpMethod()].addRoute(route).route
}

func (router *RouterService) AddRoute(route *Route.Route) *Route.Route {
	return router.routes[route.HttpMethod()].addRoute(route).route
}

func (router *RouterService) GetRouteRegistrations() map[string]*RouteRegistrar {
	return router.routes
}

func (router *RouterService) GetRouteRegistrationsForMethod(httpMethod string) *RouteRegistrar {
	if !Utilility.IsValidHttpMethod(httpMethod) {
		log.Printf("Invalid http method passed to GetRouteRegistrationsForMethod")
		return nil
	}

	return router.routes[httpMethod]
}

// FindForMethod Similar to RouteRegistrar.Find, except we'll provide the http method to lookup
func (router *RouterService) FindForMethod(httpMethod, path string) *Route.RouteMatch {
	if !Utilility.IsValidHttpMethod(httpMethod) {
		log.Printf("Invalid http method passed to FindForMethod")
		return nil
	}
	return router.GetRouteRegistrationsForMethod(httpMethod).Find(path)
}

func (router *RouterService) Build() {
	for _, group := range router.groups {
		for _, route := range group.Stack().Routes() {
			createdRoute := router.AddRoute(route)
			createdRoute.Middlewares.MergeIn(group.Middlewares)
		}
	}
}
