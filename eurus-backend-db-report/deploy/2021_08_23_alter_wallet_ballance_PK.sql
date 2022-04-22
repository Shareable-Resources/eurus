-- Deploy eurus-backend-db-report:2021_08_23_alter_wallet_ballance_PK to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE public.wallet_balances DROP CONSTRAINT IF EXISTS wallet_balances_pkey;
ALTER TABLE public.wallet_balances ADD PRIMARY KEY (wallet_type,wallet_address,asset_name,created_date,chain_id);
ALTER TABLE public.wallet_balances ALTER COLUMN created_date TYPE DATE; 
COMMIT;
