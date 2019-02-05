package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rzd/app/usecase"
	"rzd/server/http/middleware"
	"time"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// TODO: New middleware for gin logging.
func NewHandler(app *usecase.App) http.Handler {
	handler := gin.New()
	handler.Use(Logger(app.LogChan))

	eventLayer := middleware.NewEventLayer(app)
	api := handler.Group("/api/v1")

	api.GET("health", eventLayer.Health)

	api.GET("trains_list", eventLayer.GetAllTrains)

	api.GET("users_count", eventLayer.UsersCount)
	// FIXME: rewrite to post
	api.GET("new_user", eventLayer.NewUser)

	api.POST("save_one_train", eventLayer.GetAllTrains)


	return handler
}

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
