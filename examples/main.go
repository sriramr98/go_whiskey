package main

import (
	"fmt"
	"net/http"
	"os"

	whiskey2 "github.com/sriramr98/whiskey"
)

type BodyType struct {
	Key1 string `json:"key1"`
	Key2 string `json:"key2"`
}

func main() {
	whiskey := whiskey2.New()
	whiskey.GET("/hello", func(ctx whiskey2.Context) error {
		fmt.Println("Inside GET handler")
		ctx.SetHeader("CustomHeader", "CustomValue")
		ctx.String(200, "Hello, World!")
		return nil
	})

	whiskey.POST("/hello", func(ctx whiskey2.Context) error {
		fmt.Println("Inside POST handler")

		var body BodyType
		err := ctx.BindBody(&body)
		if err != nil {
			return err
		}

		fmt.Printf("Received body: %+v\n", body)
		ctx.Json(http.StatusOK, body)

		return nil
	})

	whiskey.GET("/image", func(ctx whiskey2.Context) error {
		imgPath, ok := ctx.GetQueryParam("imgPath")
		if !ok {
			return whiskey2.NewHTTPErrorWithMessage(http.StatusBadRequest, "QueryParam imgPath not found", whiskey2.BodyTypeJSON)
		}

		data, err := os.ReadFile(imgPath)
		if err != nil {
			return whiskey2.NewHTTPErrorWithMessage(http.StatusBadRequest, "Image not found at Path", whiskey2.BodyTypeJSON)
		}

		ctx.Bytes(http.StatusOK, whiskey2.MimeTypePNG, data)
		return nil
	})

	whiskey.GET("/api/{id}", func(ctx whiskey2.Context) error {
		type path struct {
			Id string `json:"id"`
		}

		var pathParams path
		err := ctx.BindPath(&pathParams)
		if err != nil {
			return whiskey2.NewHTTPErrorWithMessage(http.StatusBadRequest, "Path Param id not found", whiskey2.BodyTypeJSON)
		}

		fmt.Println("Got ID " + pathParams.Id)
		ctx.String(http.StatusOK, pathParams.Id)
		return nil
	})

	// If not paths match, this will get called
	whiskey.GlobalRequestHandler(func(ctx whiskey2.Context) error {
		return whiskey2.NewHttpError(http.StatusNotFound, whiskey2.BodyTypeJSON)
	})

	whiskey.Run(whiskey2.RunOpts{
		Port: 8080,
	})
}
