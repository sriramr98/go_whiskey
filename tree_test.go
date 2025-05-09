package whiskey

import (
	"errors"
	"testing"
)

// Test handlers that return different errors to distinguish them
func handlerOne(c Context) error {
	return errors.New("handler one called")
}

func handlerTwo(c Context) error {
	return errors.New("handler two called")
}

func handlerThree(c Context) error {
	return errors.New("handler three called")
}

func TestRouteTreeInsert(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		method      string
		shouldError bool
	}{
		{"Valid path", "/api/v1/users", "GET", false},
		{"Valid path with trailing slash", "/api/v1/users/", "POST", false},
		{"Root path", "/", "GET", false},
		{"Path with params", "/users/{id}", "GET", false},
		{"Empty path", "", "GET", true},
		{"Missing leading slash", "api/v1/users", "GET", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tree := newRouteTree()
			err := tree.insert(tc.path, tc.method, routeConfig{handlers: []HttpHandler{handlerOne}})

			if tc.shouldError && err == nil {
				t.Errorf("Expected error for path %s, but got none", tc.path)
			}

			if !tc.shouldError && err != nil {
				t.Errorf("Unexpected error for path %s: %v", tc.path, err)
			}
		})
	}
}

func TestRouteTreeGetConfig(t *testing.T) {
	tests := []struct {
		name      string
		setupTree func() *routeTree
		path      string
		method    string
		found     bool
	}{
		{
			name: "Basic route",
			setupTree: func() *routeTree {
				tree := newRouteTree()
				_ = tree.insert("/api/v1/users", "GET", routeConfig{handlers: []HttpHandler{handlerOne}})
				return tree
			},
			path:   "/api/v1/users",
			method: "GET",
			found:  true,
		},
		{
			name: "Route with trailing slash",
			setupTree: func() *routeTree {
				tree := newRouteTree()
				_ = tree.insert("/api/v1/users/", "GET", routeConfig{handlers: []HttpHandler{handlerOne}})
				return tree
			},
			path:   "/api/v1/users",
			method: "GET",
			found:  true,
		},
		{
			name: "Root path",
			setupTree: func() *routeTree {
				tree := newRouteTree()
				_ = tree.insert("/", "GET", routeConfig{handlers: []HttpHandler{handlerOne}})
				return tree
			},
			path:   "/",
			method: "GET",
			found:  true,
		},
		{
			name: "Non-existent path",
			setupTree: func() *routeTree {
				tree := newRouteTree()
				_ = tree.insert("/api/v1/users", "GET", routeConfig{handlers: []HttpHandler{handlerOne}})
				return tree
			},
			path:   "/api/v2/users",
			method: "GET",
			found:  false,
		},
		{
			name: "Method not allowed",
			setupTree: func() *routeTree {
				tree := newRouteTree()
				_ = tree.insert("/api/v1/users", "GET", routeConfig{handlers: []HttpHandler{handlerOne}})
				return tree
			},
			path:   "/api/v1/users",
			method: "POST",
			found:  false,
		},
		{
			name: "Path with params",
			setupTree: func() *routeTree {
				tree := newRouteTree()
				_ = tree.insert("/users/{id}", "GET", routeConfig{handlers: []HttpHandler{handlerOne}})
				return tree
			},
			path:   "/users/{id}",
			method: "GET",
			found:  true,
		},
		{
			name: "Empty path",
			setupTree: func() *routeTree {
				tree := newRouteTree()
				_ = tree.insert("/api", "GET", routeConfig{handlers: []HttpHandler{handlerOne}})
				return tree
			},
			path:   "",
			method: "GET",
			found:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tree := tc.setupTree()
			_, found := tree.getConfig(tc.path, tc.method)

			if found != tc.found {
				t.Errorf("Expected found=%v, got found=%v for path=%s, method=%s",
					tc.found, found, tc.path, tc.method)
			}
		})
	}
}

