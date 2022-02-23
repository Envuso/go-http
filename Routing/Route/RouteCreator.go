package Route

import (
	"log"
	"reflect"

	"gohttp/Reflection"
	"gohttp/Utilility"
)

type routeHandler struct {
	handlerType string

	handler interface{}

	controller reflect.Type
	methodName string
}

func IsValidController(controller interface{}) bool {
	var handlerType reflect.Type

	if cType, ok := controller.(reflect.Type); ok {
		handlerType = Reflection.IndirectType(cType)
	} else {
		handlerType = Reflection.IndirectType(reflect.TypeOf(controller))
	}
	k := handlerType.Kind()
	return k == reflect.Struct
}

func getRouteHandler(handlerArgs ...interface{}) *routeHandler {
	handler := handlerArgs[0]

	var handlerType reflect.Type
	if handlerT, ok := handler.(reflect.Type); ok {
		handlerType = handlerT
	} else {
		handlerType = reflect.TypeOf(handler)
	}

	// We're passing a function reference I guess?
	if len(handlerArgs) == 1 {
		if handlerType.Kind() == reflect.Func {
			return &routeHandler{
				handlerType: "func",
				handler:     reflect.ValueOf(handler).Interface(),
			}
		}
		panic("wat")
	}

	// We're passing controller struct ref + method name?
	handlerTypeI := Reflection.IndirectType(handlerType)
	if !IsValidController(handlerType) {
		log.Printf("Arg 1 must be a controller reference")
		return nil
	}

	mName, mNameOk := handlerArgs[1].(string)
	if !mNameOk {
		log.Printf("Arg 2 must be a controller method name(string)")
		return nil
	}

	if _, ok := Reflection.StructHasMethod(handlerType, mName); !ok {
		log.Printf("Arg 2(controller method name) must be a exported method of the struct %s", Reflection.GetStructName(handlerTypeI))
		return nil
	}

	return &routeHandler{
		handlerType: "controller",
		handler:     nil,
		controller:  handlerTypeI,
		methodName:  mName,
	}
}

func CreateDynamicRoute(httpMethod, path string, handlerArgs ...interface{}) *Route {
	if !Utilility.IsValidHttpMethod(httpMethod) {
		log.Printf("Invalid http method passed to Add")
		return nil
	}

	handler := getRouteHandler(handlerArgs...)
	if handler == nil {
		return nil
	}

	var route *Route
	if handler.handlerType == "func" {
		route = CreateRoute(httpMethod, path, handler.handler)
	}
	if handler.handlerType == "controller" {
		route = CreateRoute(httpMethod, path, nil).Controller(handler.controller, handler.methodName)
	}

	return route
}
