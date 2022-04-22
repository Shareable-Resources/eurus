-- Revert config:20210903_create_primary_key_assets from pg

BEGIN;

-- XXX Add DDLs here.

ALTER TABLE assets DROP CONSTRAINT assets_pk;
ALTER TABLE assets ADD CONSTRAINT assets_pk PRIMARY KEY ("currency_id");

COMMIT;
