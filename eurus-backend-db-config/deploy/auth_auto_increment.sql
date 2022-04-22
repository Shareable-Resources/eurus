-- Deploy config:config_maps_auto_increment to pg

BEGIN;


CREATE SEQUENCE IF NOT EXISTS auth_services_id_seq OWNED BY auth_services.id;

ALTER TABLE auth_services ALTER id SET DEFAULT nextval('auth_services_id_seq');

COMMIT;
