-- Verify eurus-backend-db-report:20220214_add_wallet_balance_configs on pg

BEGIN;

SELECT * FROM wallet_balance_configs LIMIT 1;

ROLLBACK;
