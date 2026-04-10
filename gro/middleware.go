package gro

import (
	"fmt"
	"log"
	"runtime"
	"strings"
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

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Println(trace(message))

			}
		}()
		c.Next()
	}
}

func trace(message string) string {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(3, pcs)

	frames := runtime.CallersFrames(pcs[:n])
	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for {
		frame, more := frames.Next()
		str.WriteString(fmt.Sprintf("%s\n\t%s:%d\n",
			frame.Function,
			frame.File,
			frame.Line,
		))
		if !more {
			break
		}
	}
	return str.String()
}

func TestLogin() HandlerFunc {
	return func(c *Context) {
		log.Println("TestLogin")
	}
}
