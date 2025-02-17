package repository

import (
	"github.com/abin-saji-2003/MicroService-Clean-Architecture-Go-/booking-service/internal/models"
	"gorm.io/gorm"
)

type BookingRepository interface {
	CreateBooking(booking *models.Booking) error
	GetBookingByID(id uint) (*models.Booking, error)
	CancelBooking(id uint) error
}

type BookingRepo struct {
	DB *gorm.DB
}

func NewBookingRepo(db *gorm.DB) *BookingRepo {
	return &BookingRepo{DB: db}
}

func (r *BookingRepo) CreateBooking(booking *models.Booking) error {
	return r.DB.Create(booking).Error
}

func (r *BookingRepo) GetBookingByID(id uint) (*models.Booking, error) {
	var booking models.Booking

	if err := r.DB.First(&booking, id).Error; err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *BookingRepo) CancelBooking(id uint) error {
	return r.DB.Model(&models.Booking{}).Where("id = ?", id).Update("status", "canceled").Error
}
