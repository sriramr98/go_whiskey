package whiskey

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestContextBindBody(t *testing.T) {
	type BodyType struct {
		Key1 string `json:"key1"`
		Key2 string `json:"key2"`
	}

	testCases := []struct {
		name        string
		body        []byte
		expected    BodyType
		expectedErr string
	}{
		{
			name:        "Valid JSON with all fields",
			body:        []byte(`{"key1": "value1", "key2": "value2"}`),
			expected:    BodyType{Key1: "value1", Key2: "value2"},
			expectedErr: "",
		},
		{
			name:        "Valid JSON with missing fields",
			body:        []byte(`{"key1": "value1"}`),
			expected:    BodyType{Key1: "value1", Key2: ""},
			expectedErr: "",
		},
		{
			name:        "Invalid JSON",
			body:        []byte(`{"key1": "value1", "key2": "value2"`), // Missing closing brace
			expected:    BodyType{},
			expectedErr: "unexpected end of JSON input",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var body BodyType
			req := HttpRequest{
				body: tc.body,
			}

			context := RequestContext{
				request: req,
			}
			err := context.BindBody(&body)
			if tc.expectedErr == "" && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.expectedErr != "" && (err == nil || err.Error() != tc.expectedErr) {
				t.Fatalf("expected error %v, got %v", tc.expectedErr, err)
			}

			if body != tc.expected {
				t.Fatalf("expected body %v, got %v", tc.expected, body)
			}
		})
	}
}

func TestContextBindQuery(t *testing.T) {
	type QueryType struct {
		StrKey   string  `json:"strKey"`
		IntKey   int     `json:"intKey"`
		BoolKey  bool    `json:"boolKey"`
		FloatKey float64 `json:"floatKey"`
	}

	testCases := []struct {
		name        string
		queryParams map[string]string
		expected    QueryType
		expectedErr string
	}{
		{
			name: "Valid query parameters",
			queryParams: map[string]string{
				"strKey": "value1",
			},
			expected:    QueryType{StrKey: "value1"},
			expectedErr: "",
		},
		{
			name: "Field Types conversion test",
			queryParams: map[string]string{
				"strKey":   "value1",
				"intKey":   "1234",
				"boolKey":  "true",
				"floatKey": "123.45",
			},
			expected:    QueryType{StrKey: "value1", IntKey: 1234, BoolKey: true, FloatKey: 123.45},
			expectedErr: "",
		},
		{
			name:        "Empty query parameters",
			queryParams: map[string]string{},
			expected:    QueryType{},
			expectedErr: "",
		},
		{
			name: "Invalid integer conversion",
			queryParams: map[string]string{
				"intKey": "invalid",
			},
			expected:    QueryType{},
			expectedErr: "invalid value type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := HttpRequest{
				queryParams: tc.queryParams,
			}

			context := RequestContext{
				request: req,
			}
			var query QueryType
			err := context.BindQuery(&query)
			if tc.expectedErr == "" && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.expectedErr != "" && (err == nil || err.Error() != tc.expectedErr) {
				t.Fatalf("expected error %v, got %v", tc.expectedErr, err)
			}

			if query != tc.expected {
				t.Fatalf("expected query %v, got %v", tc.expected, query)
			}
		})
	}
}

func TestContextBindPath(t *testing.T) {
	type PathType struct {
		StrKey   string  `json:"strKey"`
		IntKey   int     `json:"intKey"`
		BoolKey  bool    `json:"boolKey"`
		FloatKey float64 `json:"floatKey"`
	}

	testCases := []struct {
		name        string
		pathParams  map[string]string
		expected    PathType
		expectedErr string
	}{
		{
			name: "Valid path parameters",
			pathParams: map[string]string{
				"strKey": "value1",
			},
			expected:    PathType{StrKey: "value1"},
			expectedErr: "",
		},
		{
			name:        "Empty path parameters",
			pathParams:  map[string]string{},
			expected:    PathType{},
			expectedErr: "",
		},
		{
			name: "Extra path parameter is ignored",
			pathParams: map[string]string{
				"strKey":   "value1",
				"extraKey": "extraValue",
			},
			expected:    PathType{StrKey: "value1"},
			expectedErr: "",
		},
		{
			name: "Path params with type conversion",
			pathParams: map[string]string{
				"strKey":   "value1",
				"intKey":   "1234",
				"boolKey":  "true",
				"floatKey": "123.45",
			},
			expected:    PathType{StrKey: "value1", IntKey: 1234, BoolKey: true, FloatKey: 123.45},
			expectedErr: "",
		},
		{
			name: "Invalid integer conversion",
			pathParams: map[string]string{
				"intKey": "invalid",
			},
			expected:    PathType{},
			expectedErr: "invalid value type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := HttpRequest{
				pathParams: tc.pathParams,
			}

			context := RequestContext{
				request: req,
			}
			var path PathType
			err := context.BindPath(&path)
			if tc.expectedErr == "" && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.expectedErr != "" && (err == nil || err.Error() != tc.expectedErr) {
				t.Fatalf("expected error %v, got %v", tc.expectedErr, err)
			}

			if path != tc.expected {
				t.Fatalf("expected path %v, got %v", tc.expected, path)
			}
		})
	}
}

