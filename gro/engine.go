package gro

import (
	"Gro/utils"
	"net/http"
)

type HandlerFunc func(*Context)

type Engine struct {
	trees     Trees
	maxParams uint16
	RouterGroup
}

func New() *Engine {
	engine := &Engine{trees: make(Trees, 0, 4), maxParams: 0}
	engine.RouterGroup.engine = engine
	return engine
}

func (engine *Engine) addRoute(method string, path string, handlers []HandlerFunc) {
	// Validate input
	utils.Assert(len(path) > 0, "path is empty")
	utils.Assert(path[0] == '/', "path must begin with '/'")
	utils.Assert(method != "", "method is empty")
	utils.Assert(handlers != nil && len(handlers) > 0, "handler can not be nil")
	// Find root of the tree for the given HTTP method
	methodTree := engine.trees.get(method)
	// If no tree exists for this method, create a new one
	if methodTree == nil {
		methodTree = &tree{method: method, root: &node{}}
		engine.trees = append(engine.trees, methodTree)
	}
	// Insert the route into the radix tree
	methodTree.addRoute(path, handlers)
	// Update maximum parameter count
	if paramCount := countParams(path); paramCount > engine.maxParams {
		engine.maxParams = paramCount
	}
}

// Counts the number of parameter in a path
func countParams(path string) uint16 {
	var n uint16
	for i := 0; i < len(path); i++ {
		c := path[i]
		if c == ':' {
			n++
		}
	}
	return n
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(w, req)
	if methodTree := engine.trees.get(req.Method); methodTree != nil {
		params := make([]Param, 0, engine.maxParams)
		// Find handler for the given path
		handlers := methodTree.find(c.Path, &params)
		if handlers != nil {
			c.ParamMap = convertToParamMap(params)
			c.Handlers = handlers
			c.Next()
			return
		}
	}
	// If no route matches the request
	c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
}

// Converts a slice of Param into a map
func convertToParamMap(params []Param) map[string]string {
	if params == nil || len(params) == 0 {
		return nil
	}
	paramMap := make(map[string]string, len(params))
	for _, param := range params {
		paramMap[param.Key] = param.Value
	}
	return paramMap
}
