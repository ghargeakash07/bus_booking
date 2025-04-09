package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AdminRequest struct {
	AdminId   uint      `gorm:"column:admin_id;primaryKey;autoIncrement" json:"admin_id"`
	Name      string    `gorm:"column:name" json:"name" binding:"required"`
	Email     string    `gorm:"column:email" json:"email" binding:"required"`
	Password  string    `gorm:"column:password" json:"password" binding:"required"`
	MobileNo  string    `gorm:"column:mobile_no" json:"mobile_no" binding:"required"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (e *AdminRequest) Hashpassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	e.Password = string(hashedPassword)
	return nil
}

func (e *AdminRequest) Checkpassword(password string) error {

	return bcrypt.CompareHashAndPassword([]byte(e.Password), []byte(password))
}

func (m *AdminRequest) TableName() string {
	return "admin_request"
}
