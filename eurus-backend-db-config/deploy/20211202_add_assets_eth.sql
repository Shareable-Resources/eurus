-- Deploy config:20211202_add_assets_eth to pg

BEGIN;

-- XXX Add DDLs here.

INSERT INTO assets ("decimal", asset_name, currency_id, auto_update) VALUES (18, 'ETH', 'eth', false);

COMMIT;