func TestContextBindHeader(t *testing.T) {
	type HeaderType struct {
		StrKey   string  `json:"strKey"`
		IntKey   int     `json:"intKey"`
		BoolKey  bool    `json:"boolKey"`
		FloatKey float64 `json:"floatKey"`
	}

	testCases := []struct {
		name        string
		headers     map[string]string
		expected    HeaderType
		expectedErr string
	}{
		{
			name: "Valid headers",
			headers: map[string]string{
				"strKey": "value1",
			},
			expected:    HeaderType{StrKey: "value1"},
			expectedErr: "",
		},
		{
			name:        "Empty headers",
			headers:     map[string]string{},
			expected:    HeaderType{},
			expectedErr: "",
		},
		{
			name: "Extra header is ignored",
			headers: map[string]string{
				"strKey":   "value1",
				"extraKey": "extraValue",
			},
			expected:    HeaderType{StrKey: "value1"},
			expectedErr: "",
		},
		{
			name: "Headers with type conversion",
			headers: map[string]string{
				"strKey":   "value1",
				"intKey":   "1234",
				"boolKey":  "true",
				"floatKey": "123.45",
			},
			expected:    HeaderType{StrKey: "value1", IntKey: 1234, BoolKey: true, FloatKey: 123.45},
			expectedErr: "",
		},
		{
			name: "Invalid integer conversion",
			headers: map[string]string{
				"intKey": "invalid",
			},
			expected:    HeaderType{},
			expectedErr: "invalid value type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := HttpRequest{
				headers: tc.headers,
			}

			context := RequestContext{
				request: req,
			}
			var header HeaderType
			err := context.BindHeader(&header)
			if tc.expectedErr == "" && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.expectedErr != "" && (err == nil || err.Error() != tc.expectedErr) {
				t.Fatalf("expected error %v, got %v", tc.expectedErr, err)
			}

			if header != tc.expected {
				t.Fatalf("expected header %v, got %v", tc.expected, header)
			}
		})
	}
}

func TestResponseJson(t *testing.T) {
	type responseData struct {
		Key1 string `json:"key1"`
		Key2 int    `json:"key2"`
	}

	type invalidData struct {
		Channel chan int `json:"channel"`
	}

	testCases := []struct {
		name        string
		data        any
		statusCode  int
		expectedErr string
	}{
		{
			name:       "Empty response",
			data:       responseData{},
			statusCode: http.StatusOK,
		},
		{
			name:       "Non-empty response",
			data:       responseData{Key1: "value1", Key2: 42},
			statusCode: http.StatusOK,
		},
		{
			name:        "Invalid JSON data",
			data:        invalidData{Channel: make(chan int)},
			statusCode:  http.StatusOK,
			expectedErr: "json: unsupported type: chan int",
		},
		{
			name:       "No status code provided",
			data:       responseData{Key1: "value1", Key2: 42},
			statusCode: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := RequestContext{
				response: &HttpResponse{},
			}
			err := ctx.Json(tc.statusCode, tc.data)

			if tc.expectedErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				expectedBody, err := json.Marshal(tc.data)
				if err != nil {
					t.Fatalf("failed to marshal expected body: %v", err)
				}
				if !bytes.Equal(ctx.response.body, expectedBody) {
					t.Fatalf("expected body %s, got %s", string(expectedBody), string(ctx.response.body))
				}

				if ctx.response.headers[HeaderContentType] != MimeTypeJSON {
					t.Fatalf("expected content type %s, got %s", MimeTypeJSON, ctx.response.headers[HeaderContentType])
				}

				if tc.statusCode == 0 {
					tc.statusCode = http.StatusOK
				}
				if ctx.response.statusCode != tc.statusCode {
					t.Fatalf("expected status code %d, got %d", tc.statusCode, ctx.response.statusCode)
				}
			} else {
				if err == nil || err.Error() != tc.expectedErr {
					t.Fatalf("expected error %v, got %v", tc.expectedErr, err)
				}
			}
		})
	}
}

