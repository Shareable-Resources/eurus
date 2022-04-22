-- Deploy config:auth_services to pg

BEGIN;

SET search_path to public;

-- XXX Add DDLs here.
CREATE TABLE auth_services (
	id BIGINT,
	service_name VARCHAR,
    pub_key VARCHAR,
    PRIMARY KEY(id)
);

COMMIT;
