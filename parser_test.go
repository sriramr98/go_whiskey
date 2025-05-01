package whiskey

import (
	"net/http"
	"testing"
)

func TestHTTP_1_1_Parser(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameter for target function - now a single string with \r\n delimiters
		requestData string
		want        HttpRequest
		wantErr     bool
	}{
		// HAPPY PATH: Basic requests
		{
			name: "Valid GET request (No query params, no body)",
			requestData: "GET /hello HTTP/1.1\r\n" +
				"Host: localhost\r\n" +
				"Content-Type: text/plain\r\n\r\n",
			want: HttpRequest{
				path:        "/hello",
				method:      http.MethodGet,
				queryParams: map[string]string{},
				headers: map[string]string{
					"Host":         "localhost",
					"Content-Type": "text/plain",
				},
			},
			wantErr: false,
		},
		{
			name: "Valid POST request with query params no body",
			requestData: "POST /hello?query1=query2 HTTP/1.1\r\n" +
				"Host: localhost\r\n\r\n",
			want: HttpRequest{
				path:   "/hello",
				method: http.MethodPost,
				queryParams: map[string]string{
					"query1": "query2",
				},
				headers: map[string]string{
					"Host": "localhost",
				},
			},
			wantErr: false,
		},
		{
			name: "Valid PUT request with multiple query params",
			requestData: "PUT /users?id=123&action=update HTTP/1.1\r\n" +
				"Host: api.example.com\r\n" +
				"Content-Type: application/json\r\n" +
				"Authorization: Bearer token123\r\n\r\n",
			want: HttpRequest{
				path:   "/users",
				method: http.MethodPut,
				queryParams: map[string]string{
					"id":     "123",
					"action": "update",
				},
				headers: map[string]string{
					"Host":          "api.example.com",
					"Content-Type":  "application/json",
					"Authorization": "Bearer token123",
				},
				body: []byte{},
			},
			wantErr: false,
		},
		{
			name: "Valid DELETE request with no body",
			requestData: "DELETE /resources/42 HTTP/1.1\r\n" +
				"Host: api.example.com\r\n\r\n",
			want: HttpRequest{
				path:        "/resources/42",
				method:      http.MethodDelete,
				queryParams: map[string]string{},
				headers: map[string]string{
					"Host": "api.example.com",
				},
			},
			wantErr: false,
		},

		// HAPPY PATH: Headers scenarios
		{
			name: "Valid request with multiple headers",
			requestData: "GET /index.html HTTP/1.1\r\n" +
				"Host: www.example.com\r\n" +
				"User-Agent: Mozilla/5.0\r\n" +
				"Accept: text/html\r\n" +
				"Accept-Language: en-US\r\n" +
				"Accept-Encoding: gzip, deflate\r\n" +
				"Connection: keep-alive\r\n\r\n",
			want: HttpRequest{
				path:        "/index.html",
				method:      http.MethodGet,
				queryParams: map[string]string{},
				headers: map[string]string{
					"Host":            "www.example.com",
					"User-Agent":      "Mozilla/5.0",
					"Accept":          "text/html",
					"Accept-Language": "en-US",
					"Accept-Encoding": "gzip, deflate",
					"Connection":      "keep-alive",
				},
			},
			wantErr: false,
		},
		{
			name: "Valid request with header containing colon in value",
			requestData: "GET /time HTTP/1.1\r\n" +
				"Host: example.com\r\n" +
				"Authorization: Basic dXNlcjpwYXNzd29yZA==\r\n" +
				"Time-Zone: GMT+2:00\r\n\r\n",
			want: HttpRequest{
				path:        "/time",
				method:      http.MethodGet,
				queryParams: map[string]string{},
				headers: map[string]string{
					"Host":          "example.com",
					"Authorization": "Basic dXNlcjpwYXNzd29yZA==",
					"Time-Zone":     "GMT+2:00",
				},
			},
			wantErr: false,
		},

		// HAPPY PATH: Path variations
		{
			name: "Valid request with root path",
			requestData: "GET / HTTP/1.1\r\n" +
				"Host: example.com\r\n\r\n",
			want: HttpRequest{
				path:        "/",
				method:      http.MethodGet,
				queryParams: map[string]string{},
				headers: map[string]string{
					"Host": "example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "Valid request with complex path",
			requestData: "GET /api/v1/users/profile/settings HTTP/1.1\r\n" +
				"Host: api.example.com\r\n\r\n",
			want: HttpRequest{
				path:        "/api/v1/users/profile/settings",
				method:      http.MethodGet,
				queryParams: map[string]string{},
				headers: map[string]string{
					"Host": "api.example.com",
				},
			},
			wantErr: false,
		},

		// HAPPY PATH: Query parameter variations
		{
			name: "Valid request with multiple query parameters",
			requestData: "GET /search?q=test&page=1&limit=10&sort=desc HTTP/1.1\r\n" +
				"Host: search.example.com\r\n\r\n",
			want: HttpRequest{
				path:   "/search",
				method: http.MethodGet,
				queryParams: map[string]string{
					"q":     "test",
					"page":  "1",
					"limit": "10",
					"sort":  "desc",
				},
				headers: map[string]string{
					"Host": "search.example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "Valid request with empty query parameter value",
			requestData: "GET /filter?category=books&price= HTTP/1.1\r\n" +
				"Host: store.example.com\r\n\r\n",
			want: HttpRequest{
				path:   "/filter",
				method: http.MethodGet,
				queryParams: map[string]string{
					"category": "books",
					"price":    "",
				},
				headers: map[string]string{
					"Host": "store.example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "Valid request with URL-encoded query parameters",
			requestData: "GET /search?q=hello%20world&tag=example%3Atest HTTP/1.1\r\n" +
				"Host: example.com\r\n\r\n",
			want: HttpRequest{
				path:   "/search",
				method: http.MethodGet,
				queryParams: map[string]string{
					"q":   "hello world",
					"tag": "example:test",
				},
				headers: map[string]string{
					"Host": "example.com",
				},
			},
			wantErr: false,
		},

		// HAPPY PATH: Different methods
		{
			name: "Valid PATCH request with body",
			requestData: "PATCH /api/users/123 HTTP/1.1\r\n" +
				"Host: api.example.com\r\n" +
				"Content-Type: application/json\r\n" +
				"Content-Length: 22\r\n" +
				"\r\n" +
				"{\"status\": \"inactive\"}",
			want: HttpRequest{
				path:        "/api/users/123",
				method:      "PATCH",
				queryParams: map[string]string{},
				headers: map[string]string{
					"Host":           "api.example.com",
					"Content-Type":   "application/json",
					"Content-Length": "22",
				},
				body: []byte("{\"status\": \"inactive\"}"),
			},
			wantErr: false,
		},
		{
			name: "Valid HEAD request",
			requestData: "HEAD /index.html HTTP/1.1\r\n" +
				"Host: example.com\r\n\r\n",
			want: HttpRequest{
				path:        "/index.html",
				method:      "HEAD",
				queryParams: map[string]string{},
				headers: map[string]string{
					"Host": "example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "Valid OPTIONS request",
			requestData: "OPTIONS /api/resources HTTP/1.1\r\n" +
				"Host: api.example.com\r\n" +
				"Access-Control-Request-Method: POST\r\n" +
				"Access-Control-Request-Headers: Content-Type, Authorization\r\n\r\n",
			want: HttpRequest{
				path:        "/api/resources",
				method:      "OPTIONS",
				queryParams: map[string]string{},
				headers: map[string]string{
					"Host":                           "api.example.com",
					"Access-Control-Request-Method":  "POST",
					"Access-Control-Request-Headers": "Content-Type, Authorization",
				},
			},
			wantErr: false,
		},

		// HAPPY PATH: Body variations
		{
			name: "Valid POST request with JSON body",
			requestData: "POST /api/users HTTP/1.1\r\n" +
				"Host: api.example.com\r\n" +
				"Content-Type: application/json\r\n" +
				"Content-Length: 53\r\n" +
				"\r\n" +
				"{\"username\":\"johndoe\",\"email\":\"john.doe@example.com\"}",
			want: HttpRequest{
				path:        "/api/users",
				method:      http.MethodPost,
				queryParams: map[string]string{},
				headers: map[string]string{
					"Host":           "api.example.com",
					"Content-Type":   "application/json",
					"Content-Length": "53",
				},
				body: []byte("{\"username\":\"johndoe\",\"email\":\"john.doe@example.com\"}"),
			},
			wantErr: false,
		},
		{
			name: "Valid POST request with form data",
			requestData: "POST /login HTTP/1.1\r\n" +
				"Host: example.com\r\n" +
				"Content-Type: application/x-www-form-urlencoded\r\n" +
				"Content-Length: 29\r\n" +
				"\r\n" +
				"username=john&password=secret",
			want: HttpRequest{
				path:        "/login",
				method:      http.MethodPost,
				queryParams: map[string]string{},
				headers: map[string]string{
					"Host":           "example.com",
					"Content-Type":   "application/x-www-form-urlencoded",
					"Content-Length": "29",
				},
				body: []byte("username=john&password=secret"),
			},
			wantErr: false,
		},
		{
			name: "Valid POST request with multi-line body",
			requestData: "POST /api/messages HTTP/1.1\r\n" +
				"Host: api.example.com\r\n" +
				"Content-Type: text/plain\r\n" +
				"Content-Length: 64\r\n" +
				"\r\n" +
				"This is a test message.\r\n" +
				"It has multiple lines.\r\n" +
				"End of message.",
			want: HttpRequest{
				path:        "/api/messages",
				method:      http.MethodPost,
				queryParams: map[string]string{},
				headers: map[string]string{
					"Host":           "api.example.com",
					"Content-Type":   "text/plain",
					"Content-Length": "64",
				},
				body: []byte("This is a test message.\r\nIt has multiple lines.\r\nEnd of message."),
			},
			wantErr: false,
		},

		// ERROR CASES: Malformed request lines
		{
			name:        "Error - Empty request data",
			requestData: "",
			want:        HttpRequest{},
			wantErr:     true,
		},
		{
			name: "Error - Missing HTTP version",
			requestData: "GET /index.html\r\n" +
				"Host: example.com\r\n\r\n",
			want:    HttpRequest{},
			wantErr: true,
		},
		{
			name: "Error - Invalid HTTP version (not 1.1)",
			requestData: "GET /index.html HTTP/2.0\r\n" +
				"Host: example.com\r\n\r\n",
			want:    HttpRequest{},
			wantErr: true,
		},
		{
			name: "Error - Missing path",
			requestData: "GET HTTP/1.1\r\n" +
				"Host: example.com\r\n\r\n",
			want:    HttpRequest{},
			wantErr: true,
		},
		{
			name: "Error - Missing method",
			requestData: "/index.html HTTP/1.1\r\n" +
				"Host: example.com\r\n\r\n",
			want:    HttpRequest{},
			wantErr: true,
		},
		{
			name: "Error - Invalid method",
			requestData: "INVALID /index.html HTTP/1.1\r\n" +
				"Host: example.com\r\n\r\n",
			want:    HttpRequest{},
			wantErr: true,
		},

		// ERROR CASES: Header issues
		{
			name: "Error - Malformed header (missing colon)",
			requestData: "GET /index.html HTTP/1.1\r\n" +
				"Host: example.com\r\n" +
				"Content-Type text/html\r\n\r\n",
			want:    HttpRequest{},
			wantErr: true,
		},

		// ERROR CASES: Body issues
		{
			name: "Error - Content-Length doesn't match actual body length",
			requestData: "POST /api/users HTTP/1.1\r\n" +
				"Host: api.example.com\r\n" +
				"Content-Type: application/json\r\n" +
				"Content-Length: 100\r\n" + // Incorrect length
				"\r\n" +
				"{\"username\":\"johndoe\",\"email\":\"john.doe@example.com\"}",
			want:    HttpRequest{},
			wantErr: true,
		},
		{
			name: "Error - Body present for GET request with no Content-Length",
			requestData: "GET /index.html HTTP/1.1\r\n" +
				"Host: example.com\r\n" +
				"\r\n" +
				"This body shouldn't be here",
			want:    HttpRequest{},
			wantErr: true,
		},

		// ERROR CASES: Path and query parameters
		{
			name: "Error - Malformed query parameter (missing value after =)",
			requestData: "GET /search?q=test&invalid= HTTP/1.1\r\n" +
				"Host: example.com\r\n\r\n",
			want: HttpRequest{
				path:   "/search",
				method: http.MethodGet,
				queryParams: map[string]string{
					"q":       "test",
					"invalid": "",
				},
				headers: map[string]string{
					"Host": "example.com",
				},
			},
			wantErr: false, // Note: This is actually a valid case in HTTP - empty values are allowed
		},
		{
			name: "Error - Malformed path with spaces",
			requestData: "GET /path with spaces HTTP/1.1\r\n" +
				"Host: example.com\r\n\r\n",
			want:    HttpRequest{},
			wantErr: true,
		},

		// EDGE CASES
		{
			name: "Edge case - Request with empty headers",
			requestData: "GET /index.html HTTP/1.1\r\n" +
				"Host: example.com\r\n" +
				"EmptyHeader: \r\n\r\n",
			want: HttpRequest{
				path:        "/index.html",
				method:      http.MethodGet,
				queryParams: map[string]string{},
				headers: map[string]string{
					"Host":        "example.com",
					"EmptyHeader": "",
				},
			},
			wantErr: false,
		},
		{
			name: "Edge case - Request with extremely long path",
			requestData: "GET /this/is/a/very/long/path/that/goes/on/and/on/with/many/segments/to/test/parser/limits HTTP/1.1\r\n" +
				"Host: example.com\r\n\r\n",
			want: HttpRequest{
				path:        "/this/is/a/very/long/path/that/goes/on/and/on/with/many/segments/to/test/parser/limits",
				method:      http.MethodGet,
				queryParams: map[string]string{},
				headers: map[string]string{
					"Host": "example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "Edge case - Case-insensitive header names",
			requestData: "GET /index.html HTTP/1.1\r\n" +
				"HOST: example.com\r\n" + // Capitalized differently from typical "Host"
				"content-type: text/html\r\n\r\n", // Lowercase instead of typical "Content-Type"
			want: HttpRequest{
				path:        "/index.html",
				method:      http.MethodGet,
				queryParams: map[string]string{},
				headers: map[string]string{
					"HOST":         "example.com",
					"content-type": "text/html",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := HTTP_1_1_Parser(tt.requestData)

			if tt.wantErr && gotErr == nil {
				t.Errorf("Expected to error but didn't get any error")
				return
			}

			if !tt.wantErr && gotErr != nil {
				t.Errorf("Expected to not error but got error %+v", gotErr)
				return
			}

			if !tt.want.Equal(got) {
				t.Errorf("Expected %+v but got %+v", tt.want, got)
			}
		})
	}
}
