package whiskey

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func writeResponse(resp *HttpResponse, writer io.Writer) {
	if resp.statusCode == 0 {
		resp.statusCode = http.StatusOK
	}

	var responseLines []string
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
	_, err := writer.Write([]byte(response))
	if err != nil {
		return
	}
}
