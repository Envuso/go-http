package HttpContext

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type responseConfiguredState struct {
	statusCode bool
	data       bool
}

type Response struct {
	writer http.ResponseWriter

	StatusCode int
	Data       []byte

	Headers *HeaderContainer

	configured responseConfiguredState
	encoder    ContentEncoder
}

func NewResponse(writer http.ResponseWriter) *Response {
	return &Response{
		Headers: NewHeaderContainer(),

		writer: &DoneWriter{ResponseWriter: writer},
		configured: responseConfiguredState{
			statusCode: false,
			data:       false,
		},
	}
}

func (res *Response) IsSent() bool {
	if dw, ok := res.writer.(*DoneWriter); ok {
		return dw.Done
	}

	panic(errors.New("not a DoneWriter"))

	return false
}

func (res *Response) SetStatus(code int) {
	res.StatusCode = code
	res.configured.statusCode = true
}

func (res *Response) SetRawData(data []byte) *Response {
	res.Data = data
	res.configured.data = true

	return res
}

func (res *Response) Json(data interface{}) *Response {
	res.writer.Header().Set("Content-Type", HEADER_CONTENT_TYPE_JSON)
	res.writer.Header().Set("X-Content-Type-Options", "nosniff")

	dat, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode json response: %s", err)
		return res
	}

	return res.SetRawData(dat)
}

func (res *Response) SendNoContent() *Response {
	res.writer.WriteHeader(http.StatusNoContent)

	return res
}

func (res *Response) SendNotFound() *Response {
	res.SetStatus(http.StatusNotFound)
	res.Json(map[string]string{"message": "Not found"})
	res.Send()
	return res
}

func (res *Response) CanSend() bool {
	return !res.IsSent() && res.configured.data
}

func (res *Response) Send() {
	if !res.configured.statusCode {
		res.SetStatus(http.StatusOK)
	}

	if !res.Headers.Has("content-type") {
		res.Headers.SetFrom(res.encoder.HeadersForType())
	}

	if !res.Headers.IsEmpty() {
		for _, headerKey := range res.Headers.Keys() {
			res.writer.Header().Set(headerKey, res.Headers.Get(headerKey))
		}
	}

	res.writer.WriteHeader(res.StatusCode)
	res.writer.Write(res.Data)
}

func (res *Response) SendError(err error) *Response {
	res.SetRawData(res.encoder.EncodeError(err))
	res.SetStatus(http.StatusInternalServerError)
	res.Send()

	return res
}

// DoneWriter is a http.ResponseWriter which tracks its write state.
type DoneWriter struct {
	http.ResponseWriter
	Done bool
}

// WriteHeader wraps the underlying WriteHeader method.
func (w *DoneWriter) WriteHeader(status int) {
	w.Done = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *DoneWriter) Write(b []byte) (int, error) {
	w.Done = true
	return w.ResponseWriter.Write(b)
}
