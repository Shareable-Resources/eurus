-- Deploy eurus-backend-db-user:exchange_rate_add_auto to pg

BEGIN;

-- XXX Add DDLs here.
alter table exchange_rates add column auto_update BOOLEAN DEFAULT TRUE;

COMMIT;
