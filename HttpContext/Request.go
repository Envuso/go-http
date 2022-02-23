package HttpContext

import (
	"net/http"
	"strconv"
)

type Request struct {
	reader    *http.Request
	requestId uint64
	headers   *HeaderContainer
	encoder   ContentEncoder
}

func NewRequest(reader *http.Request) *Request {
	req := &Request{
		reader: reader,
	}

	req.headers = NewHeaderContainer().SetData(reader.Header)

	return req
}

func (req *Request) Id() uint64 {
	return req.requestId << 32
}

func (req *Request) IdString() string {
	return strconv.Itoa(int(req.Id()))
}

func (req *Request) Headers() *HeaderContainer {
	return req.headers
}
