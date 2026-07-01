CREATE TYPE entity_status AS ENUM ('active', 'suspended', 'deactivated');
CREATE TYPE business_status AS ENUM ('pending', 'active', 'suspended');


CREATE TABLE entities (
    id              BIGSERIAL PRIMARY KEY,
    email           TEXT NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    phone           TEXT UNIQUE,
    avatar_url      TEXT,
    status          entity_status NOT NULL DEFAULT 'active',
    last_login_at   TIMESTAMPTZ,
    deleted_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE individuals (
    id              BIGSERIAL PRIMARY KEY,
    entity_id       BIGINT NOT NULL UNIQUE REFERENCES entities(id) ON DELETE CASCADE,

    first_name      TEXT NOT NULL,
    last_name       TEXT NOT NULL,
    date_of_birth   DATE,
    bio             TEXT,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE businesses (
    id              BIGSERIAL PRIMARY KEY,
    entity_id       BIGINT NOT NULL UNIQUE REFERENCES entities(id) ON DELETE CASCADE,

    business_name   TEXT NOT NULL,
    slug            TEXT NOT NULL UNIQUE,
    description     TEXT,
    logo_url        TEXT,
    status          business_status NOT NULL DEFAULT 'pending',
    is_verified     BOOLEAN NOT NULL DEFAULT FALSE,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE addresses (
    id              BIGSERIAL PRIMARY KEY,
    entity_id       BIGINT NOT NULL REFERENCES entities(id) ON DELETE CASCADE,

    line1           TEXT NOT NULL,
    line2           TEXT,
    city            TEXT NOT NULL,
    state           TEXT,
    postal_code     TEXT,
    country         CHAR(2) NOT NULL,

    is_default      BOOLEAN NOT NULL DEFAULT FALSE,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE INDEX idx_entities_email ON entities(email);
CREATE INDEX idx_entities_deleted_at ON entities(deleted_at) WHERE deleted_at IS NULL;

CREATE INDEX idx_individuals_entity_id ON individuals(entity_id);
CREATE INDEX idx_businesses_entity_id ON businesses(entity_id);

CREATE INDEX idx_businesses_slug ON businesses(slug);

CREATE INDEX idx_addresses_entity_id ON addresses(entity_id);


CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER trg_entities_updated_at
BEFORE UPDATE ON entities
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_individuals_updated_at
BEFORE UPDATE ON individuals
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_businesses_updated_at
BEFORE UPDATE ON businesses
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_addresses_updated_at
BEFORE UPDATE ON addresses
FOR EACH ROW EXECUTE FUNCTION set_updated_at();