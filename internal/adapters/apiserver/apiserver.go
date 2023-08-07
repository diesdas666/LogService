package apiserver

import (
	"context"
	"example_consumer/internal/adapters/apiserver/internal"
	"example_consumer/internal/core/di"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Start(_ context.Context, di *di.DI) {

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(zapLoggerMiddleware(zap.S()))

	apiRoutes(e, di)

	listenAddr := fmt.Sprintf("%s:%d", di.Config.Server.Addr, di.Config.Server.Port)
	go func() {
		if err := e.Start(listenAddr); err != nil && err != http.ErrServerClosed {
			zap.S().Fatal("shutting down the server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	zap.S().Infof("Server is shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		zap.S().Warnf("Could not gracefully shutdown the server")
	}
	zap.S().Infof("Server stopped")

	zap.S().Info("Cleaning up resources")
	di.Close()
	zap.S().Infof("Resources has been cleaned up")
}

func apiRoutes(e *echo.Echo, di *di.DI) {
	e.GET("/api/version", internal.GetVersion())
	contacts := e.Group("/api/contacts")
	contacts.POST("", internal.CreateContact(di.UseCases))
	contacts.GET("", internal.ListAllContacts(di.UseCases))
	contacts.PUT("/:id", internal.UpdateContact(di.UseCases))
	contacts.GET("/:id", internal.GetContact(di.UseCases))
	contacts.DELETE("/:id", internal.DeleteContact(di.UseCases))
}
