package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rzd/app/usecase"
	"rzd/server/http/middleware"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func NewHandler(app *usecase.App) http.Handler {
	handler := gin.Default()

	middlewares := middleware.InitMiddleWares(app)
	api := handler.Group("/api/v1")
	api.GET("health", middlewares.Health)

	api.GET("test", middlewares.Test)

	return handler
}
