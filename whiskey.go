package whiskey

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

// If we get these content types, the response body will have to be stringified and sent
var stringContentTypes = []string{
	MimeTypeHTML,
	MimeTypeText,
	MimeTypeJSON,
	MimeTypeXML,
}

type Whiskey struct {
	router *router
}

type Route struct {
	Path     string
	Method   string
	Handlers []HttpHandler
}

// Default settings for the Whiskey engine.
var (
	PORT = 8080
	ADDR = "0.0.0.0" // Listens on all IP ranges by default
)

// New creates a new Whiskey engine instance with default settings.
func New() Whiskey {
	return Whiskey{
		router: newRouter(),
	}
}

// GET registers a handler for the given path with the HTTP GET method.
func (w Whiskey) GET(path string, handlers ...HttpHandler) {
	w.router.addHandler(path, http.MethodGet, handlers)
}

// POST registers a handler for the given path with the HTTP POST method.
func (w Whiskey) POST(path string, handlers ...HttpHandler) {
	w.router.addHandler(path, http.MethodPost, handlers)
}

// PUT registers a handler for the given path with the HTTP PUT method.
func (w Whiskey) PUT(path string, handlers ...HttpHandler) {
	w.router.addHandler(path, http.MethodPut, handlers)
}

// DELETE registers a handler for the given path with the HTTP DELETE method.
func (w Whiskey) DELETE(path string, handlers ...HttpHandler) {
	w.router.addHandler(path, http.MethodDelete, handlers)
}

// PATCH registers a handler for the given path with the HTTP PATCH method.
func (w Whiskey) PATCH(path string, handlers ...HttpHandler) {
	w.router.addHandler(path, http.MethodPatch, handlers)
}

func (w Whiskey) GlobalErrorHandler(handler HttpErrorHandler) {
	w.router.setErrorHandler(handler)
}

// GlobalRequestHandler handles any requests for which a handler isn't mapped. Ideal for returning custom 404 not found responses
func (w Whiskey) GlobalRequestHandler(handler HttpHandler) {
	w.router.setGlobalRequestHandler(handler)
}

// ConfigRoutes provides an easily utility to configure routes with a simple data structure
func (w Whiskey) ConfigRoutes(routes []Route) {
	for _, route := range routes {
		w.router.addHandler(route.Path, route.Method, route.Handlers)
	}
}

// Run starts the HTTP server and blocks until it is stopped
func (w Whiskey) Run(opts RunOpts) {
	if opts.Port == 0 {
		opts.Port = PORT
	}
	if opts.Addr == "" {
		opts.Addr = ADDR
	}

	// Start the HTTP server
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", opts.Addr, opts.Port))
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer ln.Close()

	log.Printf("Starting server on %s:%d\n", opts.Addr, opts.Port)

	for {
		// This blocks until a connection is accepted
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Error accepting connection:", err)
		}

		go w.handleConnection(conn)
	}
}

// HTTP 1.1 connection handler
func (w Whiskey) handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("Error closing connection:", err)
		}
	}(conn)

	// Read the request
	req, err := readRequest(conn)
	if err != nil {
		log.Println("Error reading request:", err)
		return
	}

	config, validRouteConfig := w.router.getConfig(req.path, req.method)
	handlers := config.handlers
	if !validRouteConfig {
		log.Printf("Handler not found for path %s\n", req.path)
		globalHandler, ok := w.router.getGlobalRequestHandler()
		if !ok {
			resp := &HttpResponse{
				headers:    make(map[string]string),
				statusCode: http.StatusNotFound,
				body:       []byte("Path route not found"),
			}
			log.Println("No handler found for path:", req.path)
			writeResponse(resp, conn)
			return
		}
		handlers = []HttpHandler{globalHandler}
	} else {
		req.pathParams = config.pathParams
	}

	resp := &HttpResponse{
		headers: make(map[string]string),
	}
	ctx := RequestContext{
		DataStore: NewDataStore(),
		request:   req,
		response:  resp,
	}

	var handleError error
	if validRouteConfig {
		for _, handler := range handlers {
			handleError = handler(&ctx)
			if handleError != nil {
				// We break here and let the global error handler take care of handling the error down the line
				break
			}
		}
	}

	if handleError != nil {
		if err := w.router.errorHandler(handleError, ctx); err != nil {
			log.Println("Error in error handler:", err)
			// Error handler failed, send a generic error response
			ctx.String(http.StatusInternalServerError, "Internal Server Error")
		}
	}

	// Default response type of text/plain unless overriden in the handler
	resp.SetHeader(HeaderContentType, fmt.Sprintf("%s; charset=utf-8", MimeTypeText))
	resp.SetHeader(HeaderConnection, "close") // Even if the client wants us to keep the connection alive, we close it

	writeResponse(resp, conn)
}