func TestResponseBytes(t *testing.T) {
	testCases := []struct {
		name        string
		data        []byte
		contentType string
		statusCode  int
		expectedErr string
	}{
		{
			name:        "Empty response",
			data:        []byte{},
			contentType: MimeTypeText,
			statusCode:  http.StatusOK,
		},
		{
			name:        "Non-empty response",
			data:        []byte("Hello, World!"),
			contentType: MimeTypeText,
			statusCode:  http.StatusOK,
		},
		{
			name:        "Empty content type",
			data:        []byte("Hello, World!"),
			contentType: "",
			statusCode:  http.StatusOK,
		},
		{
			name:        "Nil data",
			data:        nil,
			contentType: MimeTypeText,
			statusCode:  http.StatusOK,
		},
		{
			name:        "StatusCode Not Provided",
			data:        []byte("Hello, World!"),
			contentType: MimeTypeText,
			statusCode:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := RequestContext{
				response: &HttpResponse{},
			}
			err := ctx.Bytes(tc.statusCode, tc.contentType, tc.data)

			if tc.expectedErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if ctx.response.headers[HeaderContentType] != tc.contentType {
					t.Fatalf("expected content type %s, got %s", tc.contentType, ctx.response.headers[HeaderContentType])
				}
				if !bytes.Equal(ctx.response.body, tc.data) {
					t.Fatalf("expected body %s, got %s", string(tc.data), string(ctx.response.body))
				}
				if tc.statusCode == 0 {
					tc.statusCode = http.StatusOK
				}
				if ctx.response.statusCode != tc.statusCode {
					t.Fatalf("expected status code %d, got %d", tc.statusCode, ctx.response.statusCode)
				}
			} else {
				if err == nil || err.Error() != tc.expectedErr {
					t.Fatalf("expected error %v, got %v", tc.expectedErr, err)
				}
			}
		})
	}
}

func TestResponseString(t *testing.T) {
	testCases := []struct {
		name        string
		data        string
		statusCode  int
		expectedErr string
	}{
		{
			name:       "Valid string response",
			data:       "Hello, World!",
			statusCode: http.StatusOK,
		},
		{
			name:       "Empty string response",
			data:       "",
			statusCode: http.StatusOK,
		},
		{
			name:       "No status code provided",
			data:       "Hello, World!",
			statusCode: 0,
		},
		{
			name:        "Empty string with no status code",
			data:        "",
			statusCode:  0,
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := RequestContext{
				response: &HttpResponse{},
			}
			err := ctx.String(tc.statusCode, tc.data)

			if tc.expectedErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if string(ctx.response.body) != tc.data {
					t.Fatalf("expected body %s, got %s", tc.data, string(ctx.response.body))
				}

				if ctx.response.headers[HeaderContentType] != MimeTypeText {
					t.Fatalf("expected content type %s, got %s", MimeTypeText, ctx.response.headers[HeaderContentType])
				}
				if tc.statusCode == 0 {
					tc.statusCode = http.StatusOK
				}
				if ctx.response.statusCode != tc.statusCode {
					t.Fatalf("expected status code %d, got %d", tc.statusCode, ctx.response.statusCode)
				}
			} else {
				if err == nil || err.Error() != tc.expectedErr {
					t.Fatalf("expected error %v, got %v", tc.expectedErr, err)
				}
			}
		})
	}
}

