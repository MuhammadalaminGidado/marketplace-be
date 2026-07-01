-- Drop indexes first (optional but clean)
DROP INDEX IF EXISTS idx_business_slug;
DROP INDEX IF EXISTS idx_business_entity_id;
DROP INDEX IF EXISTS idx_individuals_entity_id;
DROP INDEX IF EXISTS idx_entities_deleted_at;
DROP INDEX IF EXISTS idx_entities_email;

-- Drop tables in dependency order
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS business;
DROP TABLE IF EXISTS individuals;
DROP TABLE IF EXISTS entities;

-- Drop enum types AFTER tables that use them
DROP TYPE IF EXISTS address_owner;
DROP TYPE IF EXISTS business_status;
DROP TYPE IF EXISTS entity_status;