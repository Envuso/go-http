package HttpContext

import (
	"net/http"
	"net/textproto"

	"golang.org/x/exp/maps"
)

type HeaderContainer struct {
	data http.Header
}

func NewHeaderContainer() *HeaderContainer {
	return &HeaderContainer{data: make(http.Header)}
}

func (p *HeaderContainer) Data() http.Header {
	return p.data
}

func (p *HeaderContainer) SetData(data http.Header) *HeaderContainer {
	p.Clear()

	for key, header := range data {
		p.Set(key, header[0])
	}

	return p
}

func (p *HeaderContainer) Set(key, value string) {
	key = textproto.CanonicalMIMEHeaderKey(key)

	textproto.MIMEHeader(p.data).Set(key, value)
}

func (p *HeaderContainer) Has(key string, value ...string) bool {
	key = textproto.CanonicalMIMEHeaderKey(key)

	dat, ok := p.data[key]

	if ok && len(dat) == 0 {
		return false
	}

	if len(value) > 0 {
		return p.Get(key) == value[0]
	}

	return ok
}

func (p *HeaderContainer) Get(key string, defaultVal ...string) string {
	key = textproto.CanonicalMIMEHeaderKey(key)

	val := textproto.MIMEHeader(p.data).Get(key)
	if val != "" {
		return val
	}

	if len(defaultVal) >= 1 {
		return defaultVal[0]
	}

	return ""
}

func (p *HeaderContainer) Keys() []string {
	return maps.Keys(p.data)
}

func (p *HeaderContainer) Values() [][]string {
	return maps.Values(p.data)
}

func (p *HeaderContainer) Clear() {
	maps.Clear(p.data)
}

func (p *HeaderContainer) IsEmpty() bool {
	return p.Length() == 0
}

func (p *HeaderContainer) Length() int {
	return len(p.data)
}

func (p *HeaderContainer) SetFrom(headers map[string]string) {
	for key, header := range headers {
		if p.Has(key) {
			continue
		}

		p.Set(key, header)
	}
}
