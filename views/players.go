package views

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/marcsello/ponyhug2-backend/db"
	"github.com/marcsello/ponyhug2-backend/db_utils"
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

	err = params.Validate()
	if err != nil {
		l.Warn("Could not validate registration params", zap.Error(err))
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}

	q := GetQueriesFromContext(ctx)

	var player db.Player
	player, err = q.CreatePlayer(ctx, params.Name)
	if err != nil {
		if db_utils.IsDuplicatedKeyErr(err) {
			l.Warn("Name already in use", zap.String("The name that's already in use", params.Name), zap.Error(err))
			ctx.Status(http.StatusConflict)
			return
		}
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

func getPlayer(ctx *gin.Context) {
	l := GetLoggerFromContext(ctx)
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		l.Warn("Invalid id provided", zap.String("idStr", idStr))
		ctx.Status(http.StatusBadRequest)
		return
	}

	q := GetQueriesFromContext(ctx)

	var p db.Player
	p, err = q.GetPlayer(ctx, int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			l.Info("No such player", zap.Int32("id", int32(id)), zap.Error(err))
			ctx.Status(http.StatusNotFound)
			return
		}
		l.Error("Error while querying player", zap.Int32("id", int32(id)), zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

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

func patchPlayer(ctx *gin.Context) {
	l := GetLoggerFromContext(ctx)
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		l.Warn("Invalid id provided", zap.String("idStr", idStr))
		ctx.Status(http.StatusBadRequest)
		return
	}

	var params model.PatchPlayerParams

	err = ctx.BindJSON(&params)
	if err != nil {
		l.Warn("Could not bind player patch params", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = params.Validate()
	if err != nil {
		l.Warn("Could not validate player patch params", zap.Error(err))
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}

	q := GetQueriesFromContext(ctx)
	if params.IsAdmin {
		err = q.PromotePlayer(ctx, int32(id))
	} else {
		err = q.DemotePlayer(ctx, int32(id))
	}

	if err != nil {
		l.Error("Error while promoting/demoting player", zap.Bool("isAdmin", params.IsAdmin), zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)

}
