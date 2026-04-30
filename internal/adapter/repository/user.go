package repository

import (
	"context"
	"errors"
	"fmt"
	"jsunnykhan/go-clean-template/internal/domain/user"
	"time"

	apperrors "jsunnykhan/go-clean-template/pkg/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userModel struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"uniqueIndex;not null;unique"`
	PasswordHash string    `gorm:"not null"`
	Role         string    `gorm:"type:varchar(50);default:'user';not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (u userModel) tableName() string {
	return "users"
}

// --- Mapping functions ---

func toUserEntity(m *userModel) *user.User {
	return &user.User{
		ID:           m.ID,
		Name:         m.Name,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		Role:         user.Role(m.Role),
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func fromUserEntity(u *user.User) *userModel {
	return &userModel{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         string(u.Role),
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

// ── Repository ───────────────────────────────────────────────────────────────
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) AutoMigrateModel() any { return &userModel{} }

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	m := fromUserEntity(u)
	result := r.db.WithContext(ctx).Create(m)
	if result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			return apperrors.Newf(apperrors.ErrConflict, "email %q already registered", u.Email)
		}
		return fmt.Errorf("postgres: create user: %w", result.Error)
	}
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var m userModel
	result := r.db.WithContext(ctx).First(&m, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.New(apperrors.ErrNotFound, "user not found")
		}
		return nil, fmt.Errorf("postgres: find user by id: %w", result.Error)
	}
	return toUserEntity(&m), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var m userModel
	result := r.db.WithContext(ctx).First(&m, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.New(apperrors.ErrNotFound, "user not found")
		}
		return nil, fmt.Errorf("postgres: find user by email: %w", result.Error)
	}
	return toUserEntity(&m), nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	m := fromUserEntity(u)
	result := r.db.WithContext(ctx).Save(m)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return apperrors.New(apperrors.ErrNotFound, "user not found")
		}
		return fmt.Errorf("postgres: update user: %w", result.Error)
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&userModel{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("postgres: delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.New(apperrors.ErrNotFound, "user not found")
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*user.User, error) {
	var models []userModel
	result := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&models)

	if result.Error != nil {
		return nil, fmt.Errorf("postgres: list users: %w", result.Error)
	}

	users := make([]*user.User, len(models))
	for i, m := range models {
		m := m 
		users[i] = toUserEntity(&m)
	}
	return users, nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

// isDuplicateKeyError checks for PostgreSQL unique constraint violations.
func isDuplicateKeyError(err error) bool {
	return err != nil && (
	// PostgreSQL error code 23505 = unique_violation
	fmt.Sprintf("%v", err) != "" &&
		containsStr(err.Error(), "23505", "duplicate key", "unique"))
}

func containsStr(s string, subs ...string) bool {
	for _, sub := range subs {
		if len(s) >= len(sub) {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
		}
	}
	return false
}
