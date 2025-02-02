package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"user_segmentation_service/internal/config"
	"user_segmentation_service/internal/db"
	"user_segmentation_service/internal/logger"
	"user_segmentation_service/internal/modules/segment_service"
	"user_segmentation_service/internal/modules/user_service"
	"user_segmentation_service/internal/server"
)

var (
	storage *pgxpool.Pool
	serv    *server.APIServer
)

// main - entry point.
func main() {
	cfg := config.MustLoad()
	logg := logger.Init(cfg.Log.Level, cfg.Log.AddSource)
	logg.Info("Application Started!")

	ctx, ctxCancel := context.WithCancel(context.Background())

	storage, err := db.NewPostgresPool(ctx, cfg.DB)
	if err != nil {
		logg.Error("db.NewPostgresPool", "err", err)
		os.Exit(1)
	}
	uu := user_service.NewUserService(storage)
	ss := segment_service.NewSegmentService(storage)
	serv := server.New(ctx, cfg.APIServer, uu, ss)

	go func() {
		if err := serv.Start(); err != nil {
			logg.Error("serv.Start", "err", err)
		}
	}()

	gracefulShutdown(ctxCancel)
}

// gracefulShutdown listens for interrupt signals (e.g., SIGTERM, os.Interrupt)
// to initiate a graceful shutdown.
func gracefulShutdown(ctxCancel context.CancelFunc) {
	// Channel for processing the completion signal.
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	slog.Debug("Shutting down gracefully...")
	ctxCancel()

	if storage != nil {
		storage.Close()
	}

	time.Sleep(2 * time.Second)
	slog.Info("Application Stopped!")
	close(signalChan)
}
