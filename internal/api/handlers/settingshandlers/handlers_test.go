package settingshandlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

type mockCache struct {
	getFunc    func(key string, dest interface{}) error
	setFunc    func(key string, value interface{}) error
	deleteFunc func(key string) error
}

func (m *mockCache) Get(ctx context.Context, key string, dest interface{}) error {
	return m.getFunc(key, dest)
}
func (m *mockCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return m.setFunc(key, value)
}
func (m *mockCache) Delete(ctx context.Context, key string) error {
	return m.deleteFunc(key)
}

func newMetrics() *metrics.Metrics {
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	return metrics.NewMetrics()
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestGetSettings_Success(t *testing.T) {
	userId := uuid.New()
	expectedSettings := models.Setting{TimeDuration: 3600, UserTz: "UTC+3"}

	mockSvc := &MockSettingsService{
		GetSettingsFunc: func(ctx context.Context, userId uuid.UUID) (models.Setting, error) {
			return expectedSettings, nil
		},
	}

	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewSettingsHandler(mockSvc, l, nil, tracer, m)

	router := setupTestRouter()
	router.GET("/settings", func(c *gin.Context) {
		c.Set("UserId", userId)
		handler.GetSettings(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/settings", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]models.Setting
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["settings"].TimeDuration != expectedSettings.TimeDuration {
		t.Errorf("expected duration %f, got %f", expectedSettings.TimeDuration, response["settings"].TimeDuration)
	}
}

func TestGetSettings_Forbidden(t *testing.T) {
	mockSvc := &MockSettingsService{}
	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewSettingsHandler(mockSvc, l, nil, tracer, m)

	router := setupTestRouter()
	router.GET("/settings", handler.GetSettings)

	req, _ := http.NewRequest(http.MethodGet, "/settings", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestGetSettings_InternalError(t *testing.T) {
	userId := uuid.New()

	mockSvc := &MockSettingsService{
		GetSettingsFunc: func(ctx context.Context, userId uuid.UUID) (models.Setting, error) {
			return models.Setting{}, errors.New("db error")
		},
	}

	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewSettingsHandler(mockSvc, l, nil, tracer, m)

	router := setupTestRouter()
	router.GET("/settings", func(c *gin.Context) {
		c.Set("UserId", userId)
		handler.GetSettings(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/settings", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestUpdateSettings_Success(t *testing.T) {
	userId := uuid.New()
	duration := 7200.0
	tz := "UTC+5"
	expectedSettings := models.Setting{TimeDuration: duration, UserTz: tz}

	mockSvc := &MockSettingsService{
		UpdateSettingsFunc: func(ctx context.Context, userId uuid.UUID, d *float64, tz *string) error {
			return nil
		},
		GetSettingsFunc: func(ctx context.Context, userId uuid.UUID) (models.Setting, error) {
			return expectedSettings, nil
		},
	}

	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewSettingsHandler(mockSvc, l, nil, tracer, m)

	router := setupTestRouter()
	router.PATCH("/settings", func(c *gin.Context) {
		c.Set("UserId", userId)
		handler.UpdateSettings(c)
	})

	body := models.UpdatedSettings{TimeDuration: &duration, UserTz: &tz}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPatch, "/settings", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
