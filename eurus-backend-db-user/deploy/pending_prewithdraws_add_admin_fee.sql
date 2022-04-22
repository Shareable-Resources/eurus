-- Deploy eurus-backend-db-user:pending_prewithdraws_add_admin_fee to pg

BEGIN;

ALTER TABLE pending_prewithdraws ADD COLUMN IF NOT EXISTS admin_fee numeric(78) DEFAULT 0;

COMMIT;
