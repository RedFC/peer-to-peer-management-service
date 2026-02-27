package middlewares

import (
	"p2p-management-service/utils"

	"github.com/kataras/iris/v12"
)

func LoggerMiddleware(ctx iris.Context) {
	// before request
	ctx.Next() // run handler

	// after handler executes → log status
	status := ctx.GetStatusCode()
	if status >= 400 {
		utils.Logger.Error().
			Int("status", status).
			Str("method", ctx.Method()).
			Str("path", ctx.Path()).
			Str("ip", ctx.RemoteAddr()).
			Str("error", ctx.Values().GetString("error")).
			Msg("request failed")
	} else {
		utils.Logger.Info().
			Int("status", status).
			Str("method", ctx.Method()).
			Str("path", ctx.Path()).
			Str("ip", ctx.RemoteAddr()).
			Msg("request succeeded")
	}
}

func RecoveryMiddleware(ctx iris.Context) {
	defer func() {
		if r := recover(); r != nil {
			utils.Logger.Error().
				Str("method", ctx.Method()).
				Str("path", ctx.Path()).
				Interface("recover", r).
				Msg("panic recovered")
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"message": "Internal Server Error"})
		}
	}()
	ctx.Next()
}
