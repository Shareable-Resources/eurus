-- Deploy eurus-backend-db-report:20211207_alter_wallet_balances_primary_key to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE wallet_balances DROP constraint wallet_balances_pkey;

ALTER TABLE wallet_balances ADD CONSTRAINT wallet_balances_pkey PRIMARY KEY(wallet_type, wallet_address, asset_name, mark_date, chain_id);

COMMIT;
