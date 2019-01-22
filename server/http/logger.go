package http

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		status := c.Writer.Status()
		url := c.Request.URL
		host := c.Request.Host
		log.Printf("GIN: |Status - %v |Latency - %v |Host - %v |Url - %v |\n", status, latency, host, url)
	}
}
