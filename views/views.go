package views

import (
	"github.com/gin-gonic/gin"
	"github.com/marcsello/ponyhug2-backend/db"
	"go.uber.org/zap"
)

func InitViews(logger *zap.Logger, queries *db.Queries) *gin.Engine {

	logger.Info("Configuring routes...")

	r := gin.New()
	r.Use(goodLoggerMiddleware(logger))
	r.Use(withDBQueries(queries))

	meGroup := r.Group("/me")
	meGroup.Use(enforceAccessMiddleware(AccessLevelRegistered, true), injectPlayerToLogger)
	meGroup.GET("", getMe)
	meGroup.GET("/cards", getMyCards)
	meGroup.POST("/obtain", doObtain)

	playerGroup := r.Group("/players")
	playerGroup.POST("", enforceAccessMiddleware(AccessLevelPublic, false), registerPlayer)
	playerGroup.GET("", enforceAccessMiddleware(AccessLevelAdmin, false), injectPlayerToLogger, listPlayers)
	playerGroup.GET("/:id", enforceAccessMiddleware(AccessLevelAdmin, false), injectPlayerToLogger, getPlayer)
	playerGroup.PATCH("/:id", enforceAccessMiddleware(AccessLevelAdmin, false), injectPlayerToLogger, patchPlayer)

	cardsGroup := r.Group("/cards")
	cardsGroup.Use(enforceAccessMiddleware(AccessLevelAdmin, false), injectPlayerToLogger)
	cardsGroup.GET("", listAllCards)
	cardsGroup.POST("", createCard)
	cardsGroup.PUT(":id/wear/:level", upsertWearLevel)

	return r
}
