package dto

type LoginRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=3,max=64"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}
