-- Revert eurus-backend-db-user:transfer_transactions from pg

BEGIN;

DROP TABLE transfer_transactions;

COMMIT;
