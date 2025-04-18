package whiskey

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"slices"
	"strings"
)

type Whiskey struct {
	router *router
}

// Default settings for the Whiskey engine.
var (
	PORT = 8080
	ADDR = "0.0.0.0" //Listens on all IP ranges by default
)

// New creates a new Whiskey engine instance with default settings.
func New() Whiskey {
	return Whiskey{
		router: newRouter(),
	}
}

// GET registers a handler for the given path with the HTTP GET method.
func (w Whiskey) GET(path string, handler HttpHandler) {
	w.router.addHandler(path, http.MethodGet, handler)
}

// POST registers a handler for the given path with the HTTP POST method.
func (w Whiskey) POST(path string, handler HttpHandler) {
	w.router.addHandler(path, http.MethodPost, handler)
}

// PUT registers a handler for the given path with the HTTP PUT method.
func (w Whiskey) PUT(path string, handler HttpHandler) {
	w.router.addHandler(path, http.MethodPut, handler)
}

// DELETE registers a handler for the given path with the HTTP DELETE method.
func (w Whiskey) DELETE(path string, handler HttpHandler) {
	w.router.addHandler(path, http.MethodDelete, handler)
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
	req, err := w.readRequest(conn)
	if err != nil {
		//TODO: Send HTTP response back
		log.Println("Error reading request:", err)
		return
	}

	handler, ok := w.router.getHandler(req.path, req.method)
	if !ok {
		handler, ok = w.router.getGlobalRequestHandler()
		if !ok {
			//TODO: Send 404 response
			log.Println("No handler found for path:", req.path)
			return
		}
	}

	resp := &HttpResponse{
		headers: make(map[string]string),
	}
	// Default response type of text/plain unless overriden in the handler
	resp.SetHeader(HeaderContentType, fmt.Sprintf("%s; charset=utf-8", MimeTypeText))
	// Call the handler
	handler(RequestContext{
		request:  req,
		response: resp,
	})

	w.writeResponse(resp, conn)
}

func (w Whiskey) readRequest(conn net.Conn) (HttpRequest, error) {
	tmp := make([]byte, 1024)

	// Size is 0 since we don't know how much total data we will read
	data := make([]byte, 0)
	length := 0

	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Read Error", err)
			}
			break
		}

		data = append(data, tmp[:n]...)
		length += n

		if n < 1024 {
			break
		}
	}

	lines := strings.Split(string(data), "\r\n")

	// Parse the request line
	return parseRequest(lines)
}

func (w Whiskey) writeResponse(resp *HttpResponse, conn net.Conn) {
	var responseLines []string

	if resp.statusCode == 0 {
		resp.statusCode = http.StatusOK
	}
	// We only support HTTP/1.1
	responseLines = append(responseLines, fmt.Sprintf("HTTP/1.1 %d OK", resp.statusCode))

	contentType, ok := resp.headers[HeaderContentType]
	if !ok {
		contentType = "text/plain; charset=utf-8"
	}
	contentLength := len(resp.body)
	resp.headers["Content-Length"] = fmt.Sprintf("%d", contentLength)

	for key, value := range resp.headers {
		responseLines = append(responseLines, fmt.Sprintf("%s: %s", key, value))
	}

	var body string
	if contentType == "text/plain" || contentType == "application/json" {
		body = string(resp.body)
	} else {
		var builder strings.Builder
		for _, b := range resp.body {
			builder.WriteString(fmt.Sprintf("%d", b))
		}
		body = builder.String()
	}

	response := fmt.Sprintf("%s\r\n\r\n%s", strings.Join(responseLines, "\r\n"), body)

	fmt.Println(response)
	_, err := conn.Write([]byte(response))
	if err != nil {
		return
	}
}

func parseRequest(requestData []string) (HttpRequest, error) {
	request := HttpRequest{
		headers: make(map[string]string),
	}

	if len(requestData) == 0 {
		return request, fmt.Errorf("invalid HTTP request")
	}

	// 1st line should contain the format {method} {path} HTTP/1.1
	protocolParts := strings.Split(strings.TrimSpace(requestData[0]), " ")
	if len(protocolParts) < 3 {
		return request, fmt.Errorf("invalid HTTP request format")
	}
	if !slices.Contains([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}, protocolParts[0]) {
		return request, fmt.Errorf("invalid HTTP method")
	}
	request.method = protocolParts[0]

	if !strings.HasPrefix(protocolParts[1], "/") {
		return request, fmt.Errorf("invalid HTTP path")
	}
	request.path = protocolParts[1]

	// We currently only support HTTP/1.1
	if protocolParts[2] != "HTTP/1.1" {
		return request, fmt.Errorf("invalid HTTP version")
	}

	// From the second line, it contains headers in the format {key}: {value} until we find an empty line

	bodyStartIdx := 1
	lastReadIdx := 1
	for idx, line := range requestData {
		if idx == 0 {
			continue
		}
		if line == "" {
			// An empty line indicates the end of headers
			bodyStartIdx = idx + 1 // The body starts after the empty line
			break
		}
		lastReadIdx = idx

		headerParts := strings.SplitN(line, ":", 2)
		if len(headerParts) < 2 {
			// ignore broken headers
			continue
		}

		key := strings.TrimSpace(headerParts[0])
		value := strings.TrimSpace(headerParts[1])

		request.headers[key] = value
	}

	// body isn't present, and we reached the end of the request
	// We do lastReadIdx+1 because the header loop starts from 1st index which is 0 inside the loop
	if lastReadIdx == len(requestData)-1 {
		return request, nil
	}

	// The body starts after the headers
	body := strings.Join(requestData[bodyStartIdx:], "")
	request.body = []byte(body)

	return request, nil
}
