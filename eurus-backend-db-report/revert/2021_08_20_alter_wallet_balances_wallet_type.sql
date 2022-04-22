-- Revert eurus-backend-db-report:2021_08_20_alter_wallet_balances_wallet_type from pg

BEGIN;

-- XXX Add DDLs here.
-- Don' t revert to varchar(2) since data is not able to save
ALTER TABLE wallet_balances ALTER COLUMN wallet_type TYPE VARCHAR(15);

COMMIT;
