package main

import (
	"fmt"
	"log"
	"net/http"

	"api-gateway/internal/handlers"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	bookingProto "github.com/abin-saji-2003/MicroService-Clean-Architecture-Go-/booking-service/api/proto"
	userProto "github.com/abin-saji-2003/MicroService-Clean-Architecture-Go-/user-service/api/proto"
)

func main() {
	// ✅ Connect to User Service
	userConn, err := grpc.Dial("user-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Failed to connect to User Service: %v", err)
	}
	defer userConn.Close()
	userClient := userProto.NewUserServiceClient(userConn)

	// ✅ Connect to Booking Service
	bookingConn, err := grpc.Dial("booking-service:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Failed to connect to Booking Service: %v", err)
	}
	defer bookingConn.Close()
	bookingClient := bookingProto.NewBookingServiceClient(bookingConn)

	// ✅ Create Gin Router
	r := gin.Default()

	// ✅ Register Handlers
	handlers.RegisterUserRoutes(r, userClient)
	handlers.RegisterBookingRoutes(r, bookingClient)

	// ✅ Health Check Route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API Gateway is Running!"})
	})

	// ✅ Start API Gateway
	fmt.Println("🚀 API Gateway running on port 8080...")
	r.Run(":8080")
}
