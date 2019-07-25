package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rzd/app/entity"
)

func (e *EventLayer) SaveOneTrain(ctx *gin.Context) {
	var trainID string
	if val := ctx.PostForm("train_id"); val == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": "empty train_id"})
		ctx.Abort()
	} else {
		trainID = val
	}

	var userID string
	if val := ctx.PostForm("user_id"); val == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": "empty user_id"})
		ctx.Abort()
	} else {
		userID = val
	}

	var userName string
	if val := ctx.PostForm("user_name"); val == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": "empty user_name"})
		ctx.Abort()
	} else {
		userName = val
	}

	trainID, err := e.App.SaveInfoAboutTrain(trainID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err})
		ctx.Abort()
		return
	}

	user := &entity.User{
		UserTelegramID: userID,
		UserName:       userName,
	}
	err = e.App.SaveTrainInUser(user, trainID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}
