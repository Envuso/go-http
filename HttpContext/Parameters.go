package HttpContext

import (
	"golang.org/x/exp/maps"
)

type Parameters[T string] struct {
	data map[string]T
}

func ParametersFromValues(data any) *Parameters[string] {
	var parameters Parameters[string] = Parameters[string]{data: make(map[string]string)}
	// for _, param := range data {
	// 	parameters.data[param.Name] = param.Value
	// }
	return &parameters
}

func (p *Parameters[T]) Set(key string, value T) {
	p.data[key] = value
}

func (p *Parameters[T]) Has(key string) bool {
	_, ok := p.data[key]
	return ok
}

func (p *Parameters[T]) Get(key string, defaultVal ...T) T {
	val, ok := p.data[key]

	if ok {
		return val
	}

	if len(defaultVal) >= 1 {
		return defaultVal[0]
	}

	return ""
}

func (p *Parameters[T]) Keys() []string {
	return maps.Keys(p.data)
}

func (p *Parameters[T]) Values() []T {
	return maps.Values(p.data)
}

func (p *Parameters[T]) Clear() {
	maps.Clear(p.data)
}

type QueryParameters struct {
	data map[string][]string
}

func QueryParametersFrom(data map[string][]string) *QueryParameters {
	return &QueryParameters{data: data}
}

func (p *QueryParameters) Has(key string) bool {
	_, ok := p.data[key]
	return ok
}

func (p *QueryParameters) Get(key string, defaultVal ...string) string {
	val, ok := p.data[key]
	def := ""
	if len(defaultVal) >= 1 {
		def = defaultVal[0]
	}

	if !ok {
		return def
	}

	if len(val) == 0 {
		return defaultVal[0]
	}

	return val[0]
}

func (p *QueryParameters) GetArray(key string, defaultVal ...string) []string {
	val, ok := p.data[key]
	def := []string{}
	if len(defaultVal) >= 1 {
		def = defaultVal
	}
	if !ok {
		return def
	}
	if len(val) == 0 {
		return def
	}

	return val
}

func (p *QueryParameters) Keys() []string {
	return maps.Keys(p.data)
}

func (p *QueryParameters) Values() [][]string {
	return maps.Values(p.data)
}

func (p *QueryParameters) Clear() {
	maps.Clear(p.data)
}
