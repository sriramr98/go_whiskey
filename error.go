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

func NewHttpError(statusCode int, bodyType BodyType) HttpError {
	var body string

	switch bodyType {
	case BodyTypeJSON:
		body = fmt.Sprintf("{\"error\": \"%s\"}", http.StatusText(statusCode))
	case BodyTypeString:
	default:
		body = http.StatusText(statusCode)
	}

	return HttpError{
		StatusCode: statusCode,
		Body:       body,
		BodyType:   bodyType,
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
