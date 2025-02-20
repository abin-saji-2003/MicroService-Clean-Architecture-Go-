package main

import (
	"fmt"
	"log"
	"net/http"

	"api-gateway/internal/handlers"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	bookingProto "github.com/abin-saji-2003/GRPC-Pkg/proto/bookingpb"
	userProto "github.com/abin-saji-2003/GRPC-Pkg/proto/userpb"
)

func main() {
	userConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to User Service: %v", err)
	}
	defer userConn.Close()
	userClient := userProto.NewUserServiceClient(userConn)

	bookingConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Booking Service: %v", err)
	}
	defer bookingConn.Close()
	bookingClient := bookingProto.NewBookingServiceClient(bookingConn)

	r := gin.Default()

	handlers.RegisterUserRoutes(r, userClient)
	handlers.RegisterBookingRoutes(r, bookingClient)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API Gateway is Running!"})
	})

	fmt.Println("API Gateway running on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}
