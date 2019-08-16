package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
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


	/*TODO: t, e := e.App.GetTrainByTrainID(trainID)
	  TODO: t.Seats[chosenType].Chosen = true
	  TODO: e.App.SaveInfoAboutTrain swith to id, e := SaveTrain(t)
	  TODO: e.App.SaveTrainIDToUser(userID, trainID)
	*/
	trainID, err := e.App.SaveInfoAboutTrain(trainID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err})
		ctx.Abort()
		return
	}

	err = e.App.SaveTrainInUser(userID, trainID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}
