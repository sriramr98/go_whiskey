package whiskey

import (
	"errors"
	"strings"
)

// A grossly simplified immutable implementation of a radix tree
// For a given path `/api/v1/auth`, api will be a root node which will have a child v1 which will have a child auth
// If there's another route `/apu/v2/auth`, then apu will be a root node separately.
// routeTree isn't thread safe as it isn't expected to be used across routines
type routeTree struct {
	roots map[string]*node
}

type routeConfig struct {
	handler HttpHandler
}

type node struct {
	key      string
	children map[string]*node
	handlers map[string]routeConfig // For every http method, there can be a handler
	end      bool                   // denotes if this node specifies the end of a valid route
}

func newRouteTree() *routeTree {
	return &routeTree{
		roots: make(map[string]*node),
	}
}

// insert creates a new set of nodes for the given path. path is expected to be in the format /{part1}/{part2}...
func (t *routeTree) insert(path string, method string, config routeConfig) error {
	if path == "" {
		return errors.New("invalid path " + path)
	}

	if !strings.HasPrefix(path, "/") {
		return errors.New("invalid path " + path)
	}

	var pathParts []string
	trimmedPath := strings.TrimPrefix(path, "/")
	if path != "/" {
		pathParts = strings.Split(strings.TrimSuffix(trimmedPath, "/"), "/")
	} else {
		pathParts = append(pathParts, "")
	}

	rootNode, rootExists := t.roots[pathParts[0]]
	isLastNode := len(pathParts) == 1
	if !rootExists {
		rootNode = &node{
			key:      pathParts[0],
			children: make(map[string]*node),
			handlers: make(map[string]routeConfig),
			end:      false,
		}

		t.roots[pathParts[0]] = rootNode
	}
	if isLastNode {
		rootNode.end = true
		rootNode.handlers[method] = config
	}

	if isLastNode {
		return nil
	}

	currNode := rootNode
	for idx := 1; idx < len(pathParts); idx++ {
		pathPart := pathParts[idx]
		childNode, childExists := currNode.children[pathPart]
		if !childExists {
			childNode = &node{
				key:      pathParts[idx],
				children: make(map[string]*node),
				handlers: make(map[string]routeConfig),
				end:      false,
			}

			currNode.children[pathPart] = childNode
		}

		currNode = childNode
	}

	currNode.end = true
	currNode.handlers[method] = config

	return nil
}

// getHandler returns the appropriate
func (t *routeTree) getHandler(path string, method string) (routeConfig, bool) {
	var empty routeConfig
	if path == "" {
		return empty, false
	}

	trimmedPath := strings.TrimPrefix(strings.TrimSuffix(path, "/"), "/")
	pathParts := strings.Split(trimmedPath, "/")

	current, ok := t.roots[pathParts[0]]
	if !ok {
		return empty, false
	}

	for idx := 1; idx < len(pathParts); idx++ {
		current, ok = current.children[pathParts[idx]]
		if !ok {
			return empty, false
		}
	}

	if !current.end {
		return empty, false
	}

	config, ok := current.handlers[method]
	if !ok {
		return empty, ok
	}

	return config, true
}
