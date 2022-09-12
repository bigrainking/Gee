package geenote

import (
	"log"
	"time"
)

// 这是一个中间件
func Logger() HandlerFunc {
	return func(c *Context) {
		// Next前
		t := time.Now()

		c.Next()

		// Next后
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
