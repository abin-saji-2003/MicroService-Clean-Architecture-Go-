package main

import (
	"log"
	"net"

	"github.com/abin-saji-2003/MicroService-Clean-Architecture-Go-/tree/main/booking-service/api/proto"
	"github.com/abin-saji-2003/MicroService-Clean-Architecture-Go-/tree/main/booking-service/db"
	"github.com/abin-saji-2003/MicroService-Clean-Architecture-Go-/tree/main/booking-service/internal/handlers"
	bookingRepo "github.com/abin-saji-2003/MicroService-Clean-Architecture-Go-/tree/main/booking-service/internal/repository"

	userProto "github.com/abin-saji-2003/MicroService-Clean-Architecture-Go-/tree/main/user-service/api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// ✅ Initialize Database
	db.InitDB()

	// ✅ Connect to gRPC User Service
	userConn, err := grpc.Dial("user-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Failed to connect to User Service: %v", err)
	}
	defer userConn.Close()

	userClient := userProto.NewUserServiceClient(userConn) // ✅ Create gRPC User Client

	// ✅ Start gRPC Server for Booking Service
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("❌ Failed to listen on port 50052: %v", err)
	}

	grpcServer := grpc.NewServer()
	bookingRepository := bookingRepo.NewBookingRepo(db.DB)

	// ✅ Pass the gRPC User Client instead of `userRepo`
	bookingHandler := handlers.NewBookingHandler(bookingRepository, userClient)

	proto.RegisterBookingServiceServer(grpcServer, bookingHandler)

	log.Println("✅ Booking Service is running on port 50052...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("❌ Failed to start gRPC server: %v", err)
	}
}
