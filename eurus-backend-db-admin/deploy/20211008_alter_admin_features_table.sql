-- Deploy eurus-backend-db-admin:20211008_alter_admin_features_table to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE admin_features ADD COLUMN is_enabled boolean NOT NULL DEFAULT true;

COMMIT;
