package main

import (
	"context"
	_ "embed"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/marcsello/ponyhug2-backend/db"
	"github.com/marcsello/ponyhug2-backend/views"
	"gitlab.com/MikeTTh/env"
	"go.uber.org/zap"
)

//go:embed schema.sql
var ddl string

var (
	// injected build time
	version        string
	commitHash     string
	buildTimestamp string
)

func main() {
	debug := env.Bool("DEBUG", false)
	var lgr *zap.Logger
	if debug {
		lgr = zap.Must(zap.NewDevelopment())
		gin.SetMode(gin.DebugMode)
		lgr.Warn("RUNNING IN DEBUG MODE!")
	} else {
		lgr = zap.Must(zap.NewProduction())
		gin.SetMode(gin.ReleaseMode)
	}
	defer lgr.Sync()

	lgr.Info("Starting PonyHug2 (working title) server...", zap.String("version", version), zap.String("commitHash", commitHash), zap.String("buildTimestamp", buildTimestamp))

	lgr.Info("Init db stuff...")
	conn, err := pgx.Connect(context.TODO(), env.StringOrPanic("DATABASE_URL"))
	if err != nil {
		lgr.Error("Failed to connect to DB", zap.Error(err))
		panic(err)
	}
	defer conn.Close(context.TODO())

	// create schema
	_, err = conn.Exec(context.TODO(), ddl)
	if err != nil {
		lgr.Error("Failed to create schema", zap.Error(err))
		panic(err)
	}

	lgr.Info("Init views...")
	r := views.InitViews(lgr, db.New(conn))

	lgr.Info("Starting...")
	err = r.Run(env.String("BIND_ADDR", ":8080"))
	if err != nil {
		panic(err)
	}
}
