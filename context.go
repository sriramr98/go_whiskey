package whiskey

import (
	"encoding/json"
	"net/http"
)

// Context is an interface that represents the context of a single request. It contains all information regarding that request and is propagated through all middlewares
type Context interface {
	BindBody(body any) error     // The body will be a struct and the function will add the body parameters to the struct. The request body is expected to be a valid JSON
	BindQuery(query any) error   // The query will be a struct and the function will add the query parameters to the struct.
	BindPath(path any) error     // The path will be a struct and the function will add the path parameters to the struct.
	BindHeader(header any) error // The header will be a struct and the function will add the header parameters to the struct.

	Json(statusCode int, data any) error                         // The function will convert the data to JSON and send it as a response
	String(statusCode int, data string) error                    // The function will send the data as a string response
	Html(statusCode int, data string) error                      // The function will send the data as a HTML response
	Bytes(statusCode int, contentType string, data []byte) error // The function will send the data as a byte array response. Since we won't know what content type to set, you need to pass in the appropriate type, If content type is empty, Content-Type header won't be sent

	GetQueryParam(key string) (string, bool) // The function will return the query parameter value for the given key. The boolean denotes whether the query param exists
	GetPathParam(key string) (string, bool)  // The function will return the path parameter value for the given key. The boolean denotes whether the path param exists
	GetHeader(key string) (string, bool)     // The function will return the header value for the given key. The boolean denotes whether the header exists

	GetQueryParams() map[string]string // The function will return all the query parameters
	GetPathParams() map[string]string  // The function will return all the path parameters
	GetHeaders() map[string]string     // The function will return all the headers

	SetHeader(key string, value string) // The function will set the header for the response.

	URL() string    // The function will return the current path for which the request is being processed.
	Method() string // The function will return the current HTTP method for the request
}

type RequestContext struct {
	*DataStore // This is used as temporary storage for the request. It is not persisted across requests, but persisted across middlewares in a single request
	request    HttpRequest
	response   *HttpResponse
}

func (r RequestContext) BindBody(body any) error {
	return json.Unmarshal(r.request.body, body)
}

func (r RequestContext) BindQuery(query any) error {
	return mapToStruct(r.request.queryParams, query)
}

func (r RequestContext) BindPath(path any) error {
	return mapToStruct(r.request.pathParams, path)
}

func (r RequestContext) BindHeader(header any) error {
	return mapToStruct(r.request.headers, header)
}

func (r RequestContext) Json(statusCode int, data any) error {
	r.response.SetHeader(HeaderContentType, MimeTypeJSON)
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	r.response.statusCode = statusCode

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	r.response.Send(b)
	return nil
}

func (r RequestContext) String(statusCode int, data string) error {
	r.response.SetHeader(HeaderContentType, MimeTypeText)
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	r.response.statusCode = statusCode

	b := []byte(data)
	r.response.Send(b)
	return nil
}

func (r RequestContext) Html(statusCode int, data string) error {
	r.response.SetHeader(HeaderContentType, MimeTypeHTML)
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	r.response.statusCode = statusCode

	b := []byte(data)
	r.response.Send(b)
	return nil
}

func (r RequestContext) GetQueryParam(key string) (string, bool) {
	if value, ok := r.request.queryParams[key]; ok {
		return value, true
	}
	return "", false
}

func (r RequestContext) GetPathParam(key string) (string, bool) {
	if value, ok := r.request.pathParams[key]; ok {
		return value, true
	}
	return "", false
}

func (r RequestContext) GetHeader(key string) (string, bool) {
	if value, ok := r.request.headers[key]; ok {
		return value, true
	}
	return "", false
}

func (r RequestContext) GetQueryParams() map[string]string {
	return r.request.queryParams
}

func (r RequestContext) GetPathParams() map[string]string {
	return r.request.pathParams
}

func (r RequestContext) GetHeaders() map[string]string {
	return r.request.headers
}

func (r RequestContext) Bytes(statusCode int, contentType string, data []byte) error {
	r.response.SetHeader(HeaderContentType, contentType)
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	r.response.statusCode = statusCode

	r.response.Send(data)
	return nil
}

func (r RequestContext) SetHeader(key string, value string) {
	r.response.SetHeader(key, value)
}

func (r RequestContext) URL() string {
	return r.request.path
}

func (r RequestContext) Method() string {
	return r.request.method
}
