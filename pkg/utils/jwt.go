package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret = []byte("RAHASIA_DAPUR_SAAS_INSPACEMAP_2025")

type JWTPayload struct {
	UserID         uuid.UUID `json:"user_id"`
	Email          string    `json:"email"`
	OrganizationID uuid.UUID `json:"org_id"`
	Role           string    `json:"role"`
	Permissions    []string  `json:"perms"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uuid.UUID, email string, orgID uuid.UUID, roleName string, permissions []string) (string, error) {
	claims := JWTPayload{
		UserID:         userID,
		Email:          email,
		OrganizationID: orgID,
		Role:           roleName,
		Permissions:    permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (*JWTPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTPayload{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTPayload); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
