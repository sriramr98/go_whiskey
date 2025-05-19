package whiskey

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

func (w *Whiskey) writeResponse(resp *HttpResponse, writer io.Writer) {
	if err := writer.(net.Conn).SetWriteDeadline(time.Now().Add(w.config.WriteTimeout)); err != nil {
		w.errorLogger.Println("Unable to set write deadline", err)
		return
	}

	if resp.statusCode == 0 {
		resp.statusCode = http.StatusOK
	}

	// We only support HTTP/1.1
	if _, err := fmt.Fprintf(writer, "HTTP/1.1 %d %s\r\n", resp.statusCode, http.StatusText(resp.statusCode)); err != nil {
		w.errorLogger.Printf("Unable to write response.. %+v", err)
		return
	}

	contentType, ok := resp.headers[HeaderContentType]
	if !ok {
		contentType = "text/plain; charset=utf-8"
	}
	contentLength := len(resp.body)
	resp.headers["Content-Length"] = fmt.Sprintf("%d", contentLength)
	resp.headers["Date"] = time.Now().Format("Mon, 02 January 2006 15:04:05 GMT")
	resp.headers["Content-Type"] = contentType

	// Write the headers to the response stream
	for key, value := range resp.headers {
		if _, err := fmt.Fprintf(writer, "%s: %s\r\n", key, value); err != nil {
			w.errorLogger.Printf("Unable to write response.. %+v", err)
			return
		}
	}

	// var body string = string(resp.body)

	if _, err := fmt.Fprintf(writer, "\r\n%s", resp.body); err != nil {
		w.errorLogger.Printf("Unable to write response.. %+v", err)
	}
}
