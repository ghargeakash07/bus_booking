package routes

import (
	//"akash-project-go/controllers"
	"akash-project-go/controllers"
	"akash-project-go/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// Create a Gin router
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "OPTIONS", "PUT", "PATCH"},
		//AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token", "x-api-key"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  false,
		AllowOriginFunc:  func(origin string) bool { return true },
		MaxAge:           86400,
	}))

	router.Use(middleware.LoggingMiddleware())

	travelerbooking := router.Group("user")
	{
		travelerbooking.POST("/registration", controllers.UserRegistration)
		travelerbooking.POST("/login", controllers.UserLoginRequest)
		travelerbooking.POST("/forgot-passworrd", controllers.ForgotPassword)
		travelerbooking.POST("/validate-otp", controllers.UserOTPValidiation)
		travelerbooking.Use(middleware.AuthMiddleware())
		{
			travelerbooking.POST("/booking", controllers.UserBookingBus)
			travelerbooking.GET("/booking-list", controllers.BookingList)
			travelerbooking.POST("/booling-delete/:id", controllers.DeleteBooking)
			travelerbooking.POST("/booking-update/:id", controllers.BookingUpdate)
			travelerbooking.POST("/reset-password", controllers.ResetPasswordUser)
			travelerbooking.POST("/change-password", controllers.ChangePasswordUser)
		}

	}

	adminrotes := router.Group("admin")
	{
		adminrotes.POST("/registration", controllers.RegistrationAdmin)
		adminrotes.POST("/admin-login", controllers.AdminLoginRequest)
		adminrotes.Use(middleware.AuthMiddleware())
		{
			adminrotes.GET("/pasenger-list", controllers.BookingAllPassengerList)
		}
	}

	return router
}
