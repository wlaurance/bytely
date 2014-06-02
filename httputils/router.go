package httputils

import (
	"net/http"
	"regexp"
	"strings"
)

// Router is a simple muxer which allows for regex captures in route
// patterns. It doesn't offer all of the bells and whistles of
// popular frameworks because I'm a simple man.
type Router struct {
	routes []*route
}

type route struct {
	re      *regexp.Regexp
	handler RouteHandlerFunc
}

type RouteHandlerFunc interface{}

// NewRouter allocates and returns a new Router.
func NewRouter() *Router {
	return &Router{nil}
}

// Route registers a handler function for the given pattern.
func (this *Router) Route(pattern string, handler RouteHandlerFunc) {
	// Bound the start and end of the regex to prevent multiple matches.
	re := regexp.MustCompile("^" + pattern + "$")
	newRoute := &route{re, handler}
	this.routes = append(this.routes, newRoute)
}

// ServeHTTP dispatches the request to the handler whose pattern
// most closely matches the request URL.
// It must determine if the handler means to capture any variables
// from the route pattern and, if so, pass those into the handler
// as an extra string parameter.
func (this *Router) ServeHTTP(response http.ResponseWriter, req *http.Request) {
	for _, route := range this.routes {
		cleanURI := strings.Split(req.RequestURI, "?")[0]
		match := route.re.FindStringSubmatch(cleanURI)
		if match != nil {
			switch handler := route.handler.(type) {
			case func(http.ResponseWriter, *http.Request):
				handler(response, req)
				return
			case func(http.ResponseWriter, *http.Request, string):
				handler(response, req, match[1])
				return
			default:
				break
			}
		}
	}
}
