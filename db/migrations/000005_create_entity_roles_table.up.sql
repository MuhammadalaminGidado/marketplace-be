CREATE TABLE entity_roles (
    entity_id BIGINT NOT NULL REFERENCES entities(id) ON DELETE CASCADE,
    role_id   BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,

    PRIMARY KEY (entity_id, role_id)
);