-- Revert eurus-backend-db-user:exchange_rate_add_auto from pg

BEGIN;

-- XXX Add DDLs here.
alter table exchange_rates drop column auto_update;

COMMIT;
