package whiskey

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// HTTP 1.1 connection handler
func (w *Whiskey) handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			w.accessLogger.Println("Error closing connection:", err)
		}
	}(conn)

	conn.SetReadDeadline(time.Now().Add(w.config.ReadTimeout))

	// Read the request
	req, err := readRequest(conn)
	if err != nil {
		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
			w.accessLogger.Println("Connection closed...")
			return
		}

		w.accessLogger.Println("Error reading request:", err)
		errResp := &HttpResponse{
			statusCode: http.StatusBadRequest,
			body:       []byte("invalid request"),
			headers:    make(map[string]string),
		}
		w.writeResponse(errResp, conn)

		return
	}

	w.accessLogger.Printf("Request on path %s", req.path)

	config, validRouteConfig := w.router.getConfig(req.path, req.method)
	handlers := config.handlers
	if !validRouteConfig {
		w.accessLogger.Printf("Handler not found for path %s\n", req.path)
		globalHandler, ok := w.router.getGlobalRequestHandler()
		if !ok {
			w.accessLogger.Println("No handler found for path:", req.path)
			w.handleNotFound(conn)
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

	var handlerErr error
	for _, handler := range handlers {
		handlerErr = handler(&ctx)
		if handlerErr != nil {
			// We break here and let the global error handler take care of handling the error down the line
			break
		}
	}

	if handlerErr != nil {
		if err := w.router.errorHandler(handlerErr, ctx); err != nil {
			w.errorLogger.Println("Error in error handler:", err)
			// Error handler failed, send a generic error response
			ctx.String(http.StatusInternalServerError, "Internal Server Error")
		}
	}

	// Default response type of text/plain unless overriden in the handler
	resp.SetHeader(HeaderContentType, fmt.Sprintf("%s; charset=utf-8", MimeTypeText))
	resp.SetHeader(HeaderConnection, "close") // Even if the client wants us to keep the connection alive, we close it

	w.writeResponse(resp, conn)
}

func (w *Whiskey) handleNotFound(conn net.Conn) {
	resp := &HttpResponse{
		headers:    make(map[string]string),
		statusCode: http.StatusNotFound,
		body:       []byte("Path route not found"),
	}
	w.writeResponse(resp, conn)
}
