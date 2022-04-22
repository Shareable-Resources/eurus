-- Deploy eurus-backend-db-user:20210903_create_unique_index to pg

BEGIN;

-- XXX Add DDLs here.

--Add Primary Key
ALTER TABLE assets DROP CONSTRAINT assets_pkey;
ALTER TABLE assets ADD CONSTRAINT assets_pk PRIMARY KEY ("asset_name");

COMMIT;


