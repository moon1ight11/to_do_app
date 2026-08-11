package taskshandlers

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

func ptr[T any](v T) *T { return &v }

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestCreateTask_Success(t *testing.T) {
	userId := uuid.New()

	mockSvc := &MockTaskService{
		CreateTaskFunc: func(ctx context.Context, task models.Task) error {
			return nil
		},
	}

	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewTasksHandler(mockSvc, l, nil, tracer, m)

	router := setupTestRouter()
	router.POST("/tasks", func(c *gin.Context) {
		c.Set("UserId", userId)
		handler.CreateTask(c)
	})

	task := models.Task{Title: ptr("test task")}
	jsonBody, _ := json.Marshal(task)

	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestCreateTask_Forbidden(t *testing.T) {
	mockSvc := &MockTaskService{}
	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewTasksHandler(mockSvc, l, nil, tracer, m)

	router := setupTestRouter()
	router.POST("/tasks", handler.CreateTask)

	req, _ := http.NewRequest(http.MethodPost, "/tasks", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestCreateTask_InternalError(t *testing.T) {
	userId := uuid.New()

	mockSvc := &MockTaskService{
		CreateTaskFunc: func(ctx context.Context, task models.Task) error {
			return errors.New("db error")
		},
	}

	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewTasksHandler(mockSvc, l, nil, tracer, m)

	router := setupTestRouter()
	router.POST("/tasks", func(c *gin.Context) {
		c.Set("UserId", userId)
		handler.CreateTask(c)
	})

	task := models.Task{Title: ptr("test task")}
	jsonBody, _ := json.Marshal(task)

	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestGetTasks_Success(t *testing.T) {
	userId := uuid.New()
	taskId := uuid.New()
	expectedTasks := []models.Task{{Id: &taskId, Title: ptr("test")}}

	mockSvc := &MockTaskService{
		GetAllTasksFunc: func(ctx context.Context, userId uuid.UUID) ([]models.Task, error) {
			return expectedTasks, nil
		},
	}

	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewTasksHandler(mockSvc, l, nil, tracer, m)

	router := setupTestRouter()
	router.GET("/tasks", func(c *gin.Context) {
		c.Set("UserId", userId)
		handler.GetTasks(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string][]models.Task
	json.Unmarshal(w.Body.Bytes(), &response)
	if len(response["tasks"]) != 1 {
		t.Errorf("expected 1 task, got %d", len(response["tasks"]))
	}
}

func TestGetTaskById_Success(t *testing.T) {
	userId := uuid.New()
	taskId := uuid.New()
	expectedTask := models.Task{Id: &taskId, Title: ptr("test")}

	mockSvc := &MockTaskService{
		GetOneTaskFunc: func(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) (models.Task, error) {
			return expectedTask, nil
		},
	}

	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewTasksHandler(mockSvc, l, nil, tracer, m)

	router := setupTestRouter()
	router.GET("/tasks/:task_id", func(c *gin.Context) {
		c.Set("UserId", userId)
		handler.GetTaskById(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/tasks/"+taskId.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetTaskById_BadRequest_InvalidUUID(t *testing.T) {
	userId := uuid.New()

	mockSvc := &MockTaskService{}
	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewTasksHandler(mockSvc, l, nil, tracer, m)

	router := setupTestRouter()
	router.GET("/tasks/:task_id", func(c *gin.Context) {
		c.Set("UserId", userId)
		handler.GetTaskById(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/tasks/not-a-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteTask_Success(t *testing.T) {
	userId := uuid.New()
	taskId := uuid.New()

	mockSvc := &MockTaskService{
		DeleteTaskFunc: func(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) error {
			return nil
		},
	}

	l := &mockLogger{}
	tracer := noop.NewTracerProvider().Tracer("test")
	m := newMetrics()

	handler := NewTasksHandler(mockSvc, l, nil, tracer, m)

	router := setupTestRouter()
	router.DELETE("/tasks/:task_id", func(c *gin.Context) {
		c.Set("UserId", userId)
		handler.DeleteTask(c)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/tasks/"+taskId.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
