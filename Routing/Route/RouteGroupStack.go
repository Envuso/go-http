package Route

import (
	"fmt"
	"net/http"
	"reflect"
)

type RouteGroupStack struct {
	routes        []*Route
	routeHandling *routeHandler
}

func (group *RouteGroupStack) Routes() []*Route {
	return group.routes
}

func CreateRouteGroupStack() *RouteGroupStack {
	return &RouteGroupStack{
		routes: []*Route{},
	}
}

func (group *RouteGroupStack) isControllerGroup() bool {
	if group.routeHandling == nil {
		return false
	}

	return group.routeHandling.handlerType == "controller" && group.routeHandling.controller != nil
}

func (group *RouteGroupStack) addRouteToStack(route *Route) *Route {
	group.routes = append(group.routes, route)
	return group.routes[len(group.routes)-1]
}

func (group *RouteGroupStack) getHandlerArgs(handlerArgs ...interface{}) []interface{} {
	returnArgs := []interface{}{}

	if len(handlerArgs) == 1 {
		handlerArg := handlerArgs[0]

		// If we pass a string, and we're working with a controller group
		// this is our reference to our method name on the controller
		// With this, we'll then set the first dynamic arg to the controller type
		// And the second, we'll pass the controller method name.
		if controllerMethodName, ok := handlerArg.(string); ok && group.isControllerGroup() {
			returnArgs = append(returnArgs, group.routeHandling.controller)
			returnArgs = append(returnArgs, controllerMethodName)

			return returnArgs
		}

		handlerType := reflect.TypeOf(handlerArg)
		if handlerType.Kind() == reflect.Func {
			returnArgs = append(returnArgs, handlerArg)
			return returnArgs
		}

	}

	fmt.Printf("RouteGroupStack.getHandlerArgs has more than 1 arg... return args are: %v .... %v", handlerArgs, returnArgs)

	return handlerArgs
}

func (group *RouteGroupStack) createStackRoute(method string, path string, handlerArgs ...interface{}) *Route {
	return CreateDynamicRoute(method, path, group.getHandlerArgs(handlerArgs...)...)
}

func (group *RouteGroupStack) Get(path string, handlerArgs ...interface{}) *Route {
	return group.addRouteToStack(group.createStackRoute(http.MethodGet, path, handlerArgs...))
}

func (group *RouteGroupStack) Post(path string, handlerArgs ...interface{}) *Route {
	return group.addRouteToStack(group.createStackRoute(http.MethodPost, path, handlerArgs...))
}
