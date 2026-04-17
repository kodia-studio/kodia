// Package response provides standardized JSON response helpers for Kodia Framework.
// All API responses follow a consistent structure for easier client-side handling.
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the standard API response envelope.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta holds pagination and other metadata.
type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// OK sends a 200 OK response with data.
func OK(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// OKWithMeta sends a 200 OK response with data and pagination metadata.
func OKWithMeta(c *gin.Context, message string, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Created sends a 201 Created response.
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// NoContent sends a 204 No Content response.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// BadRequest sends a 400 Bad Request response with validation errors.
func BadRequest(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

// Unauthorized sends a 401 Unauthorized response.
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Message: message,
	})
}

// Forbidden sends a 403 Forbidden response.
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Forbidden"
	}
	c.JSON(http.StatusForbidden, Response{
		Success: false,
		Message: message,
	})
}

// NotFound sends a 404 Not Found response.
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Message: message,
	})
}

// Conflict sends a 409 Conflict response.
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, Response{
		Success: false,
		Message: message,
	})
}

// UnprocessableEntity sends a 422 response for business logic validation failures.
func UnprocessableEntity(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusUnprocessableEntity, Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

// TooManyRequests sends a 429 Too Many Requests response.
func TooManyRequests(c *gin.Context) {
	c.JSON(http.StatusTooManyRequests, Response{
		Success: false,
		Message: "Too many requests. Please slow down.",
	})
}

// InternalServerError sends a 500 Internal Server Error response.
func InternalServerError(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error"
	}
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Message: message,
	})
}

// NewMeta constructs a Meta struct for paginated responses.
func NewMeta(page, perPage int, total int64) *Meta {
	totalPages := 0
	if perPage > 0 && total > 0 {
		totalPages = int((total + int64(perPage) - 1) / int64(perPage))
	}
	return &Meta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}
}
