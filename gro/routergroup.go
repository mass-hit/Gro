package gro

import (
	"net/http"
	"path"
)

type RouterGroup struct {
	handlers []HandlerFunc
	basePath string
	engine   *Engine
}

// Use adds middleware to the group
func (group *RouterGroup) Use(middleware ...HandlerFunc) {
	group.handlers = append(group.handlers, middleware...)
}

func (group *RouterGroup) handle(method, relativePath string, handlers []HandlerFunc) {
	group.engine.addRoute(method, group.calculateAbsolutePath(relativePath), group.mergeHandlers(handlers))
}

// Group creates a new router group.
func (group *RouterGroup) Group(relativePath string, middleware ...HandlerFunc) *RouterGroup {
	return &RouterGroup{handlers: group.mergeHandlers(middleware), basePath: group.calculateAbsolutePath(relativePath), engine: group.engine}
}

func (group *RouterGroup) GET(pattern string, handlers ...HandlerFunc) {
	group.handle(http.MethodGet, pattern, handlers)
}

func (group *RouterGroup) POST(pattern string, handlers ...HandlerFunc) {
	group.handle(http.MethodPost, pattern, handlers)
}

func (group *RouterGroup) mergeHandlers(handlers []HandlerFunc) []HandlerFunc {
	mergedHandlers := make([]HandlerFunc, 0, len(group.handlers)+len(handlers))
	mergedHandlers = append(mergedHandlers, group.handlers...)
	mergedHandlers = append(mergedHandlers, handlers...)
	return mergedHandlers
}

func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	if relativePath == "" {
		return group.basePath
	}
	finalPath := path.Join(group.basePath, relativePath)
	if relativePath[len(relativePath)-1] == '/' && finalPath[len(finalPath)-1] != '/' {
		return finalPath + "/"
	}
	return finalPath
}
