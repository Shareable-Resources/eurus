-- Revert eurus-backend-db-report:alter_table_wallet_balances from pg

BEGIN;

-- XXX Add DDLs here.
-- Add chain_id as column
ALTER TABLE wallet_balances DROP CONSTRAINT wallet_balances_pkey;
ALTER TABLE wallet_balances DROP COLUMN chain_id;
ALTER TABLE wallet_balances ALTER COLUMN balance TYPE INT;
ALTER TABLE wallet_balances ALTER COLUMN user_id DROP DEFAULT;
ALTER TABLE wallet_balances ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE wallet_balances ALTER COLUMN user_id TYPE BIGINT;
ALTER TABLE wallet_balances ADD PRIMARY KEY (wallet_address,created_date,asset_name);

COMMIT;