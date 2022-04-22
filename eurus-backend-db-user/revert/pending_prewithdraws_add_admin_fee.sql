-- Revert eurus-backend-db-user:pending_prewithdraws_add_admin_fee from pg

BEGIN;

ALTER TABLE pending_prewithdraws DROP COLUMN admin_fee;

COMMIT;
