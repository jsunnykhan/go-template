package userusecase

import (
	"context"
	"jsunnykhan/go-clean-template/internal/domain/user"
)

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type UpdateProfileInput struct {
	Name string
}

type Repository interface {
	Register(ctx context.Context, input RegisterInput) (*user.User, error)
	Login(ctx context.Context, input LoginInput) (*user.User, error)
}
