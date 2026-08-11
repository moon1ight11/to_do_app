package settingsservice

import "go.opentelemetry.io/otel/trace"

type SettingsService struct {
	settingsRepo settingsRepo
	tracer       trace.Tracer
}

func NewSettingsService(settingsRepo settingsRepo, tracer trace.Tracer) *SettingsService {
	return &SettingsService{
		settingsRepo: settingsRepo,
		tracer:       tracer,
	}
}
