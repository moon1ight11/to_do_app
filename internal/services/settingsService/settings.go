package settingsservice

import (
	"todoapp/internal/storage/repos/settingsrepos"

	"go.opentelemetry.io/otel/trace"
)

type SettingsService struct {
	settingsRepo *settingsrepos.Repo
	tracer       trace.Tracer
}

func NewSettingsService(settingsRepo *settingsrepos.Repo, tracer trace.Tracer) *SettingsService {
	return &SettingsService{
		settingsRepo: settingsRepo,
		tracer:       tracer,
	}
}
