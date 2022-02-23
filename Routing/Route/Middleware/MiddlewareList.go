package Middleware

import "golang.org/x/exp/maps"

type MiddlewareList struct {
	data map[string]Middleware
}

func NewMiddlewareList() *MiddlewareList {
	return &MiddlewareList{data: make(map[string]Middleware)}
}

func (p *MiddlewareList) Data() map[string]Middleware {
	return p.data
}

func (p *MiddlewareList) SetData(data map[string]Middleware) {
	p.data = data
}

func (p *MiddlewareList) Set(key string, value Middleware) {
	p.data[key] = value
}

func (p *MiddlewareList) Has(key string) bool {
	_, ok := p.data[key]
	return ok
}

func (p *MiddlewareList) Get(key string, defaultVal ...Middleware) Middleware {
	val, ok := p.data[key]

	if ok {
		return val
	}

	if len(defaultVal) >= 1 {
		return defaultVal[0]
	}

	return nil
}
func (p *MiddlewareList) GetOk(key string, defaultVal ...Middleware) (Middleware, bool) {
	val, ok := p.data[key]

	if ok {
		return val, ok
	}

	if len(defaultVal) >= 1 {
		return defaultVal[0], ok
	}

	return nil, ok
}

func (p *MiddlewareList) Keys() []string {
	return maps.Keys(p.data)
}

func (p *MiddlewareList) Values() []Middleware {
	return maps.Values(p.data)
}

func (p *MiddlewareList) Clear() {
	maps.Clear(p.data)
}

func (p *MiddlewareList) MergeIn(middlewares *MiddlewareList) {
	for name, middleware := range middlewares.data {
		if !p.Has(name) {
			p.Set(name, middleware)
		}
	}
}

func (p *MiddlewareList) IsEmpty() bool {
	return p.Length() == 0
}

func (p *MiddlewareList) Length() int {
	return len(p.data)
}
