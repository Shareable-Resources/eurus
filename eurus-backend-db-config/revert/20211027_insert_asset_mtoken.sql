-- Revert config:20211027_insert_asset_mtoken from pg

BEGIN;

-- XXX Add DDLs here.
DELETE from assets where asset_name in ('USDM', 'BTCM', 'ETHM');

DELETE from exchange_rates where asset_name = 'ETHM';

DELETE from asset_settings where  asset_name in ('USDM', 'BTCM', 'ETHM');

COMMIT;
