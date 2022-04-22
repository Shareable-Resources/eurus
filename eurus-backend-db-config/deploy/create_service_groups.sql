-- Deploy config:create_service_groups to pg

BEGIN;

SET search_path to public;

-- XXX Add DDLs here.
CREATE TABLE service_groups (
	group_id BIGINT,
	service_id BIGINT,
    PRIMARY KEY(group_id, service_id)
);

COMMIT;
