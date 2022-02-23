package Route

import "gohttp/Routing/Route/Middleware"

type RouteMatch struct {
	Middleware.UsesMiddlewares

	Method          string
	Path            string
	Handler         interface{}
	Name            string
	Prefix          string
	Param           RouteParameterMatch
	Params          *RouteParameters[string]
	key             string
	ControllerRoute *ControllerRoute
}

func CreateRouteMatch(route Route, params map[string]string) *RouteMatch {
	return &RouteMatch{
		UsesMiddlewares: Middleware.UsesMiddlewares{Middlewares: route.Middlewares},
		Method:          route.httpMethod,
		Path:            route.path,
		Handler:         route.handler,
		Name:            route.name,
		Prefix:          route.prefix,
		Param:           route.param,
		Params:          CreateRouteParameters(params),
		ControllerRoute: route.controllerHandler,
		key:             route.key,
	}
}
