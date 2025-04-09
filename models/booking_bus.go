package models

import (
	"time"
)

type BookingBus struct {
	BookingId     uint      `gorm:"column:booking_id;primaryKey;autoIncrement" json:"booking_id"`
	PassengerName string    `gorm:"column:passenger_name" json:"passenger_name" binding:"required"`
	Gender        string    `gorm:"column:gender" json:"gender" binding:"required"`
	BusType       string    `gorm:"column:bus_type" json:"bus_type" binding:"required"`
	BookingDate   string    `gorm:"column:booking_date" json:"booking_date" binding:"required"`
	Location      string    `gorm:"column:location" json:"location" binding:"required"`
	Destination   string    `gorm:"column:destination" json:"destination" binding:"required"`
	UserEmail     string    `gorm:"foreignKey:user_email;references:email" json:"user_email"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (b *BookingBus) TableName() string {
	return "booking_bus"
}

type PasengerBookingDetail struct {
	BookingId     uint   `json:"booking_id"`
	UserEmail     string `json:"user_email"`
	MobileNo      string `json:"mobile_no"`
	PassengerName string `json:"passenger_name" binding:"required"`
	Gender        string `json:"gender" binding:"required"`
	BusType       string `json:"bus_type" binding:"required"`
	BookingDate   string `json:"booking_date" binding:"required"`
	Location      string `json:"location" binding:"required"`
	Destination   string `json:"destination" binding:"required"`
}