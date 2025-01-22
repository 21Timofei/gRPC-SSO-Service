package app

import (
	grpcapp "gRPCAuthService/internal/app/grpc"
	"gRPCAuthService/internal/services/auth"
	"gRPCAuthService/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func NewApp(log *slog.Logger, gRPCPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.NewStorage(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.NewAuth(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.NewApp(log, authService, gRPCPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
