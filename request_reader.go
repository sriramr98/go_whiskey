package whiskey

import (
	"io"
	"log"
)

func readRequest(reader io.Reader) (HttpRequest, error) {
	tmp := make([]byte, 1024)

	// Size is 0 since we don't know how much total data we will read
	data := make([]byte, 0)
	length := 0

	for {
		n, err := reader.Read(tmp)
		if err != nil {
			if err == io.EOF {
				log.Println("Connection Closed...")
			}
			log.Printf("Error reading from connection, err: %+v", err)
			break
		}

		data = append(data, tmp[:n]...)
		length += n

		if n < 1024 {
			break
		}
	}

	// Parse the request line
	return parseRequest(string(data))
}

func parseRequest(requestData string) (HttpRequest, error) {
	return HTTP_1_1_Parser(requestData)
}