func TestMultipleRoutes(t *testing.T) {
	tree := newRouteTree()

	// Insert multiple routes
	_ = tree.insert("/api/v1/users", "GET", routeConfig{handlers: []HttpHandler{handlerOne}})
	_ = tree.insert("/api/v1/users", "POST", routeConfig{handlers: []HttpHandler{handlerTwo}})
	_ = tree.insert("/api/v2/users", "GET", routeConfig{handlers: []HttpHandler{handlerThree}})

	tests := []struct {
		path      string
		method    string
		found     bool
		handlerID string
	}{
		{"/api/v1/users", "GET", true, "handler one called"},
		{"/api/v1/users", "POST", true, "handler two called"},
		{"/api/v2/users", "GET", true, "handler three called"},
		{"/api/v1/users", "PUT", false, ""},
		{"/api/v3/users", "GET", false, ""},
	}

	for _, tc := range tests {
		t.Run(tc.path+"-"+tc.method, func(t *testing.T) {
			config, found := tree.getConfig(tc.path, tc.method)

			if found != tc.found {
				t.Errorf("Expected found=%v, got found=%v for path=%s, method=%s",
					tc.found, found, tc.path, tc.method)
				return
			}

			if found {
				// Execute the handler to verify it's the correct one
				mockContext := RequestContext{}
				err := config.handlers[0](mockContext)
				if err.Error() != tc.handlerID {
					t.Errorf("Expected handler %s, got %s", tc.handlerID, err.Error())
				}
			}
		})
	}
}

func TestNestedRoutes(t *testing.T) {
	tree := newRouteTree()

	// Insert nested routes
	_ = tree.insert("/api", "GET", routeConfig{handlers: []HttpHandler{handlerOne}})
	_ = tree.insert("/api/v1", "GET", routeConfig{handlers: []HttpHandler{handlerTwo}})
	_ = tree.insert("/api/v1/users", "GET", routeConfig{handlers: []HttpHandler{handlerThree}})

	tests := []struct {
		path      string
		method    string
		found     bool
		handlerID string
	}{
		{"/api", "GET", true, "handler one called"},
		{"/api/v1", "GET", true, "handler two called"},
		{"/api/v1/users", "GET", true, "handler three called"},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			config, found := tree.getConfig(tc.path, tc.method)

			if found != tc.found {
				t.Errorf("Expected found=%v, got found=%v for path=%s",
					tc.found, found, tc.path)
				return
			}

			if found {
				// Execute the handler to verify it's the correct one
				mockContext := RequestContext{}
				err := config.handlers[0](mockContext)
				if err.Error() != tc.handlerID {
					t.Errorf("Expected handler %s, got %s", tc.handlerID, err.Error())
				}
			}
		})
	}
}

func TestPathParams(t *testing.T) {
	tree := newRouteTree()

	// Insert routes with path parameters
	_ = tree.insert("/users/{id}", "GET", routeConfig{handlers: []HttpHandler{handlerOne}})
	_ = tree.insert("/posts/{postId}/comments/{commentId}", "GET", routeConfig{handlers: []HttpHandler{handlerTwo}})

	tests := []struct {
		path       string
		method     string
		found      bool
		pathParams map[string]string
	}{
		{"/users/{id}", "GET", true, make(map[string]string)},
		{"/posts/{postId}/comments/{commentId}", "GET", true, make(map[string]string)},
		{"/users/123", "GET", true, map[string]string{"id": "123"}},
		{"/posts/123/comments/456", "GET", true, map[string]string{"postId": "123", "commentId": "456"}},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			config, found := tree.getConfig(tc.path, tc.method)

			if found != tc.found {
				t.Errorf("Expected found=%v, got found=%v for path=%s",
					tc.found, found, tc.path)
			}

			if !CompareMaps(tc.pathParams, config.pathParams) {
				t.Errorf("Expected pathParams=%v but got %v for path %s", tc.pathParams, config.pathParams, tc.path)
			}
		})
	}
}

