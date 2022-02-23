package Routing

import (
	"regexp"
	"strings"

	"gohttp/Routing/Route"
)

var RouteParamRegex = regexp.MustCompile(`{(!)?(\w+)}`)

type RouteRegistrar struct {
	children map[string]*RouteRegistrar

	httpMethod string
	route      *Route.Route
}

func NewRouteRegistrar(httpMethod string) *RouteRegistrar {
	return &RouteRegistrar{
		httpMethod: httpMethod,
		children:   make(map[string]*RouteRegistrar),
	}
}

// getRouteParam process the section(a section is a part of a route path("/route/name/something") split by "/")
func (r *RouteRegistrar) getRouteParam(section string) Route.RouteParameterMatch {
	match := Route.CreateRouteParameterMatch(section)
	if len(section) <= 2 {
		return match
	}

	paramIndexes := RouteParamRegex.FindStringIndex(section)
	if len(paramIndexes) <= 0 {
		return match
	}

	matches := RouteParamRegex.FindStringSubmatch(section)
	if len(matches) != 3 {
		return match
	}

	match.Indexes = paramIndexes
	match.Param = matches[len(matches)-1]
	match.Matched = true

	return match
}

// addRoute Create a "trie" for the provided route definition
func (r *RouteRegistrar) addRoute(route *Route.Route) *RouteRegistrar {
	current := r

	trimmed := strings.TrimPrefix(route.Path(), "/")
	slice := strings.Split(trimmed, "/")

	for _, k := range slice {
		param := r.getRouteParam(k)
		if param.Matched {
			k = "*"
		}
		next, ok := current.children[k]
		if !ok {
			next = NewRouteRegistrar(r.httpMethod)
			next.route = Route.CreateRouteRegistration(r.httpMethod, route, param, k)
			current.children[k] = next
		}
		current = next
	}

	return current
}

// Find Match the incoming request path to a route using trie tree's
func (r *RouteRegistrar) Find(path string) *Route.RouteMatch {
	params := make(map[string]string)
	current := r

	trimmed := strings.TrimPrefix(path, "/")
	slice := strings.Split(trimmed, "/")

	for _, k := range slice {
		next := current.hasRoutePath(k)
		if next == nil {
			return nil
		}

		current = next

		// if the node has a param add it to params map.
		if current.route.Param().Matched {
			params[current.route.Param().Param] = k
		}
	}

	return Route.CreateRouteMatch(*current.route, params)
}

// hasRoutePath check the section to see if it has children/dynamic route path
func (r *RouteRegistrar) hasRoutePath(k string) *RouteRegistrar {
	next, ok := r.children[k]
	if !ok {
		next, ok = r.children["*"]
		if !ok {
			return nil
		}
	}

	return next
}

func (r *RouteRegistrar) GetChildren() map[string]*RouteRegistrar {
	return r.children
}
func (r *RouteRegistrar) CurrentRoute() *Route.Route {
	return r.route
}
