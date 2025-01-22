package app

import (
	grpcapp "gRPCAuthService/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func NewApp(log *slog.Logger, gRPCPort int, storagePath string, tokenTTL time.Duration) *App {
	grpcApp := grpcapp.NewApp(log, gRPCPort)
	return &App{
		GRPCServer: grpcApp,
	}
}
