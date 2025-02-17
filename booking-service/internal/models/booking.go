package models

import "time"

type Booking struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null;index"`
	TotalPrice float64   `json:"total_price" gorm:"not null"`
	Status     string    `json:"status" gorm:"default:pending"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}
