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
	meGroup.GET("", enforceAccessMiddleware(AccessLevelRegistered, true), injectPlayerToLogger, getMe)
	meGroup.GET("/cards", enforceAccessMiddleware(AccessLevelRegistered, true), injectPlayerToLogger, getMyCards)

	playerGroup := r.Group("/player")
	playerGroup.POST("", enforceAccessMiddleware(AccessLevelPublic, false), registerPlayer)
	playerGroup.GET("", enforceAccessMiddleware(AccessLevelAdmin, false), injectPlayerToLogger, listPlayers)
	playerGroup.GET("/:id", enforceAccessMiddleware(AccessLevelAdmin, false), injectPlayerToLogger, getAnyPlayer)

	return r
}
