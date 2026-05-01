package usersrepos

import (
	"go.opentelemetry.io/otel/trace"
	"todoapp/internal/storage"
)

type Repo struct {
	storage.DataBase
	tracer trace.Tracer
}

func NewUserRepo(db *storage.DataBase, tracer trace.Tracer) *Repo {
	return &Repo{
		DataBase: *db,
		tracer:   tracer,
	}
}
