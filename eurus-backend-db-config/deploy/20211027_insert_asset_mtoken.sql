-- Deploy config:20211027_insert_asset_mtoken to pg

BEGIN;


-- XXX Add DDLs here.
INSERT INTO assets ("decimal", asset_name, currency_id, auto_update) VALUES (6, 'USDM', 'usd-coin', true);
INSERT INTO assets ("decimal", asset_name, currency_id, auto_update) VALUES (18, 'ETHM', 'eth', false);
INSERT INTO assets ("decimal", asset_name, currency_id, auto_update) VALUES (18, 'BTCM', 'wrapped-bitcoin', true);

INSERT INTO exchange_rates(asset_name, rate, created_date, last_modified_date) VALUES ('ETHM', 1, now(), now());


INSERT INTO asset_settings (asset_name, kyc0_max_daily_withdraw_amount, kyc1_max_daily_withdraw_amount, kyc2_max_daily_withdraw_amount, kyc3_max_daily_withdraw_amount, sweep_trigger_amount, is_enabled, created_date, last_modified_date) VALUES
('USDM', 5000000000, 4700000000000, NULL, NULL,  100000000, true, now(), now()),
('ETHM', 2000000000000000000, 1356000000000000000000, NULL, NULL, 1000000000000000, true, now(), now()),
('BTCM', 1000000000000000000, 102000000000000000000, NULL, NULL, 1000000000000000, true, now(), now());

COMMIT;
