-- Deploy config:create_asset_settings to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS asset_settings
(
    id                              BIGSERIAL                   NOT NULL,
    asset_name                      VARCHAR(100)                NOT NULL,
    kyc0_max_daily_withdraw_amount  NUMERIC(78),
    kyc1_max_daily_withdraw_amount  NUMERIC(78),
    kyc2_max_daily_withdraw_amount  NUMERIC(78),
    kyc3_max_daily_withdraw_amount  NUMERIC(78),
    sweep_trigger_amount            NUMERIC(78)                 NOT NULL,
    is_enabled                      BOOLEAN                     NOT NULL,
    created_date                    TIMESTAMP WITH TIME ZONE    NOT NULL,
    last_modified_date              TIMESTAMP WITH TIME ZONE    NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (asset_name)
);

COMMIT;
