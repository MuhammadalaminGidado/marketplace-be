CREATE TABLE permissions (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,   -- e.g. user.read, user.write
    description TEXT,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);