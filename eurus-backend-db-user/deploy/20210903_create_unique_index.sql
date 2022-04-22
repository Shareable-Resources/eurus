-- Deploy eurus-backend-db-user:20210903_create_unique_index to pg

BEGIN;

-- XXX Add DDLs here.

--Add index
CREATE INDEX idx_user_status ON users ("status");
CREATE INDEX idx_deposit_transactions_customer_id_customer_type_asset_name_created_date ON deposit_transactions ("customer_id","customer_type","asset_name","created_date");
CREATE INDEX idx_withdraw_transactions_asset_name_request_date_customer_id_customer_type ON withdraw_transactions ("asset_name","request_date","customer_id","customer_type");
CREATE INDEX idx_transfer_transactions_asset_name_created_date_user_id ON transfer_transactions ("asset_name","created_date","user_id");
CREATE INDEX idx_user_faucets_transhash ON user_faucets ("trans_hash");
--Add Primary Key
ALTER TABLE verifications ADD CONSTRAINT verifications_pk PRIMARY KEY ("user_id", "type");
ALTER TABLE user_faucets ADD CONSTRAINT user_faucets_pk PRIMARY KEY ("user_id", "key");

COMMIT;


