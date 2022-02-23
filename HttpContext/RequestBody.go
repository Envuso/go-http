package HttpContext

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/envuso/go-http/Reflection"
	"github.com/tidwall/gjson"
)

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

var RequestBodyType = reflect.TypeOf(RequestBody{})

func IsRequestBodyStruct(param reflect.Type) bool {
	paramType := param

	if param.Kind() == reflect.Ptr {
		paramType = param.Elem()
	}

	if paramType.Kind() != reflect.Struct {
		return false
	}

	fields := paramType.NumField()

	for i := 0; i < fields; i++ {
		field := paramType.Field(i)
		if Reflection.IndirectType(field.Type) == Reflection.IndirectType(RequestBodyType) {
			return true
		}
	}

	return false
}

func ResolveRequestBody(param reflect.Type, ctx *RequestContext) reflect.Value {
	val := reflect.New(Reflection.IndirectType(param)).Interface()

	if rbody, ok := val.(IRequestBody); ok {
		ctx.Request.reader.Body = http.MaxBytesReader(ctx.Response.writer, ctx.Request.reader.Body, 1048576)
		rbody.SetBody(ctx.Request.reader.Body)
		rbody.Decode(&val)
	}

	if param.Kind() == reflect.Ptr {
		return reflect.ValueOf(val)
	}

	return reflect.ValueOf(val).Elem()
}

type IRequestBody interface {
	SetBody(data any)
	Decode(str interface{}) *malformedRequest
	Get(path string) gjson.Result
	Has(path string) bool
	Json() gjson.Result
}

type RequestBody struct {
	body       []byte
	jsonParsed gjson.Result
	data       any
}

func (b *RequestBody) Get(path string) gjson.Result {
	return b.jsonParsed.Get(path)
}

func (b *RequestBody) Has(path string) bool {
	return b.jsonParsed.Get(path).Exists()
}

func (b *RequestBody) Json() gjson.Result {
	return b.jsonParsed
}

func (b *RequestBody) SetBody(data any) {
	b.data = data
	b.body = []byte{}
}

func (b *RequestBody) Decode(destination interface{}) *malformedRequest {
	data, err := ioutil.ReadAll(b.data.(io.ReadCloser))
	if err != nil {
		return &malformedRequest{status: http.StatusInternalServerError, msg: err.Error()}
	}
	b.data = nil
	b.body = data
	b.jsonParsed = gjson.Parse(string(b.body))

	err = json.Unmarshal(b.body, &destination)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return &malformedRequest{status: http.StatusInternalServerError, msg: err.Error()}
		}
	}

	return nil
}
