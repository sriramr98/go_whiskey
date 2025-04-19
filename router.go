package whiskey

// Router handles figuring out which handler to be called for a given request
type router struct {
	handlers             map[string]map[string]HttpHandler // For every path, for every method, a handler
	globalRequestHandler HttpHandler                       // This gets called if no path is matched
	globalHandlerSet     bool                              // Indicates if a global handler has been set
	errorHandler         HttpErrorHandler                  // This gets called if an error occurs
}

// NewRouter creates a new router instance
func newRouter() *router {
	return &router{
		handlers:     make(map[string]map[string]HttpHandler),
		errorHandler: defaultErrorHandler,
	}
}

// AddHandler adds a handler for a given path and method
func (r *router) addHandler(path string, method string, handler HttpHandler) {
	if r.handlers[path] == nil {
		r.handlers[path] = make(map[string]HttpHandler)
	}
	r.handlers[path][method] = handler
}

func (r *router) getHandler(path string, method string) (HttpHandler, bool) {
	if r.handlers[path] == nil {
		return nil, false
	}
	handler, ok := r.handlers[path][method]
	return handler, ok
}

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
