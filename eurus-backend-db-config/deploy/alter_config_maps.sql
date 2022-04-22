-- Deploy config:alter_config_maps to pg

BEGIN;

SET search_path to public;

ALTER TABLE config_maps RENAME COLUMN server_id TO id;

ALTER TABLE config_maps ADD COLUMN is_service BOOLEAN DEFAULT TRUE;


COMMIT;