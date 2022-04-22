-- Verify eurus-backend-db-user:20211201_alter_report_audit_table on pg

BEGIN;

-- XXX Add verifications here.
SELECT mainnet_hot_wallet_current_balance, mainnet_hot_wallet_previous_balance, mainnet_cold_wallet_current_balance, mainnet_cold_wallet_previous_balance FROM report_audit limit 1;

ROLLBACK;
