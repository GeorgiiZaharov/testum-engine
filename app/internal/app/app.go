package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type App struct {
	server *http.Server
	logger *zap.Logger
}

func New(server *http.Server, logger *zap.Logger) *App {
	return &App{
		server: server,
		logger: logger,
	}
}

func (a *App) Run() error {
	go func() {
		a.logger.Info("http server started", zap.String("addr", a.server.Addr))

		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal("server crashed", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	a.logger.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("server shutdown failed", zap.Error(err))
		return err
	}

	a.logger.Info("server stopped gracefully")

	return nil
}

