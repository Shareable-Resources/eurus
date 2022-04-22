-- Deploy eurus-backend-db-user:20211201_alter_report_audit_table to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE report_audit 
ADD COLUMN IF NOT EXISTS mainnet_cold_wallet_current_balance numeric(78),
ADD COLUMN IF NOT EXISTS mainnet_cold_wallet_previous_balance NUMERIC(78),
ADD COLUMN IF NOT EXISTS mainnet_hot_wallet_current_balance numeric(78),
ADD COLUMN IF NOT EXISTS mainnet_hot_wallet_previous_balance NUMERIC(78);

COMMIT; 
