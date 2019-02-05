package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rzd/app/entity"
)

func (e *EventLayer) SaveOneTrain(ctx *gin.Context) {
	trainID := ctx.PostForm("train_id")
	userID := ctx.PostForm("user_id")
	userName := ctx.PostForm("user_name")

	trainID, err := e.App.SaveInfoAboutTrain(trainID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err})
		ctx.Abort()
		return
	}

	err = e.App.SaveTrainInUser(entity.User{
		UserTelegramID: userID,
		UserName:       userName,
	}, trainID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}
