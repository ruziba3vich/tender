package models

import "time"

type Role string

var (
	Client     Role = "client"
	Contractor Role = "contractor"
)

type (
	User struct {
		ID        string    `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Password  string    `json:"password"`
		Role      Role      `json:"role"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		DeletedAt time.Time `json:"deleted_at"`
	}

	RegisterUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     Role   `json:"role"`
		Username string `json:"username"`
	}

	ResetPassword struct {
		Email       string `json:"email"`
		NewPassword string `json:"new_password"`
		Code        int    `json:"code"`
	}

	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	TokenClaims struct {
		UserID   string `json:"user_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}
)
