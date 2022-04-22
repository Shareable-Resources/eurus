-- Deploy eurus-backend-db-report:alter_table_wallet_balances to pg


BEGIN;

-- XXX Add DDLs here.
-- Add chain_id as column
ALTER TABLE wallet_balances 
ADD COLUMN chain_id SMALLINT NOT NULL;
COMMENT ON COLUMN wallet_balances.chain_id IS 'Ethernet chain id';
-- Remove constraint and add constraint with new added column chain_id
ALTER TABLE wallet_balances DROP CONSTRAINT wallet_balances_pkey;
ALTER TABLE wallet_balances ALTER COLUMN balance TYPE NUMERIC(78);
ALTER TABLE wallet_balances ALTER COLUMN user_id DROP DEFAULT;
ALTER TABLE wallet_balances ALTER COLUMN user_id DROP NOT NULL;
ALTER TABLE wallet_balances ALTER COLUMN user_id TYPE BIGINT;
ALTER table wallet_balances ADD PRIMARY KEY (wallet_address,chain_id,created_date,asset_name);

COMMIT;
