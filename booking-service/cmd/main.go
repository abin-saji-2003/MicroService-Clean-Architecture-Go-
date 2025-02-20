package main

import (
	"log"
	"net"

	"booking-service/db"
	"booking-service/internal/handlers"
	bookingRepo "booking-service/internal/repository"
	bookingProto "github.com/abin-saji-2003/GRPC-Pkg/proto/bookingpb"
	userProto "github.com/abin-saji-2003/GRPC-Pkg/proto/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	db.InitDB()

	userConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to User Service: %v", err)
	}
	defer userConn.Close()

	userClient := userProto.NewUserServiceClient(userConn)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen on port 50052: %v", err)
	}

	grpcServer := grpc.NewServer()
	bookingRepository := bookingRepo.NewBookingRepo(db.DB)

	bookingHandler := handlers.NewBookingHandler(bookingRepository, userClient)

	bookingProto.RegisterBookingServiceServer(grpcServer, bookingHandler)

	log.Println("Booking Service is running on port 50052...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
