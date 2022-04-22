-- Revert eurus-backend-db-user:20210903_create_unique_index from pg

BEGIN;

-- XXX Add DDLs here.
DROP INDEX idx_user_status;
DROP INDEX idx_deposit_transactions_customer_id_customer_type_asset_name_created_date;
DROP INDEX idx_withdraw_transactions_asset_name_request_date_customer_id_customer_type;
DROP INDEX idx_transfer_transactions_asset_name_created_date_user_id;
DROP INDEX idx_user_faucets_transhash;
ALTER TABLE verifications DROP CONSTRAINT verifications_pk;
ALTER TABLE user_faucets DROP CONSTRAINT user_faucets_pk;

COMMIT;
