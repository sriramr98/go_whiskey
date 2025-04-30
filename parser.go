package whiskey

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
)

type RequestParser func(data []string) (HttpRequest, error) // Parse takes a list of lines as input and parses it into a valid HttpRequest

// HTTP_1_1_Parser parses the incoming data according to HTTP 1.1 Specification
func HTTP_1_1_Parser(requestData []string) (HttpRequest, error) {
	request := HttpRequest{
		headers:     make(map[string]string),
		queryParams: make(map[string]string),
		pathParams:  make(map[string]string),
	}

	if len(requestData) == 0 {
		return request, fmt.Errorf("invalid HTTP request")
	}

	// 1st line should contain the format {method} {path} HTTP/1.1
	protocolParts := strings.Split(strings.TrimSpace(requestData[0]), " ")
	if len(protocolParts) < 3 {
		return request, fmt.Errorf("invalid HTTP request format")
	}
	if !slices.Contains([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}, protocolParts[0]) {
		return request, fmt.Errorf("invalid HTTP method")
	}
	request.method = protocolParts[0]

	if !strings.HasPrefix(protocolParts[1], "/") {
		return request, fmt.Errorf("invalid HTTP path")
	}

	fullPath := protocolParts[1]
	pathParts := strings.SplitN(fullPath, "?", 2)
	fmt.Printf("Path %v", pathParts)
	request.path = strings.TrimSuffix(pathParts[0], "/") // Removes trailing / if there's one

	if len(pathParts) > 1 {
		queryParamsStr := pathParts[1]
		if queryParamsStr != "" {
			params := strings.SplitSeq(queryParamsStr, "&")
			for param := range params {
				paramParts := strings.Split(param, "=")
				request.queryParams[paramParts[0]] = paramParts[1]
			}
		}
	}

	// We currently only support HTTP/1.1
	if protocolParts[2] != ProtocolHTTP {
		return request, fmt.Errorf("invalid HTTP version")
	}

	// From the second line, it contains headers in the format {key}: {value} until we find an empty line

	bodyStartIdx := 1
	lastReadIdx := 1
	for idx, line := range requestData {
		if idx == 0 {
			continue
		}
		if line == "" {
			// An empty line indicates the end of headers
			bodyStartIdx = idx + 1 // The body starts after the empty line
			break
		}
		lastReadIdx = idx

		headerParts := strings.SplitN(line, ":", 2)
		if len(headerParts) < 2 {
			// ignore broken headers
			continue
		}

		key := strings.TrimSpace(headerParts[0])
		value := strings.TrimSpace(headerParts[1])

		request.headers[key] = value
	}

	// body isn't present, and we reached the end of the request
	// We do lastReadIdx+1 because the header loop starts from 1st index which is 0 inside the loop
	if lastReadIdx == len(requestData)-1 {
		return request, nil
	}

	// The body starts after the headers
	// TODO: Use `Content-Length` header to read body based on the value. If the body has less bytes than the `Content-Length` value, then the connection was closed as the client was writing the request. Ref - https://arc.net/l/quote/tpfzwcyo
	body := strings.Join(requestData[bodyStartIdx:], "")
	request.body = []byte(body)

	return request, nil
}
