-- Deploy config:config_maps to pg

BEGIN;

SET search_path to public;

-- XXX Add DDLs here.
CREATE TABLE config_maps (
	server_id BIGINT,
	key VARCHAR,
    value VARCHAR,
    PRIMARY KEY(server_id, key)
);

COMMIT;
