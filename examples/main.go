package main

import (
	"fmt"
	"net/http"

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
		ctx.Bytes(200, whiskey2.MimeTypeJSON, []byte("Hello, World!"))
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

	whiskey.Run(whiskey2.RunOpts{
		Port: 8080,
	})
}
