-- Verify config:20211027_insert_asset_mtoken on pg

BEGIN;

-- XXX Add verifications here.

select 1/count(*) from assets where asset_name = 'USDM';
select 1/count(*) from assets where asset_name = 'ETHM';
select 1/count(*) from assets where asset_name = 'BTCM';

select 1/count(*) from exchange_rates where asset_name = 'ETHM';

select 1/count(*) from asset_settings where asset_name = 'USDM';
select 1/count(*) from asset_settings where asset_name = 'ETHM';
select 1/count(*) from asset_settings where asset_name = 'BTCM';


ROLLBACK;
