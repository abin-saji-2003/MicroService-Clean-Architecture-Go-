package main

import (
	"user-service/db"
	"user-service/internal/handlers"
	"user-service/internal/repository"

	userProto "github.com/abin-saji-2003/GRPC-Pkg/proto/userpb"

	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	db.InitDB()

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userRepo := repository.NewUserRepo(db.DB)
	userHandler := handlers.NewUserHandler(userRepo)

	userProto.RegisterUserServiceServer(grpcServer, userHandler)

	log.Println("User Service is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
