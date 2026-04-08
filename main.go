package main

import (
	"Gro/gro"
	"net/http"
)

func main() {
	r := gro.New()
	hello := r.Group("/hello")
	hello.GET("/gro", func(context *gro.Context) {
		context.String(http.StatusOK, "Hello Gro")
	})
	hello.GET("/:name", func(context *gro.Context) {
		context.String(http.StatusOK, "hello %s", context.ParamMap["name"])
	})
	hello.POST("/*name", func(context *gro.Context) {
		context.String(http.StatusOK, "hello %s", context.ParamMap["name"])
	})
	login := r.Group("/login")
	login.POST("/", func(context *gro.Context) {
		context.JSON(http.StatusOK, gro.H{
			"username": context.PostForm("username"),
			"password": context.PostForm("password"),
		})
	})
	login.GET("/:name/:age", func(context *gro.Context) {
		context.JSON(http.StatusOK, gro.H{
			"name": context.ParamMap["name"],
			"age":  context.ParamMap["age"],
		})
	})
	r.Run(":8080")
}
