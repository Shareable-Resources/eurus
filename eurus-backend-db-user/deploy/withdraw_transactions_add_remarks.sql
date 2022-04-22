-- Deploy eurus-backend-db-user:withdraw_transactions_add_remarks to pg

BEGIN;

-- XXX Add DDLs here.
	ALTER TABLE withdraw_transactions ADD COLUMN remarks TEXT;
	
COMMIT;
