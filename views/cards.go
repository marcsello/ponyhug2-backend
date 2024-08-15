package views

import (
	"github.com/gin-gonic/gin"
	"github.com/marcsello/ponyhug2-backend/db"
	"github.com/marcsello/ponyhug2-backend/db_utils"
	"github.com/marcsello/ponyhug2-backend/model"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func listAllCards(ctx *gin.Context) {
	l := GetLoggerFromContext(ctx)
	q := GetQueriesFromContext(ctx)

	cards, err := q.GetCardBases(ctx)
	if err != nil {
		l.Error("Error while listing card bases", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, model.CardBasesForAdminsFromDB(cards))
}

func createCard(ctx *gin.Context) {
	l := GetLoggerFromContext(ctx)

	var params model.CreateCardBaseParams

	err := ctx.BindJSON(&params)
	if err != nil {
		l.Warn("Could not bind card params", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = params.Validate()
	if err != nil {
		l.Warn("Could not validate card params", zap.Error(err))
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}

	q := GetQueriesFromContext(ctx)

	base, err := q.CreateCardBase(ctx, db.CreateCardBaseParams{
		Key:    params.Key,
		Name:   params.Name,
		Source: params.Source,
		Place:  params.Place,
	})
	if err != nil {
		if db_utils.IsDuplicatedKeyErr(err) {
			l.Warn("Duplicated key while creating card", zap.Error(err))
			ctx.Status(http.StatusConflict)
			return
		}
		l.Error("Error while listing card bases", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(
		http.StatusCreated,
		model.BareCardBaseFromDBCardBase(base),
	)

}

func upsertWearLevel(ctx *gin.Context) {
	l := GetLoggerFromContext(ctx)

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil || id < 0 {
		l.Warn("Invalid id provided", zap.String("idStr", idStr))
		ctx.Status(http.StatusBadRequest)
		return
	}

	levelStr := ctx.Param("level")
	level, err := strconv.ParseInt(levelStr, 10, 16)
	if err != nil || level < 0 {
		l.Warn("Invalid id provided", zap.String("levelStr", levelStr))
		ctx.Status(http.StatusBadRequest)
		return
	}

	var params model.AssignWearLevelParams

	err = ctx.BindJSON(&params)
	if err != nil {
		l.Warn("Could not bind wear level params", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = params.Validate()
	if err != nil {
		l.Warn("Could not validate wear level params", zap.Error(err))
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}

	q := GetQueriesFromContext(ctx)

	levelData, err := q.AssignCardImageToWearLevel(ctx, db.AssignCardImageToWearLevelParams{
		BaseID:    int16(id),
		WearLevel: int16(level),
		ImageUrl:  params.ImgUrl,
	})
	if err != nil {
		l.Error("Error while upserting card img", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, model.CardWearImgFromDB(levelData))

}
