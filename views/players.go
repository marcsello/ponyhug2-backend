package views

import (
	"github.com/gin-gonic/gin"
	"github.com/marcsello/ponyhug2-backend/db"
	"github.com/marcsello/ponyhug2-backend/model"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func registerPlayer(ctx *gin.Context) {
	l := GetLoggerFromContext(ctx)
	var params model.PlayerRegister

	err := ctx.BindJSON(&params)
	if err != nil {
		l.Warn("Could not bind registration params", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}

	q := GetQueriesFromContext(ctx)

	var player db.Player
	player, err = q.CreatePlayer(ctx, params.Name)
	if err != nil {
		l.Error("Failure while creating new player", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	var token string
	token, err = generateToken(player.ID)
	if err != nil {
		l.Error("Failure while generating new token", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	resp := model.PlayerRegistrationSuccess{
		Name:  player.Name,
		Token: token,
	}

	ctx.JSON(http.StatusOK, resp)

}

func getAnyPlayer(ctx *gin.Context) {
	l := GetLoggerFromContext(ctx)
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 16)
	if err != nil {
		l.Warn("Invalid id provided", zap.String("idStr", idStr))
		ctx.Status(http.StatusBadRequest)
		return
	}

	q := GetQueriesFromContext(ctx)

	var p db.Player
	p, err = q.GetPlayer(ctx, int16(id))

	ctx.JSON(http.StatusOK, model.PlayerDataFromDB(p))
}

func listPlayers(ctx *gin.Context) {
	l := GetLoggerFromContext(ctx)
	q := GetQueriesFromContext(ctx)
	players, err := q.GetPlayers(ctx)
	if err != nil {
		l.Error("Error while querying players", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, model.PlayersDataFromDB(players))
}
