package grpcapp

import (
	"fmt"
	authgrpc "gRPCAuthService/internal/gRPC/auth"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewApp(log *slog.Logger, auth authgrpc.Auth, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, auth)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil { // Правильный вызов метода
		a.log.Error("Failed to run server", "error", err)
		os.Exit(1) // Корректный выход с ненулевым статусом
	}
}

func (a *App) Run() error {
	const op = "gRPC.Start"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("starting gRPC server", slog.String("address", l.Addr().String()))

	if err = a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "gRPC.Stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC server", slog.Int("port", a.port))
	a.gRPCServer.GracefulStop()
}
