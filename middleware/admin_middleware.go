package middleware

// import (
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/gin-gonic/gin"
// )

// var JWTsecratekey = []byte("shree_booking_admin@123") // Same secret key used for token generation
// var TokenBlackList = make(map[string]bool)

// func AuthMiddLewareUser() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tokenString := c.Request.Header.Get("Authorization")

// 		if tokenString == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"message": "No authorization token provided"})
// 			c.Abort()
// 			return
// 		}

// 		// Remove "Bearer " from the token string
// 		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

// 		// Check if the token is blacklisted ---this If work at logout session time-----
// 		if TokenBlacklist[tokenString] {
// 			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token has been invalidated"})
// 			c.Abort()
// 			return
// 		} //------end-------

// 		claims := jwt.MapClaims{}
// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 			return JWTsecrateKey, nil
// 		})

// 		if err != nil || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
// 			c.Abort()
// 			return
// 		}

// 		// Optionally set user info in the context
// 		c.Set("email", claims["email"])
// 		c.Next()
// 	}
// }

// func GenrateJWTToken(email string) (string, error) {
// 	claims := jwt.MapClaims{
// 		"email": email,
// 		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString(JWTsecrateKey)
// 	if err != nil {
// 		return "", err
// 	}

// 	return tokenString, nil
// }
