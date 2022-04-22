-- Deploy eurus-backend-db-user:transfer_transactions_alter_varchar_len to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE transfer_transactions ALTER from_address TYPE VARCHAR(255);
ALTER TABLE transfer_transactions ALTER tx_hash TYPE VARCHAR(255);

COMMIT;
