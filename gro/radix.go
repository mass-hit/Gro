package gro

type Trees []*tree

func (trees Trees) get(method string) *tree {
	for _, tree := range trees {
		if tree.method == method {
			return tree
		}
	}
	return nil
}

type tree struct {
	method string
	root   *node
}

// Insert a new route into the tree
func (tree *tree) insert(path string, handler HandlerFunc) {
	currentNode := tree.root
	for {
		pathLen := len(path)
		prefixLen := len(currentNode.prefix)
		minLen := pathLen
		if prefixLen < minLen {
			minLen = prefixLen
		}
		lcpLen := 0
		// Compute the Longest Common Prefix
		for ; lcpLen < minLen && path[lcpLen] == currentNode.prefix[lcpLen]; lcpLen++ {
		}
		if lcpLen < prefixLen {
			// Create a new child node
			child := newNode(currentNode.prefix[lcpLen:], currentNode.handler, currentNode.children)
			// Update the current node
			currentNode.prefix = currentNode.prefix[:lcpLen]
			currentNode.handler = nil
			currentNode.children = []*node{child}
			if lcpLen == pathLen {
				// Set the handler to the current node
				currentNode.handler = handler
			} else {
				currentNode.children = append(currentNode.children, newNode(path[lcpLen:], handler, nil))
			}
		} else if lcpLen < pathLen {
			// Continue search
			path = path[lcpLen:]
			if nextNode := currentNode.findChild(path[0]); nextNode == nil {
				currentNode.children = append(currentNode.children, newNode(path, handler, nil))
			} else {
				currentNode = nextNode
				continue
			}
		} else {
			// Node already exist
			if currentNode.handler != nil && handler != nil {
				panic("handlers are already registered for path '" + path + "'")
			}
			currentNode.handler = handler
		}
		return
	}
}

// Finds registered handler by path
func (tree *tree) find(path string) HandlerFunc {
	currentNode := tree.root
	for {
		prefixLen := len(currentNode.prefix)
		if len(path) < prefixLen || path[:prefixLen] != currentNode.prefix {
			return nil
		}
		path = path[prefixLen:]
		if path == "" {
			return currentNode.handler
		}
		currentNode = currentNode.findChild(path[0])
		if currentNode == nil {
			return nil
		}
	}
}

func (n *node) findChild(char byte) *node {
	for _, i := range n.children {
		if i.prefix[0] == char {
			return i
		}
	}
	return nil
}

type node struct {
	prefix   string
	handler  HandlerFunc
	children []*node
}

func newNode(prefix string, handler HandlerFunc, children []*node) *node {
	return &node{prefix: prefix, handler: handler, children: children}
}
