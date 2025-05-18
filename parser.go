package whiskey

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

type RequestParser func(data []string) (HttpRequest, error) // Parse takes a list of lines as input and parses it into a valid HttpRequest

// HTTP_1_1_Parser parses the incoming data according to HTTP 1.1 Specification
func HTTP_1_1_Parser(requestData string) (HttpRequest, error) {
	request := HttpRequest{
		headers:     make(map[string]string),
		queryParams: make(map[string]string),
		pathParams:  make(map[string]string),
	}

	if len(requestData) == 0 {
		return request, fmt.Errorf("invalid HTTP request")
	}

	requestParts := strings.SplitN(requestData, "\r\n\r\n", 2) // The headers and body are seprated by \r\n\r\n
	headerParts := strings.Split(requestParts[0], "\r\n")
	protocolLine := headerParts[0]

	var body string
	if len(requestParts) == 2 {
		body = requestParts[1]
	}

	var headers []string
	if len(headerParts) > 1 {
		headers = headerParts[1:]
	}

	// 1st line should contain the format {method} {path} HTTP/1.1
	protocolParts := strings.Split(strings.TrimSpace(protocolLine), " ")
	if len(protocolParts) < 3 {
		return request, fmt.Errorf("invalid HTTP request format")
	}
	if !slices.Contains([]string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodHead, http.MethodOptions}, protocolParts[0]) {
		return request, fmt.Errorf("invalid HTTP method")
	}
	request.method = protocolParts[0]

	if !strings.HasPrefix(protocolParts[1], "/") {
		return request, fmt.Errorf("invalid HTTP path")
	}

	fullPath := protocolParts[1]
	pathParts := strings.SplitN(fullPath, "?", 2)
	request.path = pathParts[0]

	if request.path != "/" {
		request.path = strings.TrimSuffix(request.path, "/") // Remove trailing /
	}

	if len(pathParts) > 1 {
		queryParamsStr := pathParts[1]
		request.queryParams = parseQueryParams(queryParamsStr)
	}

	// We currently only support HTTP/1.1
	if protocolParts[2] != ProtocolHTTP {
		return HttpRequest{}, fmt.Errorf("invalid HTTP version")
	}

	for _, header := range headers {
		headerParts := strings.SplitN(header, ":", 2)
		if len(headerParts) < 2 {
			return HttpRequest{}, fmt.Errorf("invalid header %s", header)
		}

		key := strings.TrimSpace(headerParts[0])
		value := strings.TrimSpace(headerParts[1])

		request.headers[key] = value
	}

	contentLength, ok := request.headers[HeaderContentLength]
	if !ok {
		if len(body) > 0 {
			return HttpRequest{}, errors.New("request body not expected since Content-Length is 0")
		}
	} else {
		expectedLength, err := strconv.Atoi(contentLength)
		if err != nil {
			return HttpRequest{}, errors.New("invalid Content-Length header value")
		}

		if expectedLength < len(body) {
			return HttpRequest{}, errors.New("body length higher than expected")
		}

		if expectedLength > len(body) {
			return HttpRequest{}, errors.New("incomplete body")
		}
	}

	request.body = []byte(body)

	return request, nil
}

func parseQueryParams(queryParamsStr string) map[string]string {
	queryParams := make(map[string]string)
	if queryParamsStr == "" {
		return queryParams
	}
	params := strings.SplitSeq(queryParamsStr, "&")
	for param := range params {
		paramParts := strings.Split(param, "=")
		value, err := url.QueryUnescape(paramParts[1])
		if err != nil {
			log.Printf("Err parsing query param %v\n", err)
			// Ignore faulty query params
			continue
		}

		queryParams[paramParts[0]] = value
	}

	return queryParams
}
