package main

import (
	"log"
	"net"

	"booking-service/api/proto"
	"booking-service/db"
	"booking-service/internal/handlers"
	bookingRepo "booking-service/internal/repository"

	userProto "github.com/abin-saji-2003/MicroService-Clean-Architecture-Go-/user-service/api/proto"

	"google.golang.org/grpc"
)

func main() {
	db.InitDB()

	userConn, err := grpc.Dial("user-service:50051", grpc.WithInsecure()) // Use the correct address
	if err != nil {
		log.Fatalf("❌ Failed to connect to User Service: %v", err)
	}
	defer userConn.Close()
	userClient := userProto.NewUserServiceClient(userConn) // ✅ Create gRPC User Client

	// ✅ Start gRPC Server for Booking Service
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("❌ Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	bookingRepository := bookingRepo.NewBookingRepo(db.DB)

	// ✅ Pass the gRPC User Client instead of `userRepo`
	bookingHandler := handlers.NewBookingHandler(bookingRepository, userClient)

	proto.RegisterBookingServiceServer(grpcServer, bookingHandler)

	log.Println("✅ Booking Service is running on port 50052...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
