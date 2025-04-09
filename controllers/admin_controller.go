package controllers

import (
	"akash-project-go/database"
	"akash-project-go/middleware"
	"akash-project-go/models"
	"akash-project-go/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func BookingAllPassengerList(c *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.Error("BookingAllPassengerList@panicInfo:", panicInfo)
			utils.InternalServerErrorResponse(c, panicInfo.(error))
			return
		}
	}()

	var request []models.PasengerBookingDetail
	if err := c.Bind(&request); err != nil {
		utils.BadRequestResponse(c, "invalid request payload", err)
		return
	}

	db := database.BookingDB.Begin()

	if err := db.Table("booking_bus").
		Select("booking_bus.booking_id,booking_bus.passenger_name,user_detail.mobile_no,booking_bus.user_email,booking_bus.gender," +
			"booking_bus.bus_type,booking_bus.bus_type,booking_bus.booking_date,booking_bus.location,booking_bus.destination").
		Joins("LEFT JOIN user_detail ON booking_bus.user_email=user_detail.email").
		Scan(&request).Error; err != nil {
		logrus.Error("Failed to Fetch Booking Detail")
		utils.InternalServerErrorResponse(c, err)
		db.Rollback()
		return

	}

	db.Commit()
	utils.SuccessResponse(c, "Booking Detail", request)
}

func RegistrationAdmin(c *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.Error("RegistrationAdmin@panicInfo:", panicInfo)
			utils.InternalServerErrorResponse(c, panicInfo.(error))
			return
		}
	}()

	var Reqbody models.AdminRequest
	if err := c.ShouldBindJSON(&Reqbody); err != nil {
		logrus.Error("Failed to bind request body: ", err)
		utils.BadRequestResponse(c, err.Error(), err)
		return
	}

	// Hash password
	if err := Reqbody.Hashpassword(Reqbody.Password); err != nil {
		utils.InternalServerErrorResponse(c, err)
		return
	}

	//initialise the database
	db := database.BookingDB.Begin()

	if err := db.Create(&Reqbody).Error; err != nil {
		logrus.Error("Failed to Create client detail")
		utils.InternalServerErrorResponse(c, err)
		db.Rollback()
		return
	}

	db.Commit()
	utils.SuccessResponse(c, "Admin detail stored..", Reqbody)
}

func AdminLoginRequest(c *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.Error("AdminLoginRequest@panicInfo:", panicInfo)
			utils.InternalServerErrorResponse(c, panicInfo.(error))
			return
		}
	}()

	var require models.LoginRequest // request accept from admin
	if err := c.ShouldBindJSON(&require); err != nil {
		logrus.Error("Failed to bind request body: ", err)
		utils.BadRequestResponse(c, err.Error(), err)
		return
	}

	var detail models.AdminRequest

	db := database.BookingDB.Begin() //Initialise Database
	defer db.Rollback()

	if err := db.Where("email = ?", require.Email).First(&detail).Error; err != nil {
		logrus.Error("Invalide email address")
		utils.InternalServerErrorResponse(c, err)
		return
	}

	if err := detail.Checkpassword(require.Password); err != nil {
		logrus.Error("Invalide Password")
		utils.InternalServerErrorResponse(c, err)
		return
	}

	//genrate Token BY User Email Address
	tokenString, err := middleware.GenrateJwtToken(detail.Email)
	if err != nil {
		utils.UnauthorizedResponse(c, "Could not getrate Token", "employe")
		return
	}

	db.Commit()

	utils.SuccessResponse(c, "Login SucessFully Token Genrate", tokenString)
}
