package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserDetail struct {
	Id        uint      `gorm:"column:id;AUTO_INCREMENT" json:"id"`
	Name      string    `gorm:"column:name" json:"name" binding:"required"`
	Email     string    `gorm:"column:email" json:"email" binding:"required"`
	Password  string    `gorm:"column:password" json:"password" binding:"required"`
	MobileNo  string    `gorm:"column:mobile_no" json:"mobile_no" binding:"required"`
	OTP       string    `gorm:"column:otp" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *UserDetail) TableName() string {
	return "user_detail"
}

func (e *UserDetail) Hashpassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	e.Password = string(hashedPassword)
	return nil
}

func (e *UserDetail) Checkpassword(password string) error {

	return bcrypt.CompareHashAndPassword([]byte(e.Password), []byte(password))
}
