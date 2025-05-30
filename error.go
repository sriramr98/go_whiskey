package whiskey

import (
	"fmt"
	"net/http"
)

type HttpError struct {
	StatusCode int
	Body       string
	Message    string
	BodyType   BodyType // Based on the type, sends a text/plain or application/json response
}

type BodyType int

const (
	BodyTypeString BodyType = iota
	BodyTypeJSON
)

func (h HttpError) Error() string {
	if h.Message != "" {
		return h.Message
	}
	return string(h.Body)
}

// NewHttpError returns a response in format { "error": "<message_of_status_code>" } if bodyType is json else "<message_of_status_code>". The message_of_status_code is equal to `http.StatusText`
func NewHttpError(statusCode int, bodyType BodyType) HttpError {
	return NewHTTPErrorWithMessage(statusCode, http.StatusText(statusCode), bodyType)
}

// NewHTTPErrorWithMessage will return a response in the format "{ "error": "<message>" }" if bodyType is JSON else "<message>"
func NewHTTPErrorWithMessage(statusCode int, message string, bodyType BodyType) HttpError {
	var body string

	switch bodyType {
	case BodyTypeJSON:
		body = fmt.Sprintf("{\"error\": \"%s\"}", message)
	case BodyTypeString:
	default:
		body = message
	}

	return HttpError{
		StatusCode: statusCode,
		Body:       body,
		BodyType:   bodyType,
		Message:    message,
	}
}

func defaultErrorHandler(err error, ctx Context) error {
	if err == nil {
		return nil
	}

	if httpErr, ok := err.(HttpError); ok {
		if httpErr.BodyType == BodyTypeString {
			ctx.String(httpErr.StatusCode, httpErr.Body)
		} else {
			ctx.Json(httpErr.StatusCode, httpErr.Body)
		}
	} else {
		ctx.String(http.StatusInternalServerError, err.Error())
	}

	return nil
}
