package gro

import (
	"net/http"
	"path"
)

type RouterGroup struct {
	basePath string
	engine   *Engine
}

func (group *RouterGroup) handle(method, relativePath string, handler HandlerFunc) {
	absolutePath := group.calculateAbsolutePath(relativePath)
	group.engine.addRoute(method, absolutePath, handler)
}

// Group creates a new router group.
func (group *RouterGroup) Group(relativePath string) *RouterGroup {
	return &RouterGroup{basePath: group.calculateAbsolutePath(relativePath), engine: group.engine}
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.handle(http.MethodGet, pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.handle(http.MethodPost, pattern, handler)
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
