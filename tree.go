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
	handler    HttpHandler
	pathParams map[string]string
}

type node struct {
	key      string
	children map[string]*node
	handlers map[string]routeConfig // For every http method, there can be a handler
	end      bool                   // denotes if this node specifies the end of a valid route
	isParam  bool
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
			isParam:  isPathParam(pathParts[0]),
		}

		t.roots[pathParts[0]] = rootNode
	}
	if isLastNode {
		rootNode.end = true
		rootNode.handlers[method] = config
	}

	if isPathParam(pathParts[0]) {
		rootNode.isParam = true
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
				isParam:  isPathParam(pathParts[idx]),
			}

			currNode.children[pathPart] = childNode

			if isPathParam(pathPart) {
				rootNode.isParam = true
			}
		}

		currNode = childNode
	}

	currNode.end = true
	currNode.handlers[method] = config

	return nil
}

// getConfig returns the appropriate
func (t *routeTree) getConfig(path string, method string) (routeConfig, bool) {
	var empty routeConfig
	if path == "" {
		return empty, false
	}

	trimmedPath := strings.TrimPrefix(strings.TrimSuffix(path, "/"), "/")
	pathParts := strings.Split(trimmedPath, "/")

	pathParams := make(map[string]string)

	current, ok := t.roots[pathParts[0]]
	if !ok {
		current, ok = t.findPathParam(t.roots)
		if !ok {
			return empty, false
		}
		pathParams[extractParam(current.key)] = pathParts[0]
	}

	for idx := 1; idx < len(pathParts); idx++ {
		children := current.children
		current, ok = children[pathParts[idx]]
		if !ok {
			// if there's no exact match, maybe it's a path param
			current, ok = t.findPathParam(children)
			if !ok {
				return empty, false
			}

			pathParams[extractParam(current.key)] = pathParts[idx]
		}
	}

	if !current.end {
		return empty, false
	}

	config, ok := current.handlers[method]
	if !ok {
		return empty, ok
	}

	config.pathParams = pathParams

	return config, true
}

func (t *routeTree) findPathParam(nodes map[string]*node) (*node, bool) {
	for _, node := range nodes {
		if node.isParam {
			return node, true
		}
	}

	return nil, false
}

func isPathParam(path string) bool {
	return strings.HasPrefix(path, "{") && strings.HasSuffix(path, "}")
}

func extractParam(pathPart string) string {
	return strings.Trim(pathPart, "{}")
}
