package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestGenerateToken(t *testing.T) {
	secret := "111222333"
	expiration := time.Hour
	service := NewJWTService(secret, expiration)

	userId := uuid.New()
	token, err := service.GenerateToken(userId, "Test", "Test@gmail.com")
	if err != nil {
		t.Error("expected no error, have:", err)
	}

	if token == "" {
		t.Error("expected token, have empty string")
	}
}

func TestParseToken(t *testing.T) {
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

	claims := &Claims{}
	TokenOut, err := service.ParseToken(tokenIn, claims)
	if err != nil {
		t.Fatalf("expected no error, have: %v", err)
	}

	if !TokenOut.Valid {
		t.Error("expected valid token, got invalid")
	}

	if *claims.UserId != userId {
		t.Errorf("User id expected %s, got %s", userId, *claims.UserId)
	}

	if claims.UserName != userName {
		t.Errorf("User name expected %v, got %v", userName, claims.UserName)
	}

	if claims.UserEmail != userEmail {
		t.Errorf("User email expected %v, got %v", userEmail, claims.UserEmail)
	}
}

func TestTokenWithInvalidSignature(t *testing.T) {
	service1 := NewJWTService("111111", time.Hour)
	token1, err := service1.GenerateToken(uuid.New(), "Test3", "Test3@gmail.com")
	if err != nil {
		t.Fatalf("expected no error, have: %v", err)
	}

	service2 := NewJWTService("222222", time.Hour)
	claims := &Claims{}
	_, err = service2.ParseToken(token1, claims)

	if err == nil {
		t.Error("expected error about invalid signature, got nil")
	}
}

func TestTokenWithWrongExpTime(t *testing.T) {
	service := NewJWTService("123456", -time.Hour)
	token, err := service.GenerateToken(uuid.New(), "Test4", "Test4@gmail.com")
	if err != nil {
		t.Fatalf("expected no error, have: %v", err)
	}

	claims := &Claims{}
	_, err = service.ParseToken(token, claims)

	if err == nil {
		t.Error("expected error about expired token, got nil")
	}
}

func TestTokenWithWrongMethod(t *testing.T) {
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token1, err := token.SignedString([]byte("654321"))
	if err != nil {
		t.Fatalf("error in create token with method HS512: %v", err)
	}

	parsedClaims := &Claims{}
	_, err = service.ParseToken(token1, parsedClaims)

	if err == nil {
		t.Error("expected error /invalid method/, got nil")
	}
}
