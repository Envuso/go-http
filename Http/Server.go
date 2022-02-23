package Http

import (
	"log"
	"net/http"
	"reflect"
	"sync/atomic"

	"github.com/envuso/go-http/HttpContext"
	"github.com/envuso/go-http/Routing"
	"github.com/envuso/go-http/Routing/Route"
	"github.com/envuso/go-http/Routing/Route/Middleware"
)

// Move this shit to HttpContext once i can refactor that

var globalRequestId uint64

type Http struct {
	router Routing.RouterContract
}

func NewHttp() HttpContract {
	return &Http{}
}

func (h *Http) UsingRouter(router Routing.RouterContract) HttpContract {
	h.router = router

	return h
}

func (h *Http) NextRequestId() uint64 {
	return atomic.AddUint64(&globalRequestId, 1)
}

func (h *Http) Listen(addr string) error {

	err := http.ListenAndServe(addr, h)
	if err != nil {
		return err
	}

	return nil
}

func (h *Http) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	id := h.NextRequestId()

	ctx := HttpContext.NewRequestContext(writer, request, id)
	ctx.Container.Singleton(func() *HttpContext.RequestContext {
		return ctx
	})
	ctx.Container.Singleton(new(HttpContext.Request), func(c *HttpContext.RequestContext) *HttpContext.Request {
		return c.GetRequest()
	})
	ctx.Container.Singleton(new(HttpContext.Response), func(c *HttpContext.RequestContext) *HttpContext.Response {
		return c.GetResponse()
	})
	ctx.Container.Bind(func(c *HttpContext.RequestContext) Route.Controller {
		return Route.Controller{Context: c}
	})

	registrar := h.router.GetRouteRegistrationsForMethod(request.Method)
	if registrar == nil {
		ctx.Response.SendNotFound()
		return
	}

	route := registrar.Find(request.URL.Path)
	if route == nil {
		ctx.Response.SendNotFound()
		return
	}

	// defer func(ctx *HttpContext.RequestContext) {
	// 	if r := recover(); r != nil {
	// 		log.Printf("Recovered from panic %v", r)
	// 		if !ctx.Response.IsSent() {
	// 			ctx.SendError(errors.New("Something went wrong.... woops."))
	// 		}
	// 	}
	// }(ctx)

	// Loop through all route middlewares
	var afterMiddlewareHandlers = []Middleware.MiddlewareHandlerFunc{}
	if !route.Middlewares.IsEmpty() {
		afterMiddlewareHandlers = h.ProcessMiddleware(ctx, route)
	}

	// We need to use the old reflection code to store/get the controller for this route
	responseValues := []reflect.Value{}
	responseProcessor := NewResponseProcessor(responseValues, ctx)
	var err error

	if route.Handler == nil {
		// Instantiate a new controller instance
		// Call the handler via reflection, injecting any method params
		err, responseValues = route.ControllerRoute.CallMethod(ctx, h.methodDependencyResolver())

		responseProcessor.setResponse(responseValues)
	}

	if err != nil {
		panic(err)
		// TODO: Figure something for this, maybe throw it through the response processor?
	}

	// Inject any dependencies
	// Call the handler via reflection, injecting any method params
	// Process the dynamic response...
	// Reeeee

	if route.Handler != nil {
		handler := Route.NewClosureRoute(route.Handler)
		err, responseValues = handler.CallClosure(ctx, h.methodDependencyResolver())
		responseProcessor.setResponse(responseValues)

		// route.Handler.(func())()
	}

	// Run the after middlewares before we send our response
	// This allows us to do final manipulations before...
	h.RunAfterMiddlewares(ctx, afterMiddlewareHandlers)

	// Finally... process the response types & send it
	responseProcessor.Process()
}

func (h *Http) ProcessMiddleware(context *HttpContext.RequestContext, route *Route.RouteMatch) []Middleware.MiddlewareHandlerFunc {
	middlewares := route.Middlewares.Values()
	size := route.Middlewares.Length() - 1
	afterMiddlewares := []Middleware.MiddlewareHandlerFunc{}

	for i := 0; i <= size; i++ {
		mw := middlewares[i]
		mw.Handle(context)
		log.Printf("Middleware %d run", i)

		if after, ok := mw.(Middleware.MiddlewareWithAfter); ok {
			log.Printf("Middleware %d has after method", i)
			afterMiddlewares = append(afterMiddlewares, after.HandleAfter)
		}
	}

	return afterMiddlewares
}

func (h *Http) RunAfterMiddlewares(context *HttpContext.RequestContext, middlewares []Middleware.MiddlewareHandlerFunc) {
	for _, middleware := range middlewares {
		middleware(context)
	}
}

func (h *Http) methodDependencyResolver() Route.MethodDependencyResolver {
	return func(ctx *HttpContext.RequestContext, dependencies *Route.RouteDependencies, method reflect.Value) []reflect.Value {
		methodType := method.Type()

		for i := dependencies.ProvidedCount; i < dependencies.RequiredCount; i++ {
			inType := methodType.In(i)

			if HttpContext.IsRequestBodyStruct(inType) {
				dependencies.Dependencies[i] = HttpContext.ResolveRequestBody(inType, ctx)
				continue
			}

			inst, err := ctx.ResolveInstanceFromContainer(inType)
			if err != nil {
				log.Printf("Failed to resolve method dep from container, %v - type: %v", err, inType)
				continue
			}
			if inst.IsValid() && !inst.IsNil() {
				dependencies.Dependencies[i] = inst
				continue
			}

			// if inType.AssignableTo(HttpContext.IRequestContextTypeElem) {
			// 	dependencies[i] = reflect.ValueOf(ctx)
			// 	continue
			// }

		}

		return dependencies.Dependencies
	}
}
