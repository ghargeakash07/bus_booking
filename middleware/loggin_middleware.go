package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		requestLogger := logrus.New()

		// Create a log file
		logFileName := "request.log"
		logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			// Configure Logrus to write to both the console and the log file
			requestLogger.SetOutput(io.MultiWriter(os.Stdout, logFile))
		} else {
			requestLogger.Warn("Failed to log to file. Using default stderr.")
		}

		defer logFile.Close()

		c.Next()

		duration := time.Since(startTime)

		requestLogger.WithFields(logrus.Fields{
			"method":      c.Request.Method,
			"url":         c.Request.URL.String(),
			"ip":          c.ClientIP(),
			"status":      c.Writer.Status(),
			"errors":      c.Errors.String(),
			"duration":    duration,
			"formData":    getRequestData(c),
			"queryParams": getQueryParams(c),
		}).Info("RequestLog::")
	}
}

func getRequestData(c *gin.Context) map[string]interface{} {
	requestData := make(map[string]interface{})

	if c.Request.Method == "POST" {
		contentType := c.Request.Header.Get("Content-Type")

		switch contentType {
		case "application/json":
			var jsonMap map[string]interface{}
			if err := json.NewDecoder(c.Request.Body).Decode(&jsonMap); err == nil {
				requestData = jsonMap
			}
		default:
			err := c.Request.ParseMultipartForm(32 << 20)
			if err == nil {
				for key, values := range c.Request.MultipartForm.Value {
					requestData[key] = values
				}
				for key, files := range c.Request.MultipartForm.File {
					fileInfo := make([]map[string]string, 0)
					for _, file := range files {
						fileInfo = append(fileInfo, map[string]string{
							"filename": file.Filename,
							"size":     fmt.Sprintf("%d", file.Size),
						})
					}
					requestData[key] = fileInfo
				}
			}
		}
	}

	return requestData
}

func getQueryParams(c *gin.Context) map[string]interface{} {
	queryParams := make(map[string]interface{})

	if c.Request.Method == "GET" {
		queryValues := c.Request.URL.Query()
		for key, values := range queryValues {
			queryParams[key] = values
		}
	}

	return queryParams
}

var JWTsecrateKey = []byte("shree_booking_@123") // Same secret key used for token generation
var TokenBlacklist = make(map[string]bool)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "No authorization token provided"})
			c.Abort()
			return
		}

		// Remove "Bearer " from the token string
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Check if the token is blacklisted ---this If work at logout session time----- 
		if TokenBlacklist[tokenString] {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token has been invalidated"})
			c.Abort()
			return
		}//------end-------

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JWTsecrateKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		// Optionally set user info in the context
		c.Set("email", claims["email"])
		c.Next()
	}
}

func GenrateJwtToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTsecrateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
