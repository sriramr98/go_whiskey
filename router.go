package whiskey

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

// AddHandler adds a handler for a given path and method
func (r *router) addHandler(path string, method string, handler HttpHandler) {
	config := routeConfig{handler: handler}
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
