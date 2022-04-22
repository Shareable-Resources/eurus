-- Revert eurus-backend-db-report:create_table_wallet_balances from pg

BEGIN;

DROP TABLE wallet_balances;

COMMIT;
