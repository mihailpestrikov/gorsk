package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo"
	"github.com/ribice/gorsk/pkg/utl/config"
	"github.com/ribice/gorsk/pkg/utl/zlog"
)

// New creates new HTTP server
func New(cfg *config.Config, log *zlog.Logger) *echo.Echo {
	e := echo.New()
	e.Server.ReadTimeout = time.Second * time.Duration(cfg.Server.ReadTimeout)
	e.Server.WriteTimeout = time.Second * time.Duration(cfg.Server.WriteTimeout)
	e.Server.Addr = fmt.Sprintf(":%d", cfg.Server.Port)

	go func() {
		if err := e.Start(e.Server.Addr); err != nil {
			log.Info("Shutting down the server", err.Error())
			// e.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.Server.ShutdownTimeout))
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err.Error())
	}
	return e
}

// Health responds with a simple OK status. (New function)
func Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
