-- Verify eurus-backend-db-report:create_table_wallet_balances on pg

BEGIN;

SELECT * FROM wallet_balances LIMIT 5;

ROLLBACK;
