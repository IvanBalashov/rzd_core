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
		logChan <- fmt.Sprintf("Http_Server: Status - %3v |Latency - %6v |Host - %10v |Url - %40v ", status, latency, host, url)
	}
}
