package handler

import (
	"jsunnykhan/go-clean-template/internal/core/config"
	userUseCase "jsunnykhan/go-clean-template/internal/usecase/user"
	"jsunnykhan/go-clean-template/pkg/response"
	"time"

	"github.com/gin-gonic/gin"

	"jsunnykhan/go-clean-template/internal/adapter/http/middleware"
	apperrors "jsunnykhan/go-clean-template/pkg/error"
)

type registerRequest struct {
	Name     string `json:"name"     binding:"required,min=2,max=100"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// AuthResponse wraps a JWT token and user info returned after login/register.
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserHandler struct {
	uc     userUseCase.Repository
	jwtCfg config.JWTConfig
}

func NewUserHandler(uc userUseCase.Repository, jwtCfg config.JWTConfig) *UserHandler {
	return &UserHandler{
		uc:     uc,
		jwtCfg: jwtCfg,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	h.uc.Register(c.Request.Context(), userUseCase.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})

	c.JSON(201, gin.H{"message": "User registered successfully"})

}

func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Call the use case to authenticate the user
	user, err := h.uc.Login(c.Request.Context(), userUseCase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := middleware.GenerateJWTToken(user.ID.String(), string(user.Role), h.jwtCfg, time.Duration(h.jwtCfg.ExpirationTimeInHours)*time.Hour)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrInternal, "failed to generate token"))
		return
	}
	refreshToken, err := middleware.GenerateJWTToken(user.ID.String(), string(user.Role), h.jwtCfg, time.Duration(h.jwtCfg.RefreshExpirationTimeInHours)*time.Hour)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrInternal, "failed to generate refresh token"))
		return
	}

	response.OK(c, AuthResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
	})
}

func (h *UserHandler) Refresh(c *gin.Context) {

}
