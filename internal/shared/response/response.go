// internal/shared/response/response.go
package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// APIError represents an error response
type APIError struct {
	Error string `json:"error"`
}

// Success represents a success response with optional data
type Success struct {
	Data interface{} `json:"data,omitempty"`
}

// JSON sends a JSON response
func JSON(c echo.Context, status int, data interface{}) error {
	return c.JSON(status, data)
}

// OK sends a 200 OK response
func OK(c echo.Context, data interface{}) error {
	return JSON(c, http.StatusOK, data)
}

// Created sends a 201 Created response
func Created(c echo.Context, data interface{}) error {
	return JSON(c, http.StatusCreated, data)
}

// NoContent sends a 204 No Content response
func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// BadRequest sends a 400 Bad Request error
func BadRequest(c echo.Context, msg string) error {
	return JSON(c, http.StatusBadRequest, APIError{Error: msg})
}

// NotFound sends a 404 Not Found error
func NotFound(c echo.Context, msg string) error {
	return JSON(c, http.StatusNotFound, APIError{Error: msg})
}

// InternalError sends a 500 Internal Server Error
func InternalError(c echo.Context, msg string) error {
	return JSON(c, http.StatusInternalServerError, APIError{Error: msg})
}
