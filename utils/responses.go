package utils

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CustomError : The error format in api response
type CustomError struct {
	Error string `json:"error"`
}

// Send : General function to send api response
func Send(c *gin.Context, status int, payload interface{}) {
	c.Header("Content-Type", "application/json")
	c.JSON(status, payload)
}

// Fail : General function to send api error response
func Fail(c *gin.Context, status int, details string) {
	response := &CustomError{
		Error: details,
	}
	result, err := json.Marshal(response)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "application/json")
	c.Status(status)
	c.Writer.Write(result)
}
