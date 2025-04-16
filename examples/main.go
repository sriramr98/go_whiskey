package main

import (
	"encoding/json"
	"fmt"
	whiskey2 "github.com/sriramr98/whiskey"
)

type BodyType struct {
	Key1 string `json:"key1"`
	Key2 string `json:"key2"`
}

func main() {
	whiskey := whiskey2.New()

	whiskey.GET("/hello", func(req whiskey2.HttpRequest, resp *whiskey2.HttpResponse) {
		fmt.Println("Inside GET handler")
		fmt.Println(req)

		resp.SetHeader("Content-Type", "application/json")
		resp.Send([]byte("{\"message\": \"Hello, World!\"}"))
	})

	whiskey.POST("/hello", func(req whiskey2.HttpRequest, resp *whiskey2.HttpResponse) {
		fmt.Println("Inside POST handler")
		fmt.Println(req)

		var body BodyType
		err := json.Unmarshal(req.Body, &body)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error unmarshalling JSON:", err)
			resp.SetHeader("Content-Type", "application/json")
			resp.Send([]byte("{\"error\": \"Invalid JSON\"}"))
			return
		}

		fmt.Printf("Received body: %+v\n", body)

		response := map[string]string{
			"key1": body.Key1,
			"key2": body.Key2,
		}
		respBytes, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			resp.SetHeader("Content-Type", "application/json")
			resp.Send([]byte("{\"error\": \"Internal Server Error\"}"))
			return
		}

		resp.SetHeader("Content-Type", "application/json")
		resp.Send(respBytes)
	})

	whiskey.Run(whiskey2.RunOpts{
		Port: 8080,
	})
}
