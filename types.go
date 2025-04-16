package whiskey

type RunOpts struct {
	Port int
	Addr string
}

type HttpRequest struct {
	Method  string
	Path    string
	Body    []byte
	Headers map[string]string
}

type HttpResponse struct {
	statusCode int
	body       []byte
	headers    map[string]string
}

func (resp *HttpResponse) SetHeader(key string, value string) {
	if resp.headers == nil {
		resp.headers = make(map[string]string)
	}
	resp.headers[key] = value
}

func (resp *HttpResponse) Send(body []byte) {
	resp.body = body
}

type HttpHandler func(HttpRequest, *HttpResponse)
