package models

type AuthDto struct {
	Email    string `validate:"required,lte=255" json:"email"`
	Password string `validate:"required,lte=255" json:"password"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}
