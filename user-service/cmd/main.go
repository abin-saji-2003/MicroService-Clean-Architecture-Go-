package main

import (
	"user-sevice/api/proto"
	"user-sevice/db"
	"user-sevice/internal/handlers"
	"user-sevice/internal/repository"

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

	proto.RegisterUserServiceServer(grpcServer, userHandler)

	log.Println("User Service is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
