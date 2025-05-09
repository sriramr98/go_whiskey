package whiskey

import (
	"net/http"
	"slices"
)

var configurableHttpMethods []string = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
}

// Router handles figuring out which handler to be called for a given request
type router struct {
	routes               *routeTree
	globalRequestHandler HttpHandler      // This gets called if no path is matched
	globalHandlerSet     bool             // Indicates if a global handler has been set
	errorHandler         HttpErrorHandler // This gets called if an error occurs
}

// NewRouter creates a new router instance
func newRouter() *router {
	return &router{
		routes:       newRouteTree(),
		errorHandler: defaultErrorHandler,
	}
}

// AddHandler adds a set of handlers for a given path and method
func (r *router) addHandler(path string, method string, handlers []HttpHandler) {
	if !slices.Contains(configurableHttpMethods, method) {
		// Since route configuration happens before server is started, panic is fine
		panic("Invalid HTTP method " + method + " configured")
	}
	config := routeConfig{handlers: handlers}
	r.routes.insert(path, method, config)
}

func (r *router) getConfig(path string, method string) (routeConfig, bool) {
	config, ok := r.routes.getConfig(path, method)
	return config, ok
}

// setGlobalRequestHandler assigns the request handler that gets called if no paths in the server match the incoming path. It's a default request handler
func (r *router) setGlobalRequestHandler(handler HttpHandler) {
	r.globalHandlerSet = true
	r.globalRequestHandler = handler
}

func (r *router) getGlobalRequestHandler() (HttpHandler, bool) {
	return r.globalRequestHandler, r.globalHandlerSet
}

func (r *router) setErrorHandler(handler HttpErrorHandler) {
	r.errorHandler = handler
}
