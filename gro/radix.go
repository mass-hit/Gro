package gro

import (
	"Gro/utils"
	"strings"
)

const (
	skind = iota
	pkind
	paramLabel = byte(':')
	slash      = byte('/')
	nilString  = ""
)

type Trees []*tree

func (trees Trees) get(method string) *tree {
	for _, tree := range trees {
		if tree.method == method {
			return tree
		}
	}
	return nil
}

func checkPathValid(path string) {
	utils.Assert(len(path) > 0, "path is empty")
	utils.Assert(path[0] == slash, "path must begin with '/'")
	pathLen := len(path)
	for i := 0; i < pathLen; i++ {
		if path[i] == paramLabel {
			if path[i-1] != slash {
				panic("param must after '/'")
			}
			if (i < pathLen-1 && path[i+1] == paramLabel) || i == (pathLen-1) {
				panic("param must be named with a non-empty name in path '" + path + "'")
			}
		}
	}
}

type tree struct {
	method string
	root   *node
}

// Parses the path and inserts it into the radix.
func (tree *tree) addRoute(path string, handler HandlerFunc) {
	checkPathValid(path)
	var params []string
	pathLen := len(path)
	for i := 0; i < pathLen; i++ {
		if path[i] == paramLabel {
			j := i + 1
			tree.insert(path[:i], nil, skind, nil)
			for ; i < pathLen && path[i] != slash; i++ {
			}
			params = append(params, path[j:i])
			path = path[:j] + path[i:]
			i, pathLen = j, len(path)
			if i == pathLen {
				tree.insert(path[:i], handler, pkind, params)
				return
			}
			tree.insert(path[:i], nil, pkind, params)
		}
	}
	// Insert static path if no params
	tree.insert(path, handler, skind, params)
}

// Insert a new route into the tree
func (tree *tree) insert(path string, handler HandlerFunc, kind int8, params []string) {
	currentNode := tree.root
	search := path
	for {
		searchLen := len(search)
		prefixLen := len(currentNode.prefix)
		minLen := searchLen
		if prefixLen < minLen {
			minLen = prefixLen
		}
		lcpLen := 0
		// Compute the Longest Common Prefix
		for ; lcpLen < minLen && search[lcpLen] == currentNode.prefix[lcpLen]; lcpLen++ {
		}
		if lcpLen < prefixLen {
			// Create a new child node
			n := newNode(currentNode.kind, currentNode.prefix[lcpLen:], currentNode.handler, currentNode, currentNode.children, currentNode.paramChild, currentNode.params)
			for _, child := range currentNode.children {
				child.parent = n
			}
			if currentNode.paramChild != nil {
				currentNode.paramChild.parent = n
			}
			// Update the current node
			currentNode.kind = skind
			currentNode.prefix = currentNode.prefix[:lcpLen]
			currentNode.handler = nil
			currentNode.children = []*node{n}
			if lcpLen == searchLen {
				// Set the handler to the current node
				currentNode.kind = kind
				currentNode.handler = handler
				currentNode.params = params
			} else {
				currentNode.children = append(currentNode.children, newNode(kind, search[lcpLen:], handler, currentNode, nil, nil, params))
			}
		} else if lcpLen < searchLen {
			// Continue search
			search = search[lcpLen:]
			if nextNode := currentNode.findChildWithLabel(search[0]); nextNode == nil {
				child := newNode(kind, search, handler, currentNode, nil, nil, params)
				if kind == skind {
					currentNode.children = append(currentNode.children, child)
				} else {
					currentNode.paramChild = child
				}
			} else {
				currentNode = nextNode
				continue
			}
		} else {
			// Node already exist
			if currentNode.handler != nil && handler != nil {
				panic("handlers are already registered for path '" + path + "'")
			}
			if handler != nil {
				currentNode.handler = handler
				currentNode.params = params
			}
		}
		return
	}
}

// Finds registered handler by path
func (tree *tree) find(path string, params *[]Param) (handler HandlerFunc) {
	var (
		currentNode = tree.root
		search      = path
		searchIndex = 0
		paramIndex  = 0
	)
	for {
		if currentNode.kind == skind {
			prefixLen := len(currentNode.prefix)
			if len(search) < prefixLen || search[:prefixLen] != currentNode.prefix {
				// Backtrack
				if currentNode = currentNode.parent; currentNode == nil {
					return
				}
				goto Param
			}
			search = search[prefixLen:]
			searchIndex += len(currentNode.prefix)
		}
		// End of path
		if search == nilString {
			handler = currentNode.handler
			break
		}
		// Static node match
		if child := currentNode.findChild(search[0]); child != nil {
			currentNode = child
			continue
		}
	Param:
		if child := currentNode.paramChild; child != nil {
			currentNode = child
			i := strings.Index(search, "/")
			if i == -1 {
				i = len(search)
			}
			// Expand params slice
			*params = (*params)[:(paramIndex + 1)]
			(*params)[paramIndex].Value = search[:i]
			paramIndex++
			// Move forward
			search = search[i:]
			searchIndex += i
			if search == nilString {
				handler = currentNode.handler
				break
			}
			continue
		}
		// Backtrack
		previous := currentNode
		currentNode = previous.parent
		if currentNode == nil {
			return
		}
		if previous.kind == skind {
			searchIndex -= len(previous.prefix)
		} else {
			paramIndex--
			searchIndex -= len((*params)[paramIndex].Value)
			*params = (*params)[:paramIndex]
		}
		search = path[:searchIndex]
	}
	// Fill parameter keys
	for i, param := range currentNode.params {
		(*params)[i].Key = param
	}
	return
}

func (n *node) findChild(char byte) *node {
	for _, i := range n.children {
		if i.prefix[0] == char {
			return i
		}
	}
	return nil
}

func (n *node) findChildWithLabel(char byte) *node {
	for _, i := range n.children {
		if i.prefix[0] == char {
			return i
		}
	}
	if char == paramLabel {
		return n.paramChild
	}
	return nil
}

type node struct {
	kind       int8
	prefix     string
	handler    HandlerFunc
	parent     *node
	children   []*node
	paramChild *node
	params     []string
}

func newNode(kind int8, prefix string, handler HandlerFunc, parent *node, children []*node, paramChild *node, params []string) *node {
	return &node{kind: kind, prefix: prefix, handler: handler, parent: parent, children: children, paramChild: paramChild, params: params}
}
