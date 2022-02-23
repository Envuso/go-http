package Http

import (
	"reflect"

	"gohttp/HttpContext"
)

type ResponseProcessor struct {
	out  []reflect.Value
	ctx  *HttpContext.RequestContext
	sent bool
}

func NewResponseProcessor(out []reflect.Value, ctx *HttpContext.RequestContext) *ResponseProcessor {
	return &ResponseProcessor{
		out:  out,
		ctx:  ctx,
		sent: false,
	}
}

func (res *ResponseProcessor) HasOutput() bool {
	return len(res.out) > 0
}

// Process when a request route is handled, it can have (almost)any response type
// We'll try to process a response type for it's return type.
//goland:noinspection VacuumLines
func (res *ResponseProcessor) Process() {
	if !res.HasOutput() {
		res.sendEarly()
		return
	}

	for i := 0; i < len(res.out); i++ {
		outType := res.out[i]

		var out interface{}
		if outType.Kind() == reflect.Ptr || outType.Kind() == reflect.Interface {
			out = outType.Elem().Interface()
		} else {
			outType.Interface()
		}

		// out := outType.Elem().Interface()

		if res.isErrorLikeResponse(out) {
			res.sendError(out)
			break
		}

		if res.isResponse(out) {
			res.sendResponse(out)
			break
		}

		if res.isJsonLikeResponse(outType) {
			res.sendJson(outType.Interface())
			break
		}

	}

	if !res.sent {
		res.ctx.Response.SendNoContent()
		res.sent = true
	}
}

// sendEarly TODO: Better name dafuq
func (res *ResponseProcessor) sendEarly() {
	if res.ctx.Response.CanSend() {
		res.ctx.Response.Send()
		return
	}

	if !res.ctx.Response.IsSent() {
		res.ctx.Response.SendNoContent()
		return
	}
}

func (res *ResponseProcessor) isJsonLikeResponse(outType reflect.Value) bool {
	outKind := outType.Kind()

	if outKind == reflect.Ptr {
		outKind = reflect.Indirect(outType).Kind()
	}

	switch outKind {
	case reflect.Map, reflect.Array, reflect.Struct:
		return true
	default:
		return false
	}
}

func (res *ResponseProcessor) sendJson(data interface{}) {
	res.ctx.Response.Json(data)
	res.ctx.Response.Send()
	res.sent = true
}

func (res *ResponseProcessor) isErrorLikeResponse(outVal interface{}) bool {
	_, ok := outVal.(error)

	return ok
}

func (res *ResponseProcessor) sendError(out any) {
	res.ctx.SendError(out.(error))
	res.sent = true
}

func (res *ResponseProcessor) isResponse(out any) bool {
	_, ok := out.(HttpContext.Response)

	return ok
}

func (res *ResponseProcessor) sendResponse(out any) {
	response := out.(HttpContext.Response)

	if response.CanSend() {
		response.Send()
		res.sent = true
		return
	}

	response.SendNoContent()
	res.sent = true
}

func (res *ResponseProcessor) setResponse(vals []reflect.Value) {
	res.out = vals
}
