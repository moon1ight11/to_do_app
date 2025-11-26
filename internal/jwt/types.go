package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"regexp"
)

type Claims struct {
	UserId    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	UserEmail string    `json:"user_email"`
	jwt.RegisteredClaims
}

// валидация кастомных полей клеймов
func (c *Claims) CustomFieldsValidate() error {
	// проверяем валидность uuid
	if c.UserId == uuid.Nil {
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
