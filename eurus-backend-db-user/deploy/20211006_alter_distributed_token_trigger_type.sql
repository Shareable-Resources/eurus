-- Deploy eurus-backend-db-user:20211006_alter_distributed_token_trigger_type to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE distributed_tokens ADD COLUMN IF NOT EXISTS trigger_type int NOT NULL DEFAULT 10;
ALTER TABLE distributed_tokens ALTER COLUMN trigger_type SET DEFAULT 0;

COMMIT;
