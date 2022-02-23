package HttpContext

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"

	container "github.com/Envuso/go-ioc-container"
)

type RequestContext struct {
	Response  *Response
	Request   *Request
	Container *container.ContainerInstance
	requestId uint64

	routeParams *Parameters[string]
	queryParams *QueryParameters
	body        Parameters[string]
}

func NewBasicRequestContext() *RequestContext {
	return &RequestContext{
		routeParams: &Parameters[string]{
			data: make(map[string]string),
		},
	}
}

func NewRequestContext(writer http.ResponseWriter, reader *http.Request, requestId uint64) *RequestContext {

	ctx := &RequestContext{
		Response:  NewResponse(writer),
		Request:   NewRequest(reader),
		Container: container.CreateChildContainer(),
		requestId: requestId,

		routeParams: ParametersFromValues(nil),
		queryParams: QueryParametersFrom(reader.URL.Query()),
	}

	encoder := Encoder.ContentTypeEncoderForRequest(ctx.Request)
	ctx.Request.encoder = encoder.Request
	ctx.Response.encoder = encoder.Response

	ctx.Request.requestId = requestId

	ctx.Container.Singleton(func() *Response {
		return ctx.Response
	})
	ctx.Container.Singleton(func() *Request {
		return ctx.Request
	})

	return ctx
}

func (r *RequestContext) GetRequest() *Request {
	return r.Request
}
func (r *RequestContext) GetResponse() *Response {
	return r.Response
}
func (r *RequestContext) Reader() *http.Request {
	return r.Request.reader
}
func (r *RequestContext) Writer() http.ResponseWriter {
	return r.Response.writer
}
func (r *RequestContext) Params() *Parameters[string] {
	return r.routeParams
}

func (r *RequestContext) GetRequestId() uint64 {
	return r.requestId
}
func (r *RequestContext) IncrReqId() uint64 {
	r.requestId += 1

	return r.requestId
}

func (r *RequestContext) Id() uint64 {
	return r.requestId << 32
}

func (r *RequestContext) IdString() string {
	return strconv.Itoa(int(r.Id()))
}

func (r *RequestContext) SendError(err error) {
	r.Response.SendError(err)
}

// ResolveInstanceFromContainer If our context container(which is on a per-request basis) has
// an instance of this type, we'll resolve it, if not, we'll look up the type in the
// application container
func (r *RequestContext) ResolveInstanceFromContainer(t reflect.Type) (reflect.Value, error) {
	val := reflect.ValueOf(r.Container.Make(t))
	if val.IsValid() {
		return val, nil
	}

	return reflect.ValueOf(nil), errors.New("type not found in either container...")
}

var IRequestContextType = reflect.TypeOf((*RequestContextContract)(nil))
var IRequestContextTypeElem = IRequestContextType.Elem()

type RequestContextContract interface {
	GetRequest() *Request
	GetResponse() *Response
	Reader() *http.Request
	Writer() http.ResponseWriter
	Params() *Parameters[string]
	GetRequestId() uint64
	Id() uint64
	IdString() string
	SendError(err error)
	ResolveInstanceFromContainer(t reflect.Type) (reflect.Value, error)
}
