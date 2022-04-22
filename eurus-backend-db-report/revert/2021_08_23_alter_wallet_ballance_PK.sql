-- Revert eurus-backend-db-report:2021_08_23_alter_wallet_ballance_PK from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE public.wallet_balances DROP CONSTRAINT wallet_balances_pkey;
ALTER TABLE wallet_balances ADD PRIMARY KEY (wallet_address,created_date,asset_name);
ALTER TABLE public.wallet_balances ALTER COLUMN created_date TYPE TIMESTAMP with time zone;
COMMIT;
