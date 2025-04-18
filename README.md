# WHISKEY
A minimalist implementation of a HTTP webserver library on top of TCP

## Why another webserver?
We use webservers like `express` in Node or `gin` in Go to build web applications every day, but very rarely understand what it takes to build something like that. 

This project is intended to be a learning experience on

1. Implementation of HTTP, HTTPS and WebSockets
2. How path routing works
3. How different protocols of HTTP like HTTP/1.1 vs HTTP/2 works

## Planned Protocol Support
1. HTTP/1.1
2. HTTP/2
3. WebSockets
4. HTTPS

## Why Whiskey?
Gin is a very popular Go web framework and I loved the name, so I thought why not name this project after a drink.

## How does it work?

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sriramr98/whiskey"
)

type RequestBody struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type RequestQuery struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

type ResponseBody struct {
	Message string `json:"message"`
}

// Execution will not continue to controller if error is not nil
// Any custom params added to context will be propogated to the next middleware / controller
func AuthMiddleware(ctx whiskey.Context) error {
	header := ctx.GetHeader("Authorization")
	if !strings.ContainsPrefix(header, "Bearer") {
		return errors.New("Invalid JWT Token")
	}

	ctx.Set("userId", "dummyId")
	return nil
}

func main() {
	server := whiskey.New()

	server.Use(whiskey.CorsMiddleware)
	server.Use(AuthMiddleware)

	server.GET("/hello/:name", func(ctx whiskey.Context) error {
		var queryParams RequestQuery
		if err := ctx.BindQuery(&queryParams); err != nil {
			// Unable to bind query params to struct
			ctx.StatusCode(http.StatusInternalServerError)
			return errors.New("Unable to bind query params to struct")
		}

		ctx.String(http.StatusOK, fmt.Sprintf("Hello %s", name))
        return nil
	})

	server.POST("/hello/:name", func(ctx whiskey.Context) error {
		var body RequestBody
		if err := ctx.BindBody(&body); err != nil {
			ctx.StatusCode(http.StatusInternalServerError)
			return errors.New("invalid request body")
		}

		resp := ResponseBody{
			Message: "Created successfully"
		}

		ctx.JSON(http.StatusCreated, resp)
	})

	// Routes can also be declared like this and given to whiskey for handling to make configuration easy
	routes := []whiskey2.Route{
			{
				Path:    "/hello/:name",
				Method:  "GET",
				Handler: getHandler,
			},
			{
				Path:   "/hello/:name",
				Method: "POST",
				Handler: postHandler,
			}
		}
	}
	server.Handle(routes)

	server.Run(whiskey.RunOpts{
		Port: 8080
	})

```

## Whiskey Context
Whiskey context is the central piece of controllers and middlewares. It contains all info about request, response and custom keys passed around by middlewares.

Here is the interface definition

```go
// Context is an interface that represents the context of a single request. It contains all information regarding that request and is propagated through all middlewares
type Context interface {
	StatusCode(statusCode int)			 // This sets the status code for the request to be sent with the response

	BindBody(body interface{}) error     // The body will be a struct and the function will add the body parameters to the struct. The request body is expected to be a valid JSON
	BindQuery(query interface{}) error   // The query will be a struct and the function will add the query parameters to the struct.
	BindPath(path interface{}) error     // The path will be a struct and the function will add the path parameters to the struct.
	BindHeader(header interface{}) error // The header will be a struct and the function will add the header parameters to the struct.

	JSON(statusCode int, data interface{}) error // The function will convert the data to JSON and send it as a response
	String(statusCode int, data string) error    // The function will send the data as a string response
	Html(statusCode int, data string) error      // The function will send the data as a HTML response
	Bytes(statusCode int, contentType string, data []byte) error // The function will send the data as a byte array response

	GetQueryParam(key string) (string, bool) // The function will return the query parameter value for the given key. The boolean denotes whether the query param exists
	GetPathParam(key string) (string, bool)  // The function will return the path parameter value for the given key. The boolean denotes whether the path param exists
	GetHeader(key string) (string, bool)     // The function will return the header value for the given key. The boolean denotes whether the header exists

	GetQueryParams() map[string]string // The function will return all the query parameters
	GetPathParams() map[string]string  // The function will return all the path parameters
	GetHeaders() map[string]string     // The function will return all the headers
}
```