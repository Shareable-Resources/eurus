-- Verify eurus-backend-db-user:transaction_indices_add_erc20 on pg

BEGIN;

SELECT currency_address FROM transaction_indices;

ROLLBACK;
