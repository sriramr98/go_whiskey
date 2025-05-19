package whiskey

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

// If we get these content types, the response body will have to be stringified and sent
var stringContentTypes = []string{
	MimeTypeHTML,
	MimeTypeText,
	MimeTypeJSON,
	MimeTypeXML,
}

type Whiskey struct {
	router       *router
	config       ServerConfig
	accessLogger *log.Logger
	errorLogger  *log.Logger
}

// Default settings for the Whiskey engine.
var (
	PORT = 8080
	ADDR = "0.0.0.0" // Listens on all IP ranges by default
)

var defaultConfig = ServerConfig{
	Port:           8080,
	Addr:           "0.0.0.0",
	MaxConcurrency: 1000,
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
}

// New creates a new Whiskey engine instance with default settings.
func New() Whiskey {
	return Whiskey{
		router:       newRouter(),
		config:       defaultConfig,
		errorLogger:  log.New(log.Writer(), "ERROR: ", log.LstdFlags),
		accessLogger: log.New(log.Writer(), "ACCESS: ", log.LstdFlags),
	}
}

func (w *Whiskey) WithConfig(config ServerConfig) *Whiskey {
	w.config = config
	return w
}

func (w *Whiskey) WithAccessLogger(logger *log.Logger) *Whiskey {
	w.accessLogger = logger
	return w
}

func (w *Whiskey) WithErrorLogger(logger *log.Logger) *Whiskey {
	w.errorLogger = logger
	return w
}

// GET registers a handler for the given path with the HTTP GET method.
func (w *Whiskey) GET(path string, handlers ...HttpHandler) {
	w.router.addHandler(path, http.MethodGet, handlers)
}

// POST registers a handler for the given path with the HTTP POST method.
func (w *Whiskey) POST(path string, handlers ...HttpHandler) {
	w.router.addHandler(path, http.MethodPost, handlers)
}

// PUT registers a handler for the given path with the HTTP PUT method.
func (w *Whiskey) PUT(path string, handlers ...HttpHandler) {
	w.router.addHandler(path, http.MethodPut, handlers)
}

// DELETE registers a handler for the given path with the HTTP DELETE method.
func (w *Whiskey) DELETE(path string, handlers ...HttpHandler) {
	w.router.addHandler(path, http.MethodDelete, handlers)
}

// PATCH registers a handler for the given path with the HTTP PATCH method.
func (w *Whiskey) PATCH(path string, handlers ...HttpHandler) {
	w.router.addHandler(path, http.MethodPatch, handlers)
}

func (w *Whiskey) GlobalErrorHandler(handler HttpErrorHandler) {
	w.router.setErrorHandler(handler)
}

// GlobalRequestHandler handles any requests for which a handler isn't mapped. Ideal for returning custom 404 not found responses
func (w *Whiskey) GlobalRequestHandler(handler HttpHandler) {
	w.router.setGlobalRequestHandler(handler)
}

// ConfigRoutes provides an easily utility to configure routes with a simple data structure
func (w *Whiskey) ConfigRoutes(routes []Route) {
	for _, route := range routes {
		w.router.addHandler(route.Path, route.Method, route.Handlers)
	}
}

// Run starts the HTTP server and blocks until it is stopped
func (w *Whiskey) Run() {
	// Start the HTTP server
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", w.config.Addr, w.config.Port))
	if err != nil {
		w.errorLogger.Fatal("Error starting server:", err)
	}
	defer ln.Close()

	w.accessLogger.Printf("Starting server on %s:%d\n", w.config.Addr, w.config.Port)

	for {
		// This blocks until a connection is accepted
		conn, err := ln.Accept()
		if err != nil {
			w.errorLogger.Fatal("Error accepting connection:", err)
		}

		go w.handleConnection(conn)
	}
}
