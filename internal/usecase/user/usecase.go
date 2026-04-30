package userusecase

import (
	"context"
	"fmt"
	"time"

	"jsunnykhan/go-clean-template/internal/domain/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// userUseCase is the concrete implementation of UseCase.
type userUseCase struct {
	repo user.Repository
}

// New creates a new instance of userUseCase with the given repository.
func New(repo user.Repository) Repository {
	return &userUseCase{repo: repo}
}

// Register implements [UseCase].
func (uc *userUseCase) Register(ctx context.Context, input RegisterInput) (*user.User, error) {
	if _, err := uc.repo.FindByEmail(ctx, input.Email); err == nil {
		return nil, fmt.Errorf("user with email %s already exists", input.Email)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	now := time.Now().UTC()
	newUser := &user.User{
		ID:           uuid.New(),
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hash),
		Role:         user.RoleUser,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	
	if err := uc.repo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil

}

// Login implements [UseCase].
func (uc *userUseCase) Login(ctx context.Context, input LoginInput) (*user.User, error) {
	user, err := uc.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	return user, nil

}
