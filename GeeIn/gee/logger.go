package gee

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		//Start timer
		t := time.Now()
		// Process request
		//c.Fail(500, "Internal Server Error")
		// Calculatie resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
