package handlers

import (
	"context"
	"fmt"

	"booking-service/internal/models"
	bookingRepo "booking-service/internal/repository"
	bookingProto "github.com/abin-saji-2003/GRPC-Pkg/proto/bookingpb"
	userProto "github.com/abin-saji-2003/GRPC-Pkg/proto/userpb"

	"gorm.io/gorm"
)

type BookingHandler struct {
	bookingRepo bookingRepo.BookingRepository
	userClient  userProto.UserServiceClient // ✅ Use gRPC Client Instead of Repository
	bookingProto.UnimplementedBookingServiceServer
}

// ✅ Updated Constructor to Accept gRPC User Client
func NewBookingHandler(bookingRepo bookingRepo.BookingRepository, userClient userProto.UserServiceClient) *BookingHandler {
	return &BookingHandler{
		bookingRepo: bookingRepo,
		userClient:  userClient, // ✅ Use gRPC client for user service calls
	}
}

// ✅ Create Booking
func (h *BookingHandler) CreateBooking(ctx context.Context, req *bookingProto.CreateBookingRequest) (*bookingProto.CreateBookingResponse, error) {
	// Validate input
	if req.UserId == 0 || req.TotalPrice <= 0 {
		return nil, fmt.Errorf("invalid input: user_id and total_price must be greater than zero")
	}

	booking := &models.Booking{
		UserID:     uint(req.UserId),
		TotalPrice: req.TotalPrice,
		Status:     "pending",
	}

	if err := h.bookingRepo.CreateBooking(booking); err != nil {
		return nil, fmt.Errorf("failed to create booking: %v", err)
	}

	return &bookingProto.CreateBookingResponse{Message: "Booking created successfully"}, nil
}

// ✅ Get Booking (Includes User Details)
func (h *BookingHandler) GetBooking(ctx context.Context, req *bookingProto.GetBookingRequest) (*bookingProto.GetBookingResponse, error) {
	if req.BookingId == 0 {
		return nil, fmt.Errorf("invalid booking ID")
	}

	// Fetch booking details
	booking, err := h.bookingRepo.GetBookingByID(uint(req.BookingId))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, fmt.Errorf("failed to retrieve booking: %v", err)
	}

	// Fetch user details via gRPC
	userResp, err := h.userClient.GetUserByID(ctx, &userProto.GetUserByIDRequest{
		UserId: uint32(booking.UserID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user details from User Service: %v", err)
	}

	// ✅ Return booking + user details
	return &bookingProto.GetBookingResponse{
		BookingId:  uint32(booking.ID),
		UserId:     userResp.UserId, // ✅ Retrieved via gRPC
		UserName:   userResp.Name,   // ✅ Retrieved via gRPC
		UserEmail:  userResp.Email,  // ✅ Retrieved via gRPC
		TotalPrice: booking.TotalPrice,
		Status:     booking.Status,
	}, nil
}

// ✅ Cancel Booking
func (h *BookingHandler) CancelBooking(ctx context.Context, req *bookingProto.CancelBookingRequest) (*bookingProto.CancelBookingResponse, error) {
	if req.BookingId == 0 {
		return nil, fmt.Errorf("invalid booking ID")
	}

	err := h.bookingRepo.CancelBooking(uint(req.BookingId))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, fmt.Errorf("failed to cancel booking: %v", err)
	}

	return &bookingProto.CancelBookingResponse{Message: "Booking canceled successfully"}, nil
}
