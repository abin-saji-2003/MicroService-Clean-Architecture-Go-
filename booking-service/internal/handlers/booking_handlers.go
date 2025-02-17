package handlers

import (
	"context"
	"fmt"

	"booking-service/api/proto"
	"booking-service/internal/models"
	bookingRepo "booking-service/internal/repository"
	userProto "user-service/api/proto"

	"gorm.io/gorm"
)

type BookingHandler struct {
	bookingRepo bookingRepo.BookingRepository
	userRepo    userProto.UserRepository
	proto.UnimplementedBookingServiceServer
}

func NewBookingHandler(bookingRepo bookingRepo.BookingRepository, userRepo userProto.UserRepository) *BookingHandler {
	return &BookingHandler{
		bookingRepo: bookingRepo,
		userRepo:    userRepo,
	}
}

// ✅ Create Booking
func (h *BookingHandler) CreateBooking(ctx context.Context, req *proto.CreateBookingRequest) (*proto.CreateBookingResponse, error) {
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

	return &proto.CreateBookingResponse{Message: "Booking created successfully"}, nil
}

// ✅ Get Booking (Includes User Details)
func (h *BookingHandler) GetBooking(ctx context.Context, req *proto.GetBookingRequest) (*proto.GetBookingResponse, error) {
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

	// Fetch user details
	user, err := h.userRepo.GetUserByID(booking.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user details: %v", err)
	}

	// ✅ Return booking + user details
	return &proto.GetBookingResponse{
		BookingId:  uint32(booking.ID),
		UserId:     uint32(user.ID),
		UserName:   user.Name,  // ✅ Added User Name
		UserEmail:  user.Email, // ✅ Added User Email
		TotalPrice: booking.TotalPrice,
		Status:     booking.Status,
	}, nil
}

// ✅ Cancel Booking
func (h *BookingHandler) CancelBooking(ctx context.Context, req *proto.CancelBookingRequest) (*proto.CancelBookingResponse, error) {
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

	return &proto.CancelBookingResponse{Message: "Booking canceled successfully"}, nil
}
