package Route

import (
	"golang.org/x/exp/maps"
)

type RouteParameterMatch struct {
	Indexes  []int
	RawParam string
	Param    string
	Matched  bool
}

func CreateRouteParameterMatch(section string) RouteParameterMatch {
	return RouteParameterMatch{
		Indexes:  nil,
		RawParam: section,
		Param:    "",
		Matched:  false,
	}
}

type RouteParameters[T string] struct {
	data map[string]T
}

func CreateRouteParameters(data map[string]string) *RouteParameters[string] {
	return &RouteParameters[string]{
		data: data,
	}
}

func (p *RouteParameters[T]) Has(key string, matchVal ...T) bool {
	val, ok := p.data[key]

	if len(matchVal) >= 1 && ok {
		return matchVal[0] == val
	}

	return ok
}

func (p *RouteParameters[T]) Get(key string, defaultVal ...T) T {
	val, ok := p.data[key]

	if ok {
		return val
	}

	if len(defaultVal) >= 1 {
		return defaultVal[0]
	}

	return ""
}

func (p *RouteParameters[T]) Keys() []string {
	return maps.Keys(p.data)
}

func (p *RouteParameters[T]) Values() []T {
	return maps.Values(p.data)
}

func (p *RouteParameters[T]) Clear() {
	maps.Clear(p.data)
}
