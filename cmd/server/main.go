package main

import (
	"log"
	"net"
	"os"

	pb "github.com/yourusername/sso-grpc/api"
	"github.com/yourusername/sso-grpc/internal/database"
	"github.com/yourusername/sso-grpc/internal/service"
	"google.golang.org/grpc"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	authService := service.NewAuthService(db)

	// Запуск gRPC-сервера
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка открытия порта: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authService)

	log.Println("gRPC-сервер запущен на порту 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка работы сервера: %v", err)
	}
}
