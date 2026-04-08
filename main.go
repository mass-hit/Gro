package main

import (
	"Gro/gro"
	"net/http"
)

func main() {
	r := gro.New()
	r.GET("/", func(context *gro.Context) {
		context.HTML(http.StatusOK, "<h1>Hello Gro</h1>")
	})
	r.GET("/hello", func(context *gro.Context) {
		context.String(http.StatusOK, "Hello Gro")
	})
	r.GET("/hellohello", func(context *gro.Context) {
		context.String(http.StatusOK, "Hello Hello Gro")
	})
	r.POST("/login", func(context *gro.Context) {
		context.JSON(http.StatusOK, gro.H{
			"username": context.PostForm("username"),
			"password": context.PostForm("password"),
		})
	})
	r.POST("/register/12", func(context *gro.Context) {
		context.String(http.StatusOK, "register 12")
	})
	r.GET("/hello/:name", func(context *gro.Context) {
		context.String(http.StatusOK, "hello %s", context.ParamMap["name"])
	})
	r.GET("/hello/:name/:age", func(context *gro.Context) {
		context.JSON(http.StatusOK, gro.H{
			"name": context.ParamMap["name"],
			"age":  context.ParamMap["age"],
		})
	})
	r.POST("/hello/*name", func(context *gro.Context) {
		context.String(http.StatusOK, "hello %s", context.ParamMap["name"])
	})
	r.GET("/helllo/123/:name/*age", func(context *gro.Context) {
		context.JSON(http.StatusOK, gro.H{
			"name": context.ParamMap["name"],
			"age":  context.ParamMap["age"],
		})
	})
	r.Run(":8080")
}
