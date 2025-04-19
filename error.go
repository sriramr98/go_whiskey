package whiskey

import "net/http"

func defaultErrorHandler(err error, ctx Context) error {
	if err == nil {
		return nil
	}
	ctx.SetHeader("Content-Type", MimeTypeText)
	ctx.String(http.StatusInternalServerError, err.Error())
	return nil
}

type HttpError struct {
	StatusCode int
	Body       []byte
}

func (h *HttpError) Error() string {
	return "HTTP Error"
}