func TestResponseHtml(t *testing.T) {
	testCases := []struct {
		name        string
		data        string
		statusCode  int
		expectedErr string
	}{
		{
			name:       "Valid HTML response",
			data:       "<h1>Hello, World!</h1>",
			statusCode: http.StatusOK,
		},
		{
			name:       "Empty HTML response",
			data:       "",
			statusCode: http.StatusOK,
		},
		{
			name:       "No status code provided",
			data:       "<h1>Hello, World!</h1>",
			statusCode: 0,
		},
		{
			name:        "Empty HTML with no status code",
			data:        "",
			statusCode:  0,
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := RequestContext{
				response: &HttpResponse{},
			}
			err := ctx.Html(tc.statusCode, tc.data)

			if tc.expectedErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if string(ctx.response.body) != tc.data {
					t.Fatalf("expected body %s, got %s", tc.data, string(ctx.response.body))
				}
				if ctx.response.headers[HeaderContentType] != MimeTypeHTML {
					t.Fatalf("expected content type %s, got %s", MimeTypeHTML, ctx.response.headers[HeaderContentType])
				}
				if tc.statusCode == 0 {
					tc.statusCode = http.StatusOK
				}
				if ctx.response.statusCode != tc.statusCode {
					t.Fatalf("expected status code %d, got %d", tc.statusCode, ctx.response.statusCode)
				}
			} else {
				if err == nil || err.Error() != tc.expectedErr {
					t.Fatalf("expected error %v, got %v", tc.expectedErr, err)
				}
			}
		})
	}
}

func TestGetQueryParam(t *testing.T) {
	req := HttpRequest{
		queryParams: map[string]string{
			"key1": "value1",
		},
	}

	ctx := RequestContext{
		request: req,
	}

	value, exists := ctx.GetQueryParam("key1")
	if !exists || value != "value1" {
		t.Fatalf("expected value 'value1', got '%s'", value)
	}

	_, exists = ctx.GetQueryParam("key2")
	if exists {
		t.Fatalf("expected key2 to not exist")
	}
}

func TestGetQueryParam_ErrorCases(t *testing.T) {
	req := HttpRequest{
		queryParams: map[string]string{},
	}

	ctx := RequestContext{
		request: req,
	}

	_, exists := ctx.GetQueryParam("nonexistent")
	if exists {
		t.Fatalf("expected nonexistent query param to not exist")
	}
}

func TestGetPathParam(t *testing.T) {
	req := HttpRequest{
		pathParams: map[string]string{
			"key1": "value1",
		},
	}

	ctx := RequestContext{
		request: req,
	}

	value, exists := ctx.GetPathParam("key1")
	if !exists || value != "value1" {
		t.Fatalf("expected value 'value1', got '%s'", value)
	}

	_, exists = ctx.GetPathParam("key2")
	if exists {
		t.Fatalf("expected key2 to not exist")
	}
}

func TestGetPathParam_ErrorCases(t *testing.T) {
	req := HttpRequest{
		pathParams: map[string]string{},
	}

	ctx := RequestContext{
		request: req,
	}

	_, exists := ctx.GetPathParam("nonexistent")
	if exists {
		t.Fatalf("expected nonexistent path param to not exist")
	}
}

func TestGetHeader(t *testing.T) {
	req := HttpRequest{
		headers: map[string]string{
			"key1": "value1",
		},
	}

	ctx := RequestContext{
		request: req,
	}

	value, exists := ctx.GetHeader("key1")
	if !exists || value != "value1" {
		t.Fatalf("expected value 'value1', got '%s'", value)
	}

	_, exists = ctx.GetHeader("key2")
	if exists {
		t.Fatalf("expected key2 to not exist")
	}
}

func TestGetHeader_ErrorCases(t *testing.T) {
	req := HttpRequest{
		headers: map[string]string{},
	}

	ctx := RequestContext{
		request: req,
	}

	_, exists := ctx.GetHeader("nonexistent")
	if exists {
		t.Fatalf("expected nonexistent header to not exist")
	}
}

func TestGetQueryParams(t *testing.T) {
	req := HttpRequest{
		queryParams: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	ctx := RequestContext{
		request: req,
	}

	params := ctx.GetQueryParams()
	if len(params) != 2 || params["key1"] != "value1" || params["key2"] != "value2" {
		t.Fatalf("expected query params map with key1=value1 and key2=value2, got %v", params)
	}
}

func TestGetPathParams(t *testing.T) {
	req := HttpRequest{
		pathParams: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	ctx := RequestContext{
		request: req,
	}

	params := ctx.GetPathParams()
	if len(params) != 2 || params["key1"] != "value1" || params["key2"] != "value2" {
		t.Fatalf("expected path params map with key1=value1 and key2=value2, got %v", params)
	}
}

func TestGetHeaders(t *testing.T) {
	req := HttpRequest{
		headers: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	ctx := RequestContext{
		request: req,
	}

	headers := ctx.GetHeaders()
	if len(headers) != 2 || headers["key1"] != "value1" || headers["key2"] != "value2" {
		t.Fatalf("expected headers map with key1=value1 and key2=value2, got %v", headers)
	}
}
