package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/WangWilly/labs-hr-go/pkgs/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

////////////////////////////////////////////////////////////////////////////////

type idKey struct{}

func IdFromCtx(ctx context.Context) (id string, ok bool) {
	switch ctx.Value(idKey{}).(type) {
	case string:
		id = ctx.Value(idKey{}).(string)
		ok = true
	case xid.ID:
		id = ctx.Value(idKey{}).(xid.ID).String()
		ok = true
	default:
		id = ""
		ok = false
	}
	return
}

func ctxWithID(ctx context.Context, id xid.ID) context.Context {
	return context.WithValue(ctx, idKey{}, id)
}

func ctxWithStringID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, idKey{}, id)
}

////////////////////////////////////////////////////////////////////////////////

func LoggingMiddleware() gin.HandlerFunc {

	return func(ginCtx *gin.Context) {
		logger := utils.GetLogger().With().
			Caller().Logger()

		logger.UpdateContext(func(logCtx zerolog.Context) zerolog.Context {
			if reqId := ginCtx.Request.Header.Get(utils.RequestIdHeader); reqId == "" {
				xidInst := xid.New()
				ginCtx.Request.Header.Set(utils.RequestIdHeader, xidInst.String())
				// Set the request ID in the context
				ginCtx.Request = ginCtx.Request.WithContext(
					ctxWithID(ginCtx.Request.Context(), xidInst),
				)
			} else {
				// If the request ID is already present in the header, use it
				// and set it in the context
				ginCtx.Request = ginCtx.Request.WithContext(
					ctxWithStringID(ginCtx.Request.Context(), reqId),
				)
			}

			////////////////////////////////////////////////////////////////////

			return logCtx.
				Str("url", ginCtx.Request.URL.String()).
				Str("method", ginCtx.Request.Method).
				Str("request_id", ginCtx.Request.Header.Get(utils.RequestIdHeader))
		})

		////////////////////////////////////////////////////////////////////////

		ginCtx.Request = ginCtx.Request.WithContext(
			logger.WithContext(ginCtx.Request.Context()),
		)

		////////////////////////////////////////////////////////////////////////

		tic := time.Now()
		ginCtx.Next()
		toc := time.Since(tic)

		////////////////////////////////////////////////////////////////////////
		// Log request details

		detailedLogger := utils.GetDetailedLogger().With().
			Str("url", ginCtx.Request.URL.String()).
			Str("method", ginCtx.Request.Method).
			Str("request_id", ginCtx.Request.Header.Get(utils.RequestIdHeader)).
			Str("client_ip", ginCtx.Request.RemoteAddr).
			Str("user_agent", ginCtx.Request.UserAgent()).
			Str("referer", ginCtx.Request.Referer()).
			Str("session_id", ginCtx.Request.Header.Get(utils.SessionIdHeader)).
			Str("latency", toc.String()).
			Logger()

		if ginCtx.Writer.Status() >= 400 {
			detailedLogger.Error().Msg(fmt.Sprintf("%s %s %d", ginCtx.Request.Method, ginCtx.Request.URL.String(), ginCtx.Writer.Status()))
		} else if ginCtx.Writer.Status() >= 200 {
			detailedLogger.Info().Msg(fmt.Sprintf("%s %s %d", ginCtx.Request.Method, ginCtx.Request.URL.String(), ginCtx.Writer.Status()))
		} else {
			detailedLogger.Debug().Msg(fmt.Sprintf("%s %s %d", ginCtx.Request.Method, ginCtx.Request.URL.String(), ginCtx.Writer.Status()))
		}
	}
}
