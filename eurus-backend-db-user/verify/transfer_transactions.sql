-- Verify eurus-backend-db-user:transfer_transactions on pg

BEGIN;

SELECT * from transfer_transactions limit 1;

ROLLBACK;
