package Route

import (
	"reflect"

	"gohttp/Routing/Route/Middleware"
)

type Route struct {
	Middleware.UsesMiddlewares

	httpMethod string
	path       string
	handler    interface{}
	name       string
	prefix     string

	key   string
	param RouteParameterMatch

	controllerHandler *ControllerRoute
}

func CreateRoute(httpMethod string, path string, handler interface{}) *Route {
	r := &Route{
		httpMethod: httpMethod,
		path:       path,
		handler:    handler,
	}
	r.Created()
	return r
}

func CreateRouteRegistration(httpMethod string, route *Route, param RouteParameterMatch, key string) *Route {
	r := &Route{
		httpMethod: httpMethod,
		path:       route.Path(),
		handler:    route.Handler(),
		param:      param,
		key:        key,
	}
	r.Created()
	r.Middlewares.SetData(route.Middlewares.Data())
	r.controllerHandler = route.controllerHandler
	return r
}

func (route *Route) Name(name string) *Route {
	route.name = name

	return route
}

func (route *Route) Prefix(prefix string) *Route {
	route.prefix = prefix

	return route
}

func (route *Route) HttpMethod() string {
	return route.httpMethod
}

func (route *Route) Key() string {
	return route.key
}

func (route *Route) Param() RouteParameterMatch {
	return route.param
}

func (route *Route) Path() string {
	return route.path
}

func (route *Route) Handler() interface{} {
	return route.handler
}

func (route *Route) Middleware(middlewares ...any) *Route {
	route.UsesMiddlewares.Middleware(middlewares...)

	return route
}

func (route *Route) HasMiddleware(middleware any) bool {
	return route.UsesMiddlewares.HasMiddleware(middleware)
}

func (route *Route) Controller(controllerType reflect.Type, name string) *Route {
	route.controllerHandler = NewControllerRoute(controllerType, name)
	return route
}
