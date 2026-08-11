package settingsservice

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace/noop"

	"todoapp/internal/storage/repos/settingsrepos"
)

func TestGetSettings_Success(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	expected := settingsrepos.Setting{TimeDuration: 3600, UserTz: "UTC+3"}

	mockRepo := &MockSettingsRepo{
		SettingsByIdFunc: func(ctx context.Context, userId uuid.UUID) (settingsrepos.Setting, error) {
			return expected, nil
		},
	}

	svc := NewSettingsService(mockRepo, noop.NewTracerProvider().Tracer("test"))

	result, err := svc.GetSettings(ctx, userId)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.TimeDuration != expected.TimeDuration {
		t.Errorf("expected duration %f, got %f", expected.TimeDuration, result.TimeDuration)
	}
	if result.UserTz != expected.UserTz {
		t.Errorf("expected tz %s, got %s", expected.UserTz, result.UserTz)
	}
}

func TestGetSettings_NotFound(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()

	mockRepo := &MockSettingsRepo{
		SettingsByIdFunc: func(ctx context.Context, userId uuid.UUID) (settingsrepos.Setting, error) {
			return settingsrepos.Setting{}, errors.New("not found")
		},
	}

	svc := NewSettingsService(mockRepo, noop.NewTracerProvider().Tracer("test"))

	_, err := svc.GetSettings(ctx, userId)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateSettings_InvalidTimezone(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	badTz := "Mars/Time"

	mockRepo := &MockSettingsRepo{}

	svc := NewSettingsService(mockRepo, noop.NewTracerProvider().Tracer("test"))

	err := svc.UpdateSettings(ctx, userId, nil, &badTz)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
