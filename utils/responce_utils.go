package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseBody defines the common structure for JSON responses.
type ResponseBody struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	DevMessage error
	Body       interface{} `json:"body"`
}

// JSONResponse sends a JSON response with the given status code.
func JSONResponse(c *gin.Context, statusCode int, message string, body interface{}) {
	response := ResponseBody{
		StatusCode: statusCode,
		Message:    message,
		Body:       body,
	}
	c.JSON(statusCode, response)
}

// SuccessResponse sends an HTTP 2xx success response.
func SuccessResponse(c *gin.Context, message string, body interface{}) {
	response := ResponseBody{
		StatusCode: http.StatusOK,
		Message:    message,
		Body:       body,
	}
	c.JSON(http.StatusOK, response)
}

func ValidationResponse(c *gin.Context, message string) {
	response := ResponseBody{
		StatusCode: http.StatusUnprocessableEntity,
		Message:    message,
		Body:       map[string]interface{}{},
	}
	c.JSON(http.StatusUnprocessableEntity, response)
}

// CreatedResponse sends an HTTP 201 response.
func CreatedResponse(c *gin.Context, message string, body interface{}) {
	response := ResponseBody{
		StatusCode: http.StatusCreated,
		Message:    message,
		Body:       body,
	}
	c.JSON(http.StatusCreated, response)
}

// AcceptedResponse sends an HTTP 202 response.
func AcceptedResponse(c *gin.Context, message string, body interface{}) {
	response := ResponseBody{
		StatusCode: http.StatusAccepted,
		Message:    message,
		Body:       body,
	}
	c.JSON(http.StatusAccepted, response)
}

// NoContentResponse sends an HTTP 204 response.
func NoContentResponse(c *gin.Context, message string, body interface{}) {
	response := ResponseBody{
		StatusCode: http.StatusNoContent,
		Message:    message,
		Body:       body,
	}
	c.JSON(http.StatusNoContent, response)
}

// BadRequestResponse sends an HTTP 400 response.
func BadRequestResponse(c *gin.Context, message string, body interface{}) {
	response := ResponseBody{
		StatusCode: http.StatusBadRequest,
		Message:    message,
		Body:       body,
	}
	c.JSON(http.StatusBadRequest, response)
}

// UnauthorizedResponse sends an HTTP 401 response.
func UnauthorizedResponse(c *gin.Context, message string, body interface{}) {
	response := ResponseBody{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
		Body:       body,
	}
	c.JSON(http.StatusUnauthorized, response)
}

// ForbiddenResponse sends an HTTP 403 response.
func ForbiddenResponse(c *gin.Context, message string, body interface{}) {
	response := ResponseBody{
		StatusCode: http.StatusForbidden,
		Message:    message,
		Body:       body,
	}
	c.JSON(http.StatusForbidden, response)
}

// NotFoundResponse sends an HTTP 404 response.
func NotFoundResponse(c *gin.Context, message string) {
	response := ResponseBody{
		StatusCode: http.StatusNotFound,
		Message:    message,
		Body:       map[string]interface{}{},
	}
	c.JSON(http.StatusNotFound, response)
}

// MethodNotAllowedResponse sends an HTTP 405 response.
func MethodNotAllowedResponse(c *gin.Context, message string, body interface{}) {
	response := ResponseBody{
		StatusCode: http.StatusMethodNotAllowed,
		Message:    message,
		Body:       body,
	}
	c.JSON(http.StatusMethodNotAllowed, response)
}

// InternalServerErrorResponse sends an HTTP 500 response.
func InternalServerErrorResponse(c *gin.Context, Error error) {
	var body interface{}
	response := ResponseBody{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
		DevMessage: Error,
		Body:       body,
	}
	c.JSON(http.StatusInternalServerError, response)
}

// ServiceUnavailableResponse sends an HTTP 503 response.
func ServiceUnavailableResponse(c *gin.Context, message string, body interface{}) {
	response := ResponseBody{
		StatusCode: http.StatusServiceUnavailable,
		Message:    message,
		Body:       body,
	}
	c.JSON(http.StatusServiceUnavailable, response)
}

// GatewayTimeoutResponse sends an HTTP 504 response.
func GatewayTimeoutResponse(c *gin.Context, message string, body interface{}) {
	response := ResponseBody{
		StatusCode: http.StatusGatewayTimeout,
		Message:    message,
		Body:       body,
	}
	c.JSON(http.StatusGatewayTimeout, response)
}

func InternalServerErrorWithMessage(c *gin.Context, Message string) {
	errorResponse := ResponseBody{
		StatusCode: 500,
		Message:    Message,
		Body:       map[string]interface{}{},
	}
	c.JSON(http.StatusInternalServerError, errorResponse)
}
