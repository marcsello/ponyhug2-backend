package views

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/marcsello/ponyhug2-backend/db"
	"gitlab.com/MikeTTh/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"strings"
	"time"
)

const (
	loggerKey    = "l"
	playerKey    = "p"
	dbQueriesKey = "q"
)

type AccessLevel uint8

const (
	AccessLevelPublic AccessLevel = iota
	AccessLevelRegistered
	AccessLevelAdmin
)

const (
	AuthHeaderKeyPrefix    = "Key "
	AuthHeaderBearerPrefix = "Bearer "
)

func goodLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		// some evil middlewares may modify this value, so we store it
		path := ctx.Request.URL.Path

		subLogger := logger.With(
			zap.String("method", ctx.Request.Method),
			zap.String("path", path),
			zap.String("query", ctx.Request.URL.RawQuery),
			zap.String("ip", ctx.ClientIP()),
			zap.String("user-agent", ctx.Request.UserAgent()),
		)

		ctx.Set(loggerKey, subLogger)

		ctx.Next() // <- execute next thing in the chain
		end := time.Now()

		latency := end.Sub(start)

		completedRequestFields := []zapcore.Field{
			zap.Int("status", ctx.Writer.Status()),
			zap.Duration("latency", latency),
		}

		if len(ctx.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range ctx.Errors.Errors() {
				subLogger.Error(e, completedRequestFields...)
			}
		}

		subLogger.Info(fmt.Sprintf("%s %s served: %d", ctx.Request.Method, path, ctx.Writer.Status()), completedRequestFields...) // <- always print this
	}
}

func GetLoggerFromContext(ctx *gin.Context) *zap.Logger { // This one panics
	var logger *zap.Logger
	l, ok := ctx.Get(loggerKey)
	if !ok {
		panic("logger was not in context")
	}
	logger = l.(*zap.Logger)
	return logger
}

func enforceAccessMiddleware(level AccessLevel, userOnly bool) gin.HandlerFunc {

	adminKey := env.String("ADMIN_KEY", "")

	return func(ctx *gin.Context) {
		l := GetLoggerFromContext(ctx).With(zap.Uint8("level", uint8(level)), zap.Bool("userOnly", userOnly))

		if level == AccessLevelPublic {
			return // accept
		}

		// check auth header
		authHeader := ctx.Request.Header.Get("Authorization")

		if authHeader == "" {
			l.Warn("No auth header was provided")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var hasValidToken bool
		var hasValidKey bool
		var isAdmin bool

		if strings.HasPrefix(authHeader, AuthHeaderBearerPrefix) {
			tokenString := strings.TrimPrefix(authHeader, AuthHeaderBearerPrefix)

			playerID, err := validateJWT(l, tokenString)
			if err != nil {
				l.Warn("Could not validate user token!", zap.Error(err))
			} else {

				q := GetQueriesFromContext(ctx)

				var player db.Player
				player, err = q.GetPlayer(ctx, playerID)
				if err != nil {
					l.Error("Failure while querying player", zap.Error(err))
					ctx.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				hasValidToken = true
				isAdmin = player.IsAdmin
				ctx.Set(playerKey, &player) // don't store the entire player in the context, only a pointer to it
			}

		} else if strings.HasPrefix(authHeader, AuthHeaderKeyPrefix) { // key auth is easy

			keyString := strings.TrimPrefix(authHeader, AuthHeaderKeyPrefix)
			if adminKey != "" && keyString == adminKey {
				hasValidKey = true
				isAdmin = true
			} else {
				l.Debug("Admin key unset, or invalid key provided", zap.Bool("adminKeySet", adminKey != ""))
			}

		} else { // invalid auth type
			l.Warn("Invalid auth type")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var hasValidCreds bool
		if userOnly {
			hasValidCreds = hasValidToken
		} else {
			hasValidCreds = hasValidToken || hasValidKey
		}

		l.Debug("Authentication completed",
			zap.Bool("hasValidToken", hasValidToken),
			zap.Bool("hasValidKey", hasValidKey),
			zap.Bool("isAdmin", isAdmin),
			zap.Bool("hasValidCreds", hasValidCreds),
		)

		switch level {
		case AccessLevelRegistered: // accept key or token
			if !hasValidCreds {
				l.Warn("No valid creds presented")
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}

		case AccessLevelAdmin: // accept key or token of an admin
			if !hasValidCreds {
				l.Warn("No valid creds presented", zap.Bool("isAdmin", isAdmin), zap.Bool("hasValidCreds", hasValidCreds))
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			if !isAdmin {
				l.Warn("Not an admin", zap.Bool("isAdmin", isAdmin), zap.Bool("hasValidCreds", hasValidCreds))
				ctx.AbortWithStatus(http.StatusForbidden)
				return
			}
		default:
			panic("invalid level configured")
		}

	}
}

// GetPlayerFromContext retrives the player obect from the context. You should only set allowNil to false, when userOnly is set to true in enforceAccessMiddleware
func GetPlayerFromContext(ctx *gin.Context, allowNil bool) *db.Player {
	p, ok := ctx.Get(playerKey)
	if !ok {
		if allowNil {
			return nil
		} else {
			panic("no player in context")
		}
	}
	if p == nil && (!allowNil) {
		panic("nil is stored in place of player in context")
	}
	return p.(*db.Player) // this may panic
}

func withDBQueries(q *db.Queries) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(dbQueriesKey, q)
	}
}

func GetQueriesFromContext(ctx *gin.Context) *db.Queries {
	q, ok := ctx.Get(dbQueriesKey)
	if !ok {
		panic("no queries in context")
	}
	return q.(*db.Queries)
}

// injectUsernameToLogger should be only added after both authN and logger middleware were invoked, it only updates the stored logger to include the username
func injectPlayerToLogger(ctx *gin.Context) {
	l := GetLoggerFromContext(ctx)
	p := GetPlayerFromContext(ctx, true)
	if p != nil {
		ctx.Set(loggerKey, l.With(zap.String("player", p.Name)))
	}
}
