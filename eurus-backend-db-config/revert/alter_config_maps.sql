-- Revert config:alter_config_maps from pg

BEGIN;

SET search_path to public;

ALTER TABLE config_maps RENAME COLUMN id TO server_id;

ALTER TABLE config_maps DROP COLUMN is_service;


COMMIT;