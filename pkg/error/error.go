package error

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound     = errors.New("resource not found")
	ErrConflict     = errors.New("resource already exists")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrBadRequest   = errors.New("bad request")
	ErrInternal     = errors.New("internal error")
)

// ── AppError ─────────────────────────────────────────────────────────────────

type AppError struct {
	Code    string // e.g. "NOT_FOUND", "CONFLICT"
	Message string // safe for API response bodies
	Err     error  // wrapped sentinel — use errors.Is() to inspect
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Err }

func New(sentinel error, message string) *AppError {
	return &AppError{
		Code:    codeFor(sentinel),
		Message: message,
		Err:     sentinel,
	}
}

func Newf(sentinel error, format string, args ...any) *AppError {
	return New(sentinel, fmt.Sprintf(format, args...))
}

func Is(err, target error) bool { return errors.Is(err, target) }

// ── helpers ───────────────────────────────────────────────────────────────────

func codeFor(err error) string {
	switch {
	case errors.Is(err, ErrNotFound):
		return "NOT_FOUND"
	case errors.Is(err, ErrConflict):
		return "CONFLICT"
	case errors.Is(err, ErrUnauthorized):
		return "UNAUTHORIZED"
	case errors.Is(err, ErrForbidden):
		return "FORBIDDEN"
	case errors.Is(err, ErrBadRequest):
		return "BAD_REQUEST"
	default:
		return "INTERNAL_ERROR"
	}
}
