-- Revert config:config_maps_auto_increment from pg

BEGIN;

ALTER SEQUENCE auth_services_id_seq OWNED BY NONE;

COMMIT;
