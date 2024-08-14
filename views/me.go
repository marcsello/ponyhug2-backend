package views

import (
	"github.com/gin-gonic/gin"
	"github.com/marcsello/ponyhug2-backend/model"
	"net/http"
)

func getMyCards(ctx *gin.Context) {

}

func getMe(ctx *gin.Context) {
	l := GetLoggerFromContext(ctx)
	p := GetPlayerFromContext(ctx)
	if p == nil {
		l.Error("Could not get player from context!")
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, model.PlayerSelfFromDB(*p))
}
