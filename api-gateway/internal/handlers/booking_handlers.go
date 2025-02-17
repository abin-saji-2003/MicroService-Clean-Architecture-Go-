package handlers

import (
	"context"
	"net/http"
	"strconv"

	bookingProto "github.com/abin-saji-2003/MicroService-Clean-Architecture-Go-/booking-service/api/proto"
	"github.com/gin-gonic/gin"
)

// ✅ Register Booking Routes
func RegisterBookingRoutes(r *gin.Engine, bookingClient bookingProto.BookingServiceClient) {
	bookingRoutes := r.Group("/api/bookings")
	{
		bookingRoutes.POST("/", func(c *gin.Context) {
			var req bookingProto.CreateBookingRequest
			if err := c.BindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			resp, err := bookingClient.CreateBooking(context.Background(), &req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": resp.Message})
		})

		bookingRoutes.GET("/:id", func(c *gin.Context) {
			bookingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
				return
			}

			resp, err := bookingClient.GetBooking(context.Background(), &bookingProto.GetBookingRequest{
				BookingId: uint32(bookingID),
			})
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
				return
			}
			c.JSON(http.StatusOK, resp)
		})

		bookingRoutes.DELETE("/:id", func(c *gin.Context) {
			bookingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
				return
			}

			resp, err := bookingClient.CancelBooking(context.Background(), &bookingProto.CancelBookingRequest{
				BookingId: uint32(bookingID),
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel booking"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": resp.Message})
		})
	}
}
