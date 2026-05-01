package settingsrepos

import (
	"todoapp/internal/storage"

	"go.opentelemetry.io/otel/trace"
)

type Repo struct {
	storage.DataBase
	tracer trace.Tracer
}

func NewSettingsRepo(db *storage.DataBase, tracer trace.Tracer) *Repo {
	return &Repo{
		DataBase: *db,
		tracer:   tracer,
	}
}
