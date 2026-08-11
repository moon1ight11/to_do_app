package tasksrepos

import (
	"database/sql"

	"go.opentelemetry.io/otel/trace"
	"todoapp/internal/storage"
)

type Repo struct {
	storage.DataBase
	tracer trace.Tracer
}

func NewTasksRepo(db *storage.DataBase, tracer trace.Tracer) *Repo {
	return &Repo{
		DataBase: *db,
		tracer:   tracer,
	}
}

func (r *Repo) DB() *sql.DB {
	return r.DataBase.DB
}