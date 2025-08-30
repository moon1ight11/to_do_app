-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    users (
        id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        created_at  TIMESTAMPTZ DEFAULT now(),
        updated_at  TIMESTAMPTZ,
        deleted_at  TIMESTAMPTZ,
        name        VARCHAR NOT NULL UNIQUE,
        email       VARCHAR NOT NULL UNIQUE,
        pass        VARCHAR NOT NULL
    );

CREATE TABLE
    settings (
        user_id             UUID PRIMARY KEY,
        default_tz          VARCHAR,
        default_duration    TIME,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );

CREATE TABLE
    tasks (
        id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        parent_task_id  UUID,
        owner_id        UUID,
        created_at      TIMESTAMPTZ DEFAULT now(),
        updated_at      TIMESTAMPTZ,
        deleted_at      TIMESTAMPTZ,
        start_at        TIMESTAMPTZ,
        end_at          TIMESTAMPTZ,
        completed_at    TIMESTAMPTZ,
        title           VARCHAR,
        description     VARCHAR,
        FOREIGN KEY (parent_task_id) REFERENCES tasks(id),
        FOREIGN KEY (owner_id) REFERENCES users(id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tasks;
DROP TABLE settings;
DROP TABLE users;
-- +goose StatementEnd