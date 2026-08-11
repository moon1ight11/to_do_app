package jwt

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service struct {
	secret     []byte
	expiration time.Duration
}

func NewJWTService(secret string, expiration time.Duration) TokenService {
	return &Service{
		secret:     []byte(secret),
		expiration: expiration,
	}
}

func (c *Claims) CustomFieldsValidate() error {
	if c.UserId == nil {
		return fmt.Errorf("jwt.CustomFieldsValidate: user id is empty")
	}

	if c.UserName == "" {
		return fmt.Errorf("jwt.CustomFieldsValidate: user name is empty")
	}

	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, c.UserEmail)
	if err != nil {
		return fmt.Errorf("jwt.CustomFieldsValidate: match email: %w", err)
	}
	if !matched {
		return fmt.Errorf("jwt.CustomFieldsValidate: email not valid")
	}

	return nil
}

func (j *Service) GenerateToken(userId uuid.UUID, name string, email string) (string, error) {
	claims := &Claims{
		UserId:    &userId,
		UserName:  name,
		UserEmail: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiration)),
		},
	}

	if err := claims.CustomFieldsValidate(); err != nil {
		return "", fmt.Errorf("jwt.GenerateToken: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *Service) ParseToken(tokenString string, claims *Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid signing method")
		}
		return j.secret, nil
	})
}
