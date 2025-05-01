package whiskey

import (
	"bytes"
	"maps"
)

type RunOpts struct {
	Port int
	Addr string
}

type HttpRequest struct {
	method      string
	path        string
	body        []byte
	headers     map[string]string
	queryParams map[string]string
	pathParams  map[string]string
}

func (h HttpRequest) Equal(other HttpRequest) bool {
	return h.method == other.method &&
		h.path == other.path &&
		bytes.Equal(h.body, other.body) &&
		maps.Equal(h.headers, other.headers) &&
		maps.Equal(h.queryParams, other.queryParams) &&
		maps.Equal(h.pathParams, other.pathParams)
}

type HttpResponse struct {
	statusCode int
	body       []byte
	headers    map[string]string
}

func (resp *HttpResponse) SetHeader(key string, value string) {
	if resp.headers == nil {
		resp.headers = make(map[string]string)
	}
	resp.headers[key] = value
}

func (resp *HttpResponse) Send(body []byte) {
	resp.body = body
}

type HttpHandler func(ctx Context) error

type HttpErrorHandler func(err error, ctx Context) error
