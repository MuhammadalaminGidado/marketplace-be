DROP TRIGGER IF EXISTS trg_entities_updated_at ON entities;
DROP FUNCTION IF EXISTS set_updated_at();

DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS entities;

DROP TYPE IF EXISTS entity_status;

DROP EXTENSION IF EXISTS pgcrypto;