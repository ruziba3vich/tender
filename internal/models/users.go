package models

type Role string

var (
	Client     Role = "Client"
	Contractor Role = "Contractor"
)

type (
	User struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     Role   `json:"role"`
	}
)
