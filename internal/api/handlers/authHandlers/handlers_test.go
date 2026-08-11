package authhandlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace/noop"

	"todoapp/internal/api/jwt"
	"todoapp/internal/api/models"
	"todoapp/internal/metrics"
)

type mockLogger struct{}

func (m *mockLogger) Info(msg string, fields ...any)  {}
func (m *mockLogger) Error(msg string, fields ...any) {}
func (m *mockLogger) Fatal(msg string, fields ...any) {}
func (m *mockLogger) Close() error                    { return nil }

func newMetrics() *metrics.Metrics {
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	return metrics.NewMetrics()
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestSignUp_Success(t *testing.T) {
	mockSvc := &MockUserService{
		AddUserFunc: func(ctx context.Context, user models.UserAuth) (uuid.UUID, error) {
			return uuid.New(), nil
		},
	}

	jwtSvc := jwt.NewJWTService("test-secret", 3600)
	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewAuthHandler(mockSvc, jwtSvc, l, tracer, m)

	router := setupTestRouter()
	router.POST("/sign-up", handler.SignUp)

	body := models.UserAuth{Name: "test", Email: "test@test.com", Pass: "password"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/sign-up", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestSignUp_BadRequest_EmptyFields(t *testing.T) {
	mockSvc := &MockUserService{}
	jwtSvc := jwt.NewJWTService("test-secret", 3600)
	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewAuthHandler(mockSvc, jwtSvc, l, tracer, m)

	router := setupTestRouter()
	router.POST("/sign-up", handler.SignUp)

	body := models.UserAuth{Name: "", Email: "test@test.com", Pass: "password"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/sign-up", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestSignUp_InternalError(t *testing.T) {
	mockSvc := &MockUserService{
		AddUserFunc: func(ctx context.Context, user models.UserAuth) (uuid.UUID, error) {
			return uuid.Nil, errors.New("db error")
		},
	}

	jwtSvc := jwt.NewJWTService("test-secret", 3600)
	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewAuthHandler(mockSvc, jwtSvc, l, tracer, m)

	router := setupTestRouter()
	router.POST("/sign-up", handler.SignUp)

	body := models.UserAuth{Name: "test", Email: "test@test.com", Pass: "password"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/sign-up", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestSignIn_Success(t *testing.T) {
	testUser := models.UserRequest{
		Id:    uuid.New(),
		Name:  "test",
		Email: "test@test.com",
	}

	mockSvc := &MockUserService{
		CheckAndGetUserFunc: func(ctx context.Context, user models.UserAuth) (models.UserRequest, error) {
			return testUser, nil
		},
	}

	jwtSvc := jwt.NewJWTService("test-secret", 3600)
	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewAuthHandler(mockSvc, jwtSvc, l, tracer, m)

	router := setupTestRouter()
	router.POST("/sign-in", handler.SignIn)

	body := models.UserAuth{Email: "test@test.com", Pass: "password"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/sign-in", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]models.UserRequest
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["user"].Id != testUser.Id {
		t.Errorf("expected user id %v, got %v", testUser.Id, response["user"].Id)
	}
}

func TestSignIn_Unauthorized(t *testing.T) {
	mockSvc := &MockUserService{
		CheckAndGetUserFunc: func(ctx context.Context, user models.UserAuth) (models.UserRequest, error) {
			return models.UserRequest{}, errors.New("invalid credentials")
		},
	}

	jwtSvc := jwt.NewJWTService("test-secret", 3600)
	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewAuthHandler(mockSvc, jwtSvc, l, tracer, m)

	router := setupTestRouter()
	router.POST("/sign-in", handler.SignIn)

	body := models.UserAuth{Email: "test@test.com", Pass: "wrong"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/sign-in", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}