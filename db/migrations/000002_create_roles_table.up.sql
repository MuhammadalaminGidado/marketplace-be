CREATE TABLE roles (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,   -- e.g. admin, user, manager
    description TEXT,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);