-- Deploy eurus-backend-db-user:transfer_transactions_add_status to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE transfer_transactions ADD COLUMN IF NOT EXISTS status smallint NOT NULL DEFAULT(50);


COMMIT;
