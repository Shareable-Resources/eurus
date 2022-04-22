-- Deploy eurus-backend-db-user:create_user_asset_setting to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE user_asset_settings (
    id SERIAL PRIMARY KEY,
    asset_name TEXT,
    kyc0_max_daily_withdraw_amount BIGINT,
    kyc1_max_daily_withdraw_amount BIGINT,
    kyc2_max_daily_withdraw_amount BIGINT,
    kyc3_max_daily_withdraw_amount BIGINT,
    is_enabled BOOLEAN,
    created_date TIMESTAMP WITH TIME ZONE NOT NULL,
	last_modified_date TIMESTAMP WITH TIME ZONE NOT NULL
);

COMMIT;
