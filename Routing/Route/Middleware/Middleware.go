package Middleware

import (
	"reflect"

	"github.com/envuso/go-http/HttpContext"
	"github.com/envuso/go-http/Reflection"
)

type MiddlewareHandlerFunc = func(ctx *HttpContext.RequestContext)

var Middlewares = NewMiddlewareList()

type Middleware interface {
	Handle(ctx *HttpContext.RequestContext)
}
type MiddlewareWithAfter interface {
	Handle(ctx *HttpContext.RequestContext)
	HandleAfter(ctx *HttpContext.RequestContext)
}

func FromParam(middleware any) Middleware {
	switch middleware.(type) {
	case Middleware:
		return middleware.(Middleware)
	case string:
		if mw, ok := Middlewares.GetOk(middleware.(string)); ok {
			return mw
		}
	}

	return nil
}
func FromParamWithName(middleware any) (Middleware, string) {
	switch middleware.(type) {
	case Middleware:
		mwFunc := middleware.(Middleware)
		name := Reflection.IndirectType(reflect.TypeOf(mwFunc)).Name()
		return mwFunc, name
	case string:
		if mw, ok := Middlewares.GetOk(middleware.(string)); ok {
			return mw, middleware.(string)
		}
	}

	return nil, ""
}

func ArrayFromVariadic(middlewares ...any) []Middleware {
	mws := []Middleware{}

	for _, middleware := range middlewares {
		mw := FromParam(middleware)
		if mw != nil {
			mws = append(mws, mw)
		}
	}

	return mws
}

func MapFromVariadic(middlewares ...any) map[string]Middleware {
	mws := make(map[string]Middleware)

	for _, middleware := range middlewares {
		mw, name := FromParamWithName(middleware)

		if mw != nil {
			mws[name] = mw
		}
	}

	return mws
}
