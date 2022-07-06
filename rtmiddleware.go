package trex

import (
	"time"

	"github.com/gin-gonic/gin"
)

type ReqTrFunc func(
	context interface{},
	method string,
	path string,
	query string,
	agent string,
	ip string,
	status int,
	bytes int,
	start time.Time,
	end time.Time,
)

// Injects a request tracer middleware into the gin context
// It collects information from the request and response, and calls
// a parameter supplied handler function with the information
// The function is called with a transaction context of interface{} type
func RequestTracerMiddleware(trf ReqTrFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		txc := c.MustGet("tx-context")

		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		agent := c.Request.UserAgent()
		ip := c.ClientIP()

		start := time.Now()
		c.Next()
		end := time.Now()

		status := c.Writer.Status()
		size := c.Writer.Size()

		trf(txc, method, path, query, agent, ip, status, size, start, end)
	}
}
