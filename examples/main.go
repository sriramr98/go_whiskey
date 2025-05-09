package main

import (
	"net/http"

	whiskey2 "github.com/sriramr98/whiskey"
)

func main() {
	whiskey := whiskey2.New()
	whiskey.GET("/hello", HelloGetHandler)

	whiskey.POST("/hello", HelloPostHandler)

	whiskey.GET("/image", ImageHandler)

	whiskey.GET("/api/{id}", ApiPathHandler)

	whiskey.GET("/protected", AuthMiddleware, ProtectedRouteHandler)

	// If not paths match, this will get called
	whiskey.GlobalRequestHandler(func(ctx whiskey2.Context) error {
		return whiskey2.NewHttpError(http.StatusNotFound, whiskey2.BodyTypeJSON)
	})

	whiskey.Run(whiskey2.RunOpts{
		Port: 8080,
	})
}
