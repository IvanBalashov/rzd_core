package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rzd/app/entity"
)

func (e *EventLayer) NewUser(ctx *gin.Context) {
	userID := ctx.PostForm("user_id")

	userName := ctx.PostForm("user_name")

	user := &entity.User{
		UserTelegramID: userID,
		UserName:       userName,
		Notify:         true,
	}
	ok, err := e.App.AddUser(user)
	if err != nil {
		if ok {
			ctx.JSON(http.StatusOK, gin.H{"status": "user_exist"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "created"})
}
