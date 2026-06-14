package utils

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
}

func OK(c *gin.Context, message string, data interface{}) {
	c.JSON(200, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(201, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func BadRequest(c *gin.Context, message string, error interface{}) {
	c.JSON(400, Response{
		Success: false,
		Message: message,
		Error:   error,
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(401, Response{
		Success: false,
		Message: message,
	})
}

func InternalError(c *gin.Context, message string, error interface{}) {
	c.JSON(500, Response{
		Success: false,
		Message: message,
		Error:   error,
	})
}

func OkPaginated(c *gin.Context, message string, data interface{}, pagination Pagination) {
	c.JSON(200, PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}
