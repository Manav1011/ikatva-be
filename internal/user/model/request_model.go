package model

// LoginRequest is the body for POST /users/login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SignupRequest struct {
	Name     string `json:"name" binding:"omitempty"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
