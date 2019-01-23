package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func Logger(logChan chan string) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		status := c.Writer.Status()
		url := c.Request.URL
		host := c.Request.Host
		logChan <- fmt.Sprintf("GIN: Status - %v |Latency - %v |Host - %v |Url - %v ", status, latency, host, url)
	}
}
