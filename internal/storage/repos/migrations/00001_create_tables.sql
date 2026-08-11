-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS todo_app;

CREATE TABLE
    todo_app.users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        created_at TIMESTAMPTZ DEFAULT now(),
        updated_at TIMESTAMPTZ,
        name VARCHAR NOT NULL UNIQUE,
        email VARCHAR NOT NULL UNIQUE,
        pass VARCHAR NOT NULL
    );

CREATE INDEX idx_users_name ON todo_app.users (name);

CREATE INDEX idx_users_email ON todo_app.users (email);

CREATE TABLE
    todo_app.settings (
        user_id UUID PRIMARY KEY REFERENCES todo_app.users (id),
        default_tz VARCHAR DEFAULT 'UTC',
        default_duration NUMERIC DEFAULT 43200
    );

CREATE TABLE
    todo_app.tasks (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        parent_task_id UUID REFERENCES todo_app.tasks (id) ON DELETE CASCADE,
        owner_id UUID REFERENCES todo_app.users (id) ON DELETE CASCADE,
        created_at TIMESTAMPTZ DEFAULT now(),
        updated_at TIMESTAMPTZ,
        start_at TIMESTAMPTZ,
        end_at TIMESTAMPTZ,
        completed_at TIMESTAMPTZ,
        title VARCHAR,
        description VARCHAR
    );

CREATE INDEX idx_tasks_owner_id ON todo_app.tasks (owner_id);

CREATE INDEX idx_tasks_completed_at ON todo_app.tasks (completed_at);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE todo_app.tasks;

DROP TABLE todo_app.settings;

DROP TABLE todo_app.users;

DROP SCHEMA todo_app;

-- +goose StatementEnd