package usershandlers

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

func TestGetUser_Success(t *testing.T) {
	userId := uuid.New()
	expectedUser := models.UserRequest{Id: userId, Name: "test", Email: "test@test.com"}

	mockSvc := &MockUserService{
		GetUserFunc: func(ctx context.Context, userId uuid.UUID) (models.UserRequest, error) {
			return expectedUser, nil
		},
	}

	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewUserHandler(mockSvc, l, nil, tracer, m)

	router := setupTestRouter()
	router.GET("/users", func(c *gin.Context) {
		c.Set("UserId", userId)
		handler.GetUser(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]models.UserRequest
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["user"].Id != expectedUser.Id {
		t.Errorf("expected user id %v, got %v", expectedUser.Id, response["user"].Id)
	}
}

func TestGetUser_Forbidden(t *testing.T) {
	mockSvc := &MockUserService{}
	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewUserHandler(mockSvc, l, nil, tracer, m)

	router := setupTestRouter()
	router.GET("/users", handler.GetUser)

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	userId := uuid.New()
	newName := "updated"
	updatedUser := models.UserRequest{Id: userId, Name: newName, Email: "test@test.com"}

	mockSvc := &MockUserService{
		UpdateUserFunc: func(ctx context.Context, name *string, pass *string, email *string, userId uuid.UUID) error {
			return nil
		},
		GetUserFunc: func(ctx context.Context, userId uuid.UUID) (models.UserRequest, error) {
			return updatedUser, nil
		},
	}

	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewUserHandler(mockSvc, l, nil, tracer, m)

	router := setupTestRouter()
	router.PATCH("/users", func(c *gin.Context) {
		c.Set("UserId", userId)
		handler.UpdateUser(c)
	})

	body := models.UserUpdate{Name: &newName}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPatch, "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDeleteUser_Success(t *testing.T) {
	userId := uuid.New()

	mockSvc := &MockUserService{
		DeleteUserFunc: func(ctx context.Context, userId uuid.UUID) error {
			return nil
		},
	}

	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewUserHandler(mockSvc, l, nil, tracer, m)

	router := setupTestRouter()
	router.DELETE("/users", func(c *gin.Context) {
		c.Set("UserId", userId)
		handler.DeleteUser(c)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDeleteUser_InternalError(t *testing.T) {
	userId := uuid.New()

	mockSvc := &MockUserService{
		DeleteUserFunc: func(ctx context.Context, userId uuid.UUID) error {
			return errors.New("db error")
		},
	}

	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewUserHandler(mockSvc, l, nil, tracer, m)

	router := setupTestRouter()
	router.DELETE("/users", func(c *gin.Context) {
		c.Set("UserId", userId)
		handler.DeleteUser(c)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
