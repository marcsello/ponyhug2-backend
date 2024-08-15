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
)

func getMyCards(ctx *gin.Context) {
	l := GetLoggerFromContext(ctx)
	p := GetPlayerFromContext(ctx, false)
	q := GetQueriesFromContext(ctx)

	cards, err := q.GetPlayerCards(ctx, p.ID)
	if err != nil {
		l.Error("Failure while querying player's cards", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(
		http.StatusOK,
		model.CardCopiesVisibleByPlayerFromDBPlayerCardsRows(cards),
	)
}

func doObtain(ctx *gin.Context) {
	l := GetLoggerFromContext(ctx)

	var params model.PlayerObtainCard
	err := ctx.BindJSON(&params)
	if err != nil {
		l.Warn("Could not bind obtain params", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = params.Validate()
	if err != nil {
		l.Warn("Could not validate obtain params", zap.Error(err))
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}

	p := GetPlayerFromContext(ctx, false)
	q := GetQueriesFromContext(ctx)

	var newCopy db.CardCopy
	switch len(params.Key) {
	case 10:
		newCopy, err = q.MakeSubsequentCopy(ctx, db.MakeSubsequentCopyParams{
			PlayerID: p.ID,
			Key:      params.Key,
		})
	case 9:
		newCopy, err = q.MakeFirstCopy(ctx, db.MakeFirstCopyParams{
			PlayerID: p.ID,
			Key:      &params.Key,
		})
	default:
		l.Warn("Invalid key length", zap.Error(err), zap.Int("keyLen", len(params.Key)))
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}

	if err != nil {
		if db_utils.IsDuplicatedKeyErr(err) {
			l.Info("This player already has this card", zap.Error(err))
			ctx.Status(http.StatusConflict)
			return
		}
		if db_utils.IsNotNullViolation(err) { // our magic query would try to use null value when there is no such card
			l.Info("This player tried to copy non-existing code", zap.Error(err))
			ctx.Status(http.StatusNotFound)
			return
		}
		if errors.Is(err, pgx.ErrNoRows) {
			// I have no idea why exactly this happens
			// But I think it is because the player already has a higher level of this card
			l.Info("Player card copy query returned no rows.... what now???", zap.Error(err))
			ctx.Status(http.StatusNoContent)
			return
		}
		l.Error("Failure while processing card obtainment", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	l.Info("Player copied a new card",
		zap.Int32("CopyID", newCopy.ID),
		zap.Int16("BaseID", newCopy.BaseID),
		zap.Int16("WearLevel", newCopy.WearLevel),
	)

	// And this is where I gave up....
	newCopyDetails, err := q.GetCardCopy(ctx, newCopy.ID)
	if err != nil {
		l.Error("Failure while processing card obtainment", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(
		http.StatusOK,
		model.CardCopyVisibleByPlayerFromDBGetCardCopyRow(newCopyDetails),
	)

}

func getMe(ctx *gin.Context) {
	p := GetPlayerFromContext(ctx, false)
	ctx.JSON(http.StatusOK, model.PlayerSelfFromDB(*p))
}
