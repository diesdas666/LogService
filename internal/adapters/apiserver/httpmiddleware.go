package apiserver

import (
	"example_consumer/internal/core/app"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"time"
)

const httpLogFormat = `"[END] %s %s %s" from %s`

// zapLoggerMiddleware provides middleware for adding zap logger into the context of request handler
func zapLoggerMiddleware(logger *zap.SugaredLogger) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		fn := func(c echo.Context) error {
			req := c.Request()
			reqID := c.Response().Header().Get(echo.HeaderXRequestID)
			l := logger.With(zap.String("requestId", reqID))
			ctx := app.ContextWithLogger(req.Context(), l)
			c.SetRequest(c.Request().WithContext(ctx))
			tbegin := time.Now()
			defer func() {
				l.With(zap.String("duration", time.Since(tbegin).String())).Infof(httpLogFormat,
					req.Method,
					req.URL.Path,
					req.Proto,
					req.RemoteAddr,
				)
			}()
			return next(c)
		}
		return fn
	}
}
