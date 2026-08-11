package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserId    *uuid.UUID `json:"user_id"`
	UserName  string     `json:"user_name"`
	UserEmail string     `json:"user_email"`
	jwt.RegisteredClaims
}