func TestIsPathParam(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"{id}", true},
		{"{userId}", true},
		{"id", false},
		{"{id", false},
		{"id}", false},
		{"", false},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := isPathParam(tc.input)
			if result != tc.expected {
				t.Errorf("isPathParam(%s) = %v; want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestSlashHandling(t *testing.T) {
	tree := newRouteTree()

	// Insert routes with and without trailing slashes
	_ = tree.insert("/api/v1/users", "GET", routeConfig{handlers: []HttpHandler{handlerOne}})
	_ = tree.insert("/api/v2/posts/", "GET", routeConfig{handlers: []HttpHandler{handlerTwo}})

	tests := []struct {
		name      string
		path      string
		method    string
		found     bool
		handlerID string
	}{
		{"Path without trailing slash", "/api/v1/users", "GET", true, "handler one called"},
		{"Request with trailing slash", "/api/v1/users/", "GET", true, "handler one called"},
		{"Path with trailing slash", "/api/v2/posts", "GET", true, "handler two called"},
		{"Request without trailing slash", "/api/v2/posts/", "GET", true, "handler two called"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config, found := tree.getConfig(tc.path, tc.method)

			if found != tc.found {
				t.Errorf("Expected found=%v, got found=%v for path=%s",
					tc.found, found, tc.path)
				return
			}

			if found {
				// Execute the handler to verify it's the correct one
				mockContext := RequestContext{}
				err := config.handlers[0](mockContext)
				if err.Error() != tc.handlerID {
					t.Errorf("Expected handler %s, got %s", tc.handlerID, err.Error())
				}
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("Multiple inserts on same path", func(t *testing.T) {
		tree := newRouteTree()

		// Insert the same path twice with different methods
		_ = tree.insert("/api/v1/users", "GET", routeConfig{handlers: []HttpHandler{handlerOne}})
		_ = tree.insert("/api/v1/users", "POST", routeConfig{handlers: []HttpHandler{handlerTwo}})

		// Verify both handlers exist
		getConfig, getFound := tree.getConfig("/api/v1/users", "GET")
		postConfig, postFound := tree.getConfig("/api/v1/users", "POST")

		if !getFound || !postFound {
			t.Errorf("Expected both handlers to be found, got GET=%v, POST=%v", getFound, postFound)
			return
		}

		// Verify they're different handlers
		mockContext := RequestContext{}
		getErr := getConfig.handlers[0](mockContext)
		postErr := postConfig.handlers[0](mockContext)

		if getErr.Error() == postErr.Error() {
			t.Errorf("Expected different handlers, got the same: %s", getErr.Error())
		}
	})

	t.Run("Root path", func(t *testing.T) {
		tree := newRouteTree()

		// Insert root path
		_ = tree.insert("/", "GET", routeConfig{handlers: []HttpHandler{handlerOne}})

		// Verify it can be retrieved
		config, found := tree.getConfig("/", "GET")

		if !found {
			t.Error("Expected root path to be found")
			return
		}

		// Verify it's the correct handler
		mockContext := RequestContext{}
		err := config.handlers[0](mockContext)
		if err.Error() != "handler one called" {
			t.Errorf("Expected 'handler one called', got %s", err.Error())
		}
	})
}

func TestCaseSensitivity(t *testing.T) {
	tree := newRouteTree()

	// Insert different case routes
	_ = tree.insert("/api", "GET", routeConfig{handlers: []HttpHandler{handlerOne}})
	_ = tree.insert("/API", "GET", routeConfig{handlers: []HttpHandler{handlerTwo}})

	tests := []struct {
		path      string
		handlerID string
	}{
		{"/api", "handler one called"},
		{"/API", "handler two called"},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			config, found := tree.getConfig(tc.path, "GET")

			if !found {
				t.Errorf("Path %s should be found", tc.path)
				return
			}

			mockContext := RequestContext{}
			err := config.handlers[0](mockContext)
			if err.Error() != tc.handlerID {
				t.Errorf("Expected handler %s, got %s", tc.handlerID, err.Error())
			}
		})
	}
}

func BenchmarkInsert(b *testing.B) {
	tree := newRouteTree()
	paths := []string{
		"/api/v1/users",
		"/api/v1/users/{id}",
		"/api/v1/posts",
		"/api/v1/posts/{id}/comments",
		"/api/v2/users",
		"/api/v2/posts",
		"/api/health",
		"/static/css",
		"/static/js",
		"/static/images",
	}

	b.ResetTimer()
	for b.Loop() {
		for _, path := range paths {
			_ = tree.insert(path, "GET", routeConfig{handlers: []HttpHandler{handlerOne}})
		}
	}
}

func BenchmarkGetConfig(b *testing.B) {
	tree := newRouteTree()
	paths := []string{
		"/api/v1/users",
		"/api/v1/users/{id}",
		"/api/v1/posts",
		"/api/v1/posts/{id}/comments",
		"/api/v2/users",
		"/api/v2/posts",
		"/api/health",
		"/static/css",
		"/static/js",
		"/static/images",
	}

	for _, path := range paths {
		_ = tree.insert(path, "GET", routeConfig{handlers: []HttpHandler{handlerOne}})
	}

	b.ResetTimer()
	for b.Loop() {
		for _, path := range paths {
			_, _ = tree.getConfig(path, "GET")
		}
	}
}
