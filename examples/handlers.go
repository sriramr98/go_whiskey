package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sriramr98/whiskey"
)

type BodyType struct {
	Key1 string `json:"key1"`
	Key2 string `json:"key2"`
}

func HelloGetHandler(ctx whiskey.Context) error {
	fmt.Println("Inside GET handler")
	ctx.SetHeader("CustomHeader", "CustomValue")
	ctx.String(200, "Hello, World!")
	return nil
}

func HelloPostHandler(ctx whiskey.Context) error {
	fmt.Println("Inside POST handler")

	var body BodyType
	err := ctx.BindBody(&body)
	if err != nil {
		return err
	}

	fmt.Printf("Received body: %+v\n", body)
	ctx.Json(http.StatusOK, body)

	return nil
}

func ImageHandler(ctx whiskey.Context) error {
	imgPath, ok := ctx.GetQueryParam("imgPath")
	if !ok {
		return whiskey.NewHTTPErrorWithMessage(http.StatusBadRequest, "QueryParam imgPath not found", whiskey.BodyTypeJSON)
	}

	data, err := os.ReadFile(imgPath)
	if err != nil {
		return whiskey.NewHTTPErrorWithMessage(http.StatusBadRequest, "Image not found at Path", whiskey.BodyTypeJSON)
	}

	ctx.Bytes(http.StatusOK, whiskey.MimeTypePNG, data)
	return nil
}

func ApiPathHandler(ctx whiskey.Context) error {
	type path struct {
		Id string `json:"id"`
	}

	var pathParams path
	err := ctx.BindPath(&pathParams)
	if err != nil {
		return whiskey.NewHTTPErrorWithMessage(http.StatusBadRequest, "Path Param id not found", whiskey.BodyTypeJSON)
	}

	fmt.Println("Got ID " + pathParams.Id)
	ctx.String(http.StatusOK, pathParams.Id)
	return nil
}

func ProtectedRouteHandler(ctx whiskey.Context) error {
	userId, hasUserId := ctx.GetString("userId")
	if !hasUserId {
		return whiskey.NewHttpError(http.StatusUnauthorized, whiskey.BodyTypeJSON)
	}

	fmt.Println("Got UserID " + userId)
	ctx.Json(http.StatusOK, whiskey.Json{
		"userId": userId,
	})
	return nil
}
