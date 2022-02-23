package Route

import (
	"log"
	"reflect"

	container "github.com/Envuso/go-ioc-container"
	"gohttp/HttpContext"
	"gohttp/Reflection"
)

type ControllerRoute struct {
	methodName             string
	controllerType         reflect.Type
	controllerTypeIndirect reflect.Type

	processingType string
	closureValue   reflect.Value
	closureType    reflect.Type
}

func NewControllerRoute(controllerType reflect.Type, methodName string) *ControllerRoute {

	container.Bind(controllerType)

	return &ControllerRoute{
		methodName:             methodName,
		controllerType:         controllerType,
		controllerTypeIndirect: Reflection.IndirectType(controllerType),

		processingType: "controller",
	}
}

func NewClosureRoute(closure interface{}) *ControllerRoute {

	// container.Bind(reflect.TypeOf(closure), closure)

	return &ControllerRoute{
		processingType: "controller",
		closureValue:   reflect.ValueOf(closure),
		closureType:    reflect.TypeOf(closure),
	}
}

// func (route *ControllerRoute) InstantiateController() reflect.Value {
// 	return reflect.New(route.controllerTypeIndirect)
// }

func (route *ControllerRoute) GetControllerMethod(controllerInstance reflect.Value) reflect.Method {
	controllerMethod, ok := controllerInstance.Type().MethodByName(route.methodName)
	if !ok {
		log.Printf("Failed to find controller method...")
		return controllerMethod
	}

	return controllerMethod
}

type MethodDependencyResolver = func(ctx *HttpContext.RequestContext, dependencies *RouteDependencies, method reflect.Value) []reflect.Value

func (route *ControllerRoute) injectCtxField(ctx *HttpContext.RequestContext, controllerInst reflect.Value) reflect.Value {
	// Manually inject our context into the controller instance
	var ctxField reflect.Value
	if controllerInst.Kind() == reflect.Ptr {
		ctxField = controllerInst.Elem().FieldByName("Context")
	} else {
		ctxField = controllerInst.FieldByName("Context")
	}
	ctxType := reflect.ValueOf(ctx)
	if ctxField.IsValid() && ctxField.Type().AssignableTo(ctxType.Type()) {
		ctxField.Set(ctxType)
	}

	invocable := container.CreateInvocable(Reflection.IndirectType(controllerInst.Type()))

	return reflect.ValueOf(invocable.InstantiateWith(ctx.Container))
	// err := Ioc.TypeContainer.ResolveStructTypes(controllerInst, controllerInst)
	// if err != nil {
	// 	panic(err)
	// }

	// return controllerInst
}

type RouteDependencies struct {
	RequiredCount int
	ProvidedCount int
	Dependencies  []reflect.Value
}

func (route *ControllerRoute) CallMethod(ctx *HttpContext.RequestContext, dependencyResolver MethodDependencyResolver) (error, []reflect.Value) {

	invocable := container.CreateInvocable(route.controllerType)

	responseItems := invocable.CallMethodByNameWithArgInterceptor(
		route.methodName,
		ctx.Container,
		func(index int, argType reflect.Type, typeZeroVal reflect.Value) (reflect.Value, bool) {
			if HttpContext.IsRequestBodyStruct(argType) {
				return HttpContext.ResolveRequestBody(argType, ctx), true
			}

			return typeZeroVal, false
		},
	)
	// .Call(container.resolveFunctionArgs(method, parameters...))

	// invocable.CallMethodByNameWith(route.methodName, ctx.Container)

	// controllerInstance := route.injectCtxField(ctx, route.InstantiateController())
	// controllerMethod := route.GetControllerMethod(controllerInstance)
	//
	// responseItems := []reflect.Value{}
	//
	// controllerMethodFunc := controllerMethod.Func
	// if !controllerMethodFunc.IsValid() {
	// 	return errors.New("Failed to find controller method..."), responseItems
	// }
	//
	// requiredCount := controllerMethodFunc.Type().NumIn()
	//
	// for i := 0; i < controllerMethodFunc.Type().NumIn(); i++ {
	// 	t := container.ContainerTypes.Of(controllerMethodFunc.Type().In(i))
	// 	log.Printf("Type: %v", t)
	// }
	//
	// dependencies := make([]reflect.Value, requiredCount)
	// dependencies[0] = controllerInstance
	// rDependencies := &RouteDependencies{
	// 	RequiredCount: requiredCount,
	// 	ProvidedCount: 1,
	// 	Dependencies:  dependencies,
	// }
	//
	// args := dependencyResolver(ctx, rDependencies, controllerMethodFunc)
	//
	// responseItems = controllerMethodFunc.Call(args)

	return nil, responseItems
}

func (route *ControllerRoute) CallClosure(ctx *HttpContext.RequestContext, dependencyResolver MethodDependencyResolver) (error, []reflect.Value) {
	responseItems := []reflect.Value{}

	requiredCount := route.closureType.NumIn()
	rDependencies := &RouteDependencies{
		RequiredCount: requiredCount,
		ProvidedCount: 0,
		Dependencies:  make([]reflect.Value, requiredCount),
	}

	args := dependencyResolver(ctx, rDependencies, route.closureValue)

	responseItems = route.closureValue.Call(args)

	return nil, responseItems
}
