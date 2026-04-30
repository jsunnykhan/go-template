package response

import (
	"errors"
	"net/http"

	apperrors "jsunnykhan/go-clean-template/pkg/error"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    any         `json:"data,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
	Meta    *Pagination `json:"meta,omitempty"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Pagination struct {
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// ── Success helpers ───────────────────────────────────────────────────────────

// OK sends 200 with data.
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{Success: true, Data: data})
}

// Created sends 201 with data.
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{Success: true, Data: data})
}

// List sends 200 with data + pagination meta.
func List(c *gin.Context, data any, total, offset, limit int) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta:    &Pagination{Total: total, Offset: offset, Limit: limit},
	})
}

// NoContent sends 204.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// ── Error helpers ─────────────────────────────────────────────────────────────

func Error(c *gin.Context, err error) {
	var appErr *apperrors.AppError
	if !errors.As(err, &appErr) {
		// Unexpected error — don't leak internals.
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   &ErrorBody{Code: "INTERNAL_ERROR", Message: "an unexpected error occurred"},
		})
		return
	}

	status := statusFor(appErr.Err)
	c.JSON(status, Response{
		Success: false,
		Error:   &ErrorBody{Code: appErr.Code, Message: appErr.Message},
	})
}

func statusFor(sentinel error) int {
	switch {
	case errors.Is(sentinel, apperrors.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(sentinel, apperrors.ErrConflict):
		return http.StatusConflict
	case errors.Is(sentinel, apperrors.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(sentinel, apperrors.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(sentinel, apperrors.ErrBadRequest):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
