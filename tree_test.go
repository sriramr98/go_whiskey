package whiskey

import "testing"

func TestTreeInserts(t *testing.T) {
	tree := newRouteTree()

	handlerHit := false
	handler := func(ctx Context) error {
		handlerHit = true
		return nil
	}

	tree.insert("/api/v1/auth", "GET", routeConfig{handler: handler})

	node, ok := tree.roots["api"]
	if !ok {
		t.Fatal("Expected root node with 'api' to be present")
	}

	if node.end {
		t.Fatal("Expected node 'api' to not be end")
	}

	node, ok = node.children["v1"]
	if !ok {
		t.Fatal("Expected node 'v1' to be a child of node 'api'")
	}

	if node.end {
		t.Fatal("Expected node 'v1' to not be end")
	}

	node, ok = node.children["auth"]
	if !ok {
		t.Fatal("Expected node 'auth' to be a child of node 'v1'")
	}

	if !node.end {
		t.Fatal("Expected node 'auth' to be an end node")
	}

	config, ok := node.handlers["GET"]
	if !ok {
		t.Fatal("Expected node 'auth' to have a GET handler")
	}

	config.handler(RequestContext{})

	if !handlerHit {
		t.Fatal("Expected handler to be hit for route /api/v1/auth")
	}

	handlerHit = false
	tree.insert("/api/v1", "GET", routeConfig{handler: handler})

	node, ok = tree.roots["api"]
	if !ok {
		t.Fatal("Expected node 'api' to be a root node")
	}

	if node.end {
		t.Fatal("Expected node 'api' to not be end")
	}

	node, ok = node.children["v1"]
	if !ok {
		t.Fatal("Expected node 'v1' to be a child of node 'api'")
	}

	if !node.end {
		t.Fatal("Expected node 'v1' to be an end node")
	}

	config, ok = node.handlers["GET"]
	if !ok {
		t.Fatal("Expected handler for route /api/v1")
	}

	config.handler(RequestContext{})
	if !handlerHit {
		t.Fatal("Expected handler to be hit for /api/v1")
	}
}

func TestTreeGet(t *testing.T) {
	testCases := []struct {
		name             string
		insertPath       string
		getPath          string
		method           string
		shouldFindRouter bool
	}{
		{
			name:             "test for valid route insert and get",
			insertPath:       "/api/v1/auth",
			getPath:          "/api/v1/auth",
			method:           "GET",
			shouldFindRouter: true,
		},
		{
			name:             "test for valid route with path ending in /",
			insertPath:       "/api/v1/auth/",
			getPath:          "/api/v1/auth/",
			method:           "GET",
			shouldFindRouter: true,
		},
		{
			name:             "test for route for which handler isn't present but route with same prefix exists",
			insertPath:       "/api/v1/auth",
			getPath:          "/api/v1",
			method:           "GET",
			shouldFindRouter: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			routeTree := newRouteTree()
			handlerHit := false
			handler := func(ctx Context) error {
				handlerHit = true
				return nil
			}

			routeTree.insert(tc.insertPath, tc.method, routeConfig{handler: handler})

			config, ok := routeTree.getHandler(tc.getPath, tc.method)
			if tc.shouldFindRouter && !ok {
				t.Fatal("Expected to get router config for path")
			}

			if tc.shouldFindRouter {
				config.handler(RequestContext{})
				if tc.shouldFindRouter && !handlerHit {
					t.Fatal("Expected handler to get called but didn't get called")
				}
			}
		})
	}
}
