package controllers

import (
	"akash-project-go/database"
	"akash-project-go/models"
	"akash-project-go/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func UserBookingBus(c *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.Error("UserBookingBus@panicInfo:", panicInfo)
			utils.InternalServerErrorResponse(c, panicInfo.(error))
		}
	}()

	var request models.BookingBus
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Error("Failed to bind request body: ", err)
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	userEmail, exists := c.Get("email")
	if !exists {
		logrus.Error("User email not found in context")
		utils.BadRequestResponse(c, "User not authenticated", nil)
		return
	}

	request.UserEmail = userEmail.(string)

	db := database.BookingDB.Begin()
	defer func() {
		if r := recover(); r != nil || db.Error != nil {
			db.Rollback()
		}
	}()

	// Create the booking
	if err := db.Table("booking_bus").Create(&request).Error; err != nil {
		logrus.Error("Booking Failed: ", err)
		utils.InternalServerErrorResponse(c, err)
		return
	}

	// Get admin email from the database
	var adminEmail models.AdminRequest
	if err := db.Table("admin_request").Select("email").First(&adminEmail).Error; err != nil {
		logrus.Error("Failed to retrieve admin email: ", err)
		utils.InternalServerErrorResponse(c, err)
		return
	}

	//logrus.Debug("Retrieved Admin Email: ", adminEmail.Email) // Add this line

	// Trigger email notifications
	userDetail := models.UserDetail{Email: request.UserEmail}
	if err := utils.BookingBusMailTrigger(userDetail, request); err != nil {
		logrus.Error("Failed to send user booking confirmation email: ", err)
	}

	adminDetail := models.UserDetail{Email: adminEmail.Email}
	if err := utils.BookingBusMailTrigger(adminDetail, request); err != nil {
		logrus.Error("Failed to send admin booking notification email: ", err)
	}

	if err := db.Commit().Error; err != nil {
		logrus.Error("Failed to commit transaction: ", err)
		utils.InternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "Ticket Booking Successfully", request)
}

func BookingList(c *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.Error("UserBookingBus@panicInfo:", panicInfo)
			utils.InternalServerErrorResponse(c, panicInfo.(error))
			return
		}
	}()

	// Extract email from context set by middleware
	userEmail, exists := c.Get("email")
	if !exists {
		utils.BadRequestResponse(c, "Email not found in request", nil)
		return
	}
	//logurs("",userEmail)
	var bookings []models.BookingBus

	// Use the email to filter booking details
	db := database.BookingDB.Begin()
	if err := db.Table("booking_bus").Where("user_email = ?", userEmail).Find(&bookings).Error; err != nil {
		logrus.Error("Failed to Fetch Booking Detail for Email: ", userEmail)
		utils.InternalServerErrorResponse(c, err)
		db.Rollback()
		return
	}

	db.Commit()
	utils.SuccessResponse(c, "Booking Detail", bookings)
}

func DeleteBooking(c *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.Error("UserBookingBus@panicInfo:", panicInfo)
			utils.InternalServerErrorResponse(c, panicInfo.(error))
			return
		}
	}()
	var DelbyId = c.Param("id")
	if DelbyId == "" {
		logrus.Error("Booking ID is required")
		utils.ValidationResponse(c, "Booking Id is required")
		return
	}

	var request models.BookingBus
	db := database.BookingDB.Begin() // Initialise the database

	if err := db.Table("booking_bus").Where("booking_id=?", DelbyId).First(&request).Error; err != nil {
		logrus.Error("Failed to Find Booking detail")
		utils.InternalServerErrorResponse(c, err)
		db.Rollback()
		return
	}

	db.Delete(request)
	db.Commit()

	utils.SuccessResponse(c, "Booking Delete Sucessfully", nil)
}

func BookingUpdate(c *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.Error("BookingUpdate@panicInfo:", panicInfo)
			utils.InternalServerErrorResponse(c, panicInfo.(error))
			return
		}
	}()

	// Retrieve booking ID from URL parameters
	id := c.Param("id")
	if id == "" {
		logrus.Error("Booking ID is required")
		utils.ValidationResponse(c, "Booking ID is required")
		return
	}

	var booking models.BookingBus

	// Bind JSON request to the booking variable
	if err := c.ShouldBindJSON(&booking); err != nil {
		logrus.Error("Failed to bind request body: ", err)
		utils.BadRequestResponse(c, "Invalid input data", nil)
		return
	}

	// Start a database transaction
	db := database.BookingDB.Begin()
	defer db.Rollback() // Rollback if we do not commit

	// Fetch existing booking
	var existingBooking models.BookingBus
	if err := db.Where("booking_id = ?", id).First(&existingBooking).Error; err != nil {
		logrus.WithField("error", err).Error("Booking not found")
		utils.NotFoundResponse(c, "Booking not found")
		return
	}

	// Prepare updates
	updates := models.BookingBus{
		PassengerName: booking.PassengerName,
		Gender:        booking.Gender,
		BusType:       booking.BusType,
		BookingDate:   booking.BookingDate,
		Location:      booking.Location,
		Destination:   booking.Destination,
	}

	// Update booking details
	if err := db.Model(&existingBooking).Updates(updates).Error; err != nil {
		logrus.Error("Failed to update booking details: ", err)
		utils.InternalServerErrorResponse(c, err)
		return
	}
	db.Commit()
	utils.SuccessResponse(c, "Booking detail updated", existingBooking)
}
