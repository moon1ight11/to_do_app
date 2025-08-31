-- +goose Up
-- +goose StatementBegin

CREATE SCHEMA IF NOT EXISTS todo_app;

CREATE TABLE
    todo_app.users (
        id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        created_at  TIMESTAMPTZ DEFAULT now(),
        updated_at  TIMESTAMPTZ,
        deleted_at  TIMESTAMPTZ,
        name        VARCHAR NOT NULL UNIQUE,
        email       VARCHAR NOT NULL UNIQUE,
        pass        VARCHAR NOT NULL
    );

CREATE TABLE
    todo_app.settings (
        user_id             UUID PRIMARY KEY REFERENCES users(id),
        default_tz          VARCHAR,
        default_duration    INTERVAL
    );

CREATE TABLE
    todo_app.tasks (
        id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        parent_task_id  UUID REFERENCES tasks(id),
        owner_id        UUID REFERENCES users(id),
        created_at      TIMESTAMPTZ DEFAULT now(),
        updated_at      TIMESTAMPTZ,
        deleted_at      TIMESTAMPTZ,
        start_at        TIMESTAMPTZ,
        end_at          TIMESTAMPTZ,
        completed_at    TIMESTAMPTZ,
        title           VARCHAR,
        description     VARCHAR
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE todo_app.tasks;
DROP TABLE todo_app.settings;
DROP TABLE todo_app.users;
DROP SCHEMA todo_app;
-- +goose StatementEnd