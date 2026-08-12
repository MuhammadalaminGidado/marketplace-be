DROP TABLE IF EXISTS otp_codes;

ALTER TABLE entities
    DROP COLUMN IF EXISTS email_verified_at;