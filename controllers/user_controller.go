package controllers

import (
	"akash-project-go/database"
	"akash-project-go/middleware"
	"akash-project-go/models"
	"akash-project-go/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ResetPasswordRequest struct {
	Email       string `json:"email"`
	Otp         string `json:"otp"`
	PasswordOld string `json:"password_old"`
	PasswordNew string `json:"password_new"`
}

func UserRegistration(c *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.Error("UserRegistration@panicInfo:", panicInfo)
			utils.InternalServerErrorResponse(c, panicInfo.(error))
			return
		}
	}()

	var Reqbody models.UserDetail
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
		logrus.Error("Failed to Create user detail")
		utils.InternalServerErrorResponse(c, err)
		db.Rollback()
		return
	}

	db.Commit()
	utils.SuccessResponse(c, "User detail stored..", Reqbody)
}

func UserLoginRequest(c *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.Error("UserLoginRequest@panicInfo:", panicInfo)
			utils.InternalServerErrorResponse(c, panicInfo.(error))
			return
		}
	}()

	var require models.LoginRequest // request accept from user
	if err := c.ShouldBindJSON(&require); err != nil {
		logrus.Error("Failed to bind request body: ", err)
		utils.BadRequestResponse(c, err.Error(), err)
		return
	}

	var user models.UserDetail

	db := database.BookingDB.Begin() //Initialise Database
	defer db.Rollback()

	if err := db.Where("email = ?", require.Email).First(&user).Error; err != nil {
		logrus.Error("Invalide email address")
		utils.InternalServerErrorResponse(c, err)
		return
	}

	if err := user.Checkpassword(require.Password); err != nil {
		logrus.Error("Invalide Password")
		utils.InternalServerErrorResponse(c, err)
		return
	}

	//genrate Token BY User Email Address
	tokenString, err := middleware.GenrateJwtToken(user.Email)
	if err != nil {
		utils.UnauthorizedResponse(c, "Could not getrate Token", "employe")
		return
	}

	db.Commit()

	utils.SuccessResponse(c, "Login SucessFully Token Genrate", tokenString)
}

func ForgotPassword(c *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.Error("ForgotPassword@panicInfo:", panicInfo)
			utils.InternalServerErrorResponse(c, panicInfo.(error))
			return
		}
	}()

	// Email request
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Error("Failed to bind request body: ", err)
		utils.BadRequestResponse(c, err.Error(), err)
		return
	}

	var user models.UserDetail
	// Initialize database
	db := database.BookingDB.Begin()

	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		logrus.Error("Invalid email address: ", req.Email)
		utils.BadRequestResponse(c, "Invalid email address", err)
		db.Rollback()
		return
	}

	// Generate OTP
	otp := utils.GenerateOTP()
	// Store OTP in the database
	user.OTP = otp
	if err := db.Where("otp = ?", user.OTP).Save(&user).Error; err != nil {
		logrus.Error("Failed to save OTP database: ", err)
		utils.InternalServerErrorResponse(c, err)
		db.Rollback()
		return
	}

	// Send OTP via email
	if err := utils.SendOTPEmail(user); err != nil {
		logrus.Error("Failed to send OTP email: ", err)
		utils.InternalServerErrorResponse(c, err)
		return
	}

	db.Commit()
	utils.SuccessResponse(c, "OTP Send On Email", otp)
}

func UserOTPValidiation(c *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.Error("UserOTPValidiation@panicInfo:", panicInfo)
			utils.InternalServerErrorResponse(c, panicInfo.(error))
			return
		}
	}()

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Error("Failed to bind request body: ", err)
		utils.BadRequestResponse(c, err.Error(), err)
		return
	}

	var user models.UserDetail
	db := database.BookingDB.Begin() //initialise database

	if err := db.Where("otp = ?", req.Otp).First(&user).Error; err != nil {
		logrus.Error("OTP is not Valide: ", err)
		utils.InternalServerErrorResponse(c, err)
		db.Rollback()
		return
	}

	tokenString, err := middleware.GenrateJwtToken(user.Email) // Implement this function to generate a secure token
	if err != nil {
		logrus.Error("Failed to generate reset token: ", err)
		utils.InternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "OTP is Valid And Genrate token", tokenString)

}

func ResetPasswordUser(c *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.Error("ResetPasswordUser@panicInfo:", panicInfo)
			utils.InternalServerErrorResponse(c, panicInfo.(error))
		}
	}()

	var request ResetPasswordRequest
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
	

	UserEmail := userEmail.(string)

	var user models.UserDetail
	db := database.BookingDB.Begin()
	defer func() {
		if r := recover(); r != nil || db.Error != nil {
			db.Rollback()
		}
	}()

	// Find the user by email
	if err := db.Where("email = ?", UserEmail).First(&user).Error; err != nil {
		logrus.Error("User not found: ", err)
		utils.BadRequestResponse(c, "User not found", nil)
		return
	}

	if err := user.Hashpassword(request.PasswordNew); err != nil {
		logrus.Error("Failed to hash password: ", err)
		utils.InternalServerErrorResponse(c, err)
		return
	}

	// Update user's password
	if err := db.Save(&user).Error; err != nil {
		logrus.Error("Failed to update password: ", err)
		utils.InternalServerErrorResponse(c, err)
		return
	}

	if err := db.Commit().Error; err != nil {
		logrus.Error("Failed to commit transaction: ", err)
		utils.InternalServerErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "Password reset successfully", user)
}

func ChangePasswordUser(c *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.Error("ChangePasswordUser@panicInfo:", panicInfo)
			utils.InternalServerErrorResponse(c, panicInfo.(error))
		}
	}()

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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
	UserEmail := userEmail.(string)

	var user models.UserDetail
	db := database.BookingDB.Begin()

	// Find the user by email
	if err := db.Where("email = ?", UserEmail).First(&user).Error; err != nil {
		logrus.Error("User not found: ", err)
		utils.BadRequestResponse(c, "User not found", nil)
		return
	}
	//ckecck password
	if err := user.Checkpassword(req.PasswordOld); err != nil {
		logrus.Error("Invalide Password")
		utils.InternalServerErrorResponse(c, err)
		return
	}

	if err := user.Hashpassword(req.PasswordNew); err != nil {
		logrus.Error("Failed to hash password: ", err)
		utils.InternalServerErrorResponse(c, err)
		return
	}

	// Update user's password
	if err := db.Save(&user).Error; err != nil {
		logrus.Error("Failed to Change password: ", err)
		utils.InternalServerErrorResponse(c, err)
		return
	}

	db.Commit()

	utils.SuccessResponse(c, "user change password", nil)
}
