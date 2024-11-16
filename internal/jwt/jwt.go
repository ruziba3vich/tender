package jwttokens

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/zohirovs/internal/models"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func GenerateAccessToken(jwtKey string, tokenClaims *models.TokenClaims) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:   tokenClaims.UserID,
		Username: tokenClaims.Username,
		Role:     tokenClaims.Role,
		Email:    tokenClaims.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Convert jwtKey to []byte
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
