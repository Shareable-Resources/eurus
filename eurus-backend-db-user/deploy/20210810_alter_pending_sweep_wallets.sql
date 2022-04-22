-- Deploy eurus-backend-db-user:20210810_alter_pending_sweep_wallets to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE public.pending_sweep_wallets ADD COLUMN IF NOT EXISTS previous_gas_fee_cap NUMERIC(78);
ALTER TABLE public.pending_sweep_wallets ADD COLUMN IF NOT EXISTS previous_gas_limit BIGINT;

COMMIT;
