package main

import (
	"gRPCAuthService/internal/app"
	"gRPCAuthService/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting server",
		slog.String("env", cfg.Env),
	)

	application := app.NewApp(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCServer.MustRun()

	//TODO: init gRPC app server

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	application.GRPCServer.Stop()
	log.Info("stopping server")
}

func setupLogger(env string) *slog.Logger {
	switch env {
	case envLocal:
		log := slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		return log

	case envDev:
		log := slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		return log

	case envProd:
		log := slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		return log
	}
	return nil
}
