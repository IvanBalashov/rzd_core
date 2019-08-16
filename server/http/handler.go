package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rzd/app/usecase"
	"rzd/server/http/handlers"
	"strconv"
	"time"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// TODO: New handlers for gin logging.
func NewHandler(app *usecase.App) http.Handler {
	handler := gin.New()
	handler.Use(Logger(app.LogChan))

	eventLayer := handlers.NewEventLayer(app)
	api := handler.Group("/api/v1")

	api.GET("health", eventLayer.Health)

	api.GET("trains_list", eventLayer.GetAllTrains)

	api.GET("users_count", eventLayer.UsersCount)
	// FIXME: rewrite to post
	api.POST("new_user", eventLayer.NewUser)

	api.POST("save_one_train", eventLayer.SaveOneTrain)

	return handler
}

func Logger(logChan chan string) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := strconv.FormatFloat(time.Since(t).Seconds()*100, 'f', 10, 64)
		status := c.Writer.Status()
		url := c.Request.URL
		host := c.Request.Host
		method := c.Request.Method
		logChan <- fmt.Sprintf("Http_Server: Method - %6v |Status - %3v |Latency - %3v sec|Host - %10v |Url - %40v ", method, status, latency[:6], host, url)
	}
}
