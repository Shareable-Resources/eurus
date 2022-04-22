-- Deploy eurus-backend-db-user:transfer_trans_add_gas_price to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE transfer_transactions ADD COLUMN gas_price NUMERIC(78);


COMMIT;
