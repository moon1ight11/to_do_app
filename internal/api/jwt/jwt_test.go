package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"testing"
	"time"
)

// проверка правильной генерации токена
func TestGenerateToken(t *testing.T) {
	// создаем клеймы для токена и генерируем его
	secret := "111222333"
	expiration := time.Hour
	service := NewJWTService(secret, expiration)

	userId := uuid.New()
	token, err := service.GenerateToken(userId, "Test", "Test@gmail.com")
	if err != nil {
		t.Error("expected no error, have:", err)
	}

	// если сгенерировалась пустая строка
	if token == "" {
		t.Error("expected token, have empty string")
	}
}

// проверка правильного парсинга токена
func TestParseToken(t *testing.T) {
	// создаем клеймы для токена и генерируем его
	secret := "111222333"
	expiration := time.Hour
	service := NewJWTService(secret, expiration)

	userId := uuid.New()
	userName := "Test2"
	userEmail := "Test2@gmail.com"

	tokenIn, err := service.GenerateToken(userId, userName, userEmail)
	if err != nil {
		t.Fatalf("expected no error, have: %v", err)
	}

	// парсим токен
	claims := &Claims{}
	TokenOut, err := service.ParseToken(tokenIn, claims)
	if err != nil {
		t.Fatalf("expected no error, have: %v", err)
	}

	// ошибка валидации токена
	if !TokenOut.Valid {
		t.Error("expected valid token, got invalid")
	}

	// не совпал идентификатор
	if *claims.UserId != userId {
		 t.Errorf("User id expected %s, got %s", userId, *claims.UserId)
	}

	// не совпало имя
	if claims.UserName != userName {
		t.Errorf("User name expected %v, got %v", userName, claims.UserName)
	}

	// не совпала почта
	if claims.UserEmail != userEmail {
		t.Errorf("User email expected %v, got %v", userEmail, claims.UserEmail)
	}
}

// проверка парсинга токена с неправильной подписью
func TestTokenWithInvalidSignature(t *testing.T) {
	// создаем сервис и генерируем токен с шифром №1
	service1 := NewJWTService("111111", time.Hour)
	token1, err := service1.GenerateToken(uuid.New(), "Test3", "Test3@gmail.com")
	if err != nil {
		t.Fatalf("expected no error, have: %v", err)
	}

	// создаем сервис и парсим токен с шифром №2
	service2 := NewJWTService("222222", time.Hour)
	claims := &Claims{}
	_, err = service2.ParseToken(token1, claims)

	// если ошибок нет - плохо
	if err == nil {
		t.Error("expected error about invalid signature, got nil")
	}
}

// проверка парсинга токена с истекшим сроком годности
func TestTokenWithWrongExpTime(t *testing.T) {
	// создаем сервис, который генерирует просроченные токены
	service := NewJWTService("123456", -time.Hour)
	token, err := service.GenerateToken(uuid.New(), "Test4", "Test4@gmail.com")
	if err != nil {
		t.Fatalf("expected no error, have: %v", err)
	}

	// парсим просроченный токен
	claims := &Claims{}
	_, err = service.ParseToken(token, claims)

	// если ошибок нет - плохо
	if err == nil {
		t.Error("expected error about expired token, got nil")
	}
}

// проверка парсинга токена с неправаильным методом шифрования
func TestTokenWithWrongMethod(t *testing.T) {
	// создаем сервис, который использует HS256
	service := NewJWTService("654321", time.Hour)
	newUUID := uuid.New()

	claims := Claims{
		UserId:    &newUUID,
		UserName:  "Test5",
		UserEmail: "Test5@gmail.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	// создаем токен, который шифруется HS512
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token1, err := token.SignedString([]byte("654321"))
	if err != nil {
		t.Fatalf("error in create token with method HS512: %v", err)
	}

	// парсим токен, зашифрованный HS512, с помощью HS256
	parsedClaims := &Claims{}
	_, err = service.ParseToken(token1, parsedClaims)

	// если ошибок нет - плохо
	if err == nil {
		t.Error("expexted error /invalid method/, got nil")
	}
}
