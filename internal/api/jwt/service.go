package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"regexp"
	"time"
)

// сам jwt-сервис
type Service struct {
	secret     []byte
	expiration time.Duration
}

// конструктор jwt-сервиса
func NewJWTService(secret string, expiration time.Duration) TokenService {
	return &Service{
		secret:     []byte(secret),
		expiration: expiration,
	}
}

// валидация кастомных полей клеймов
func (c *Claims) CustomFieldsValidate() error {
	// проверяем валидность uuid
	if c.UserId == nil {
		return fmt.Errorf("User id is empty")
	}

	// проверяем, что имя пользователя не пустое
	if c.UserName == "" {
		return fmt.Errorf("Invalid user name")
	}

	// проверяем валидность почты
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, c.UserEmail)
	if err != nil {
		return fmt.Errorf("Error in matchString: %w", err)
	}
	if !matched {
		return fmt.Errorf("User email not looks like email")
	}

	return nil
}

// создание токена
func (j *Service) GenerateToken(user_id uuid.UUID, user_name string, user_email string) (string, error) {
	claims := &Claims{
		UserId:    &user_id,
		UserName:  user_name,
		UserEmail: user_email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiration)),
		},
	}

	// валидация кастомных полей
	if err := claims.CustomFieldsValidate(); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// декодировка токена
func (j *Service) ParseToken(tokenString string, claims *Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid method")
		}
		return j.secret, nil
	})
}
