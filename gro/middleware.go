package gro

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		log.Printf("request :%s", c.Req.RequestURI)
		c.Next()
		time.Sleep(100 * time.Millisecond)
		log.Printf("use time :%v", time.Since(t))
	}
}

func TestLogin() HandlerFunc {
	return func(c *Context) {
		log.Println("TestLogin")
	}
}
