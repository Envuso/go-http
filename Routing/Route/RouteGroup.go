package Route

import (
	"reflect"

	"github.com/envuso/go-http/Routing/Route/Middleware"
)

type RouteGroup struct {
	Middleware.UsesMiddlewares

	name   string
	prefix string
	stack  *RouteGroupStack

	// This is set when we're creating a controller group
	controllerType reflect.Type
}

func CreateRouteGroup() *RouteGroup {
	rg := &RouteGroup{
		name:   "",
		prefix: "",
		stack:  CreateRouteGroupStack(),
	}
	rg.Created()
	return rg
}

func CreateRouteControllerGroup(controllerType reflect.Type) *RouteGroup {
	rg := CreateRouteGroup()
	rg.controllerType = controllerType
	rg.stack.routeHandling = &routeHandler{
		handlerType: "controller",
		controller:  controllerType,
	}

	return rg
}

func (group *RouteGroup) Stack() *RouteGroupStack {
	return group.stack
}

func (group *RouteGroup) Name(name string) *RouteGroup {
	group.name = name
	return group
}

func (group *RouteGroup) Middleware(middlewares ...any) *RouteGroup {
	group.UsesMiddlewares.Middleware(middlewares...)
	return group
}

func (group *RouteGroup) HasMiddleware(middleware any) bool {
	return group.UsesMiddlewares.HasMiddleware(middleware)
}
