package gro

import (
	"net/http"
)

type HandlerFunc func(*Context)

type Engine struct {
	trees Trees
}

func New() *Engine {
	return &Engine{trees: make(Trees, 0, 4)}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	// Find root of the tree for the given HTTP method
	methodTree := engine.trees.get(method)
	// If no tree exists for this method, create a new one
	if methodTree == nil {
		methodTree = &tree{method: method, root: &node{}}
		engine.trees = append(engine.trees, methodTree)
	}
	if pattern == "" {
		panic("empty path")
	}
	if pattern[0] != '/' {
		panic("path must begin with '/'")
	}
	// Insert the route into the radix tree
	methodTree.insert(pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(w, req)
	if methodTree := engine.trees.get(req.Method); methodTree != nil {
		// Find handler for the given path
		handler := methodTree.find(c.Path)
		if handler != nil {
			handler(c)
			return
		}
	}
	// If no route matches the request
	c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
}
