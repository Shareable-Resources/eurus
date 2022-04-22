-- Revert eurus-backend-db-user:transfer_trans_add_gas_price from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE transfer_transactions DROP COLUMN gas_price;

COMMIT;
