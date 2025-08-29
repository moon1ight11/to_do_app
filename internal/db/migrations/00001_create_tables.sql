-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    users (
        id UUID PRIMARY KEY,
        created_at TIMESTAMPTZ DEFAULT now (),
        updated_at TIMESTAMPTZ,
        deleted_at TIMESTAMPTZ name VARCHAR NOT NULL UNIQUE,
        name VARCHAR NOT NULL UNIQUE,
        email VARCHAR NOT NULL UNIQUE,
        pass VARCHAR NOT NULL,
    );

CREATE TABLE
    settings (
        user_id UUID REFERENCES users (id),
        default_tz VARCHAR,
        default_duration TIME,
    );

CREATE TABLE
    tasks (
        id UUID PRIMARY KEY,
        subtask_id VARCHAR REFERENCES tasks (id),
        owner_id UUID REFERENCES users (id),
        created_at TIMESTAMPTZ DEFAULT now (),
        updated_at TIMESTAMPTZ,
        deleted_at TIMESTAMPTZ,
        start_at TIMESTAMPTZ,
        end_at TIMESTAMPTZ,
        completed_at TIMESTAMPTZ,
        title VARCHAR,
        description VARCHAR,
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;

DROP TABLE tasks;

DROP TABLE settings;
-- +goose StatementEnd