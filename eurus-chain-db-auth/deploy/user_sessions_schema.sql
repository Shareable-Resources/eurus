-- Deploy eurus-backend:user_sessions_schema to pg

BEGIN;

SET search_path to public;

-- XXX Add DDLs here.
CREATE TABLE user_sessions (
	token VARCHAR(256) PRIMARY KEY,
	service_id BIGINT,
	user_id VARCHAR(1024),
	expired_time TIMESTAMP with TIME zone,
	created_date TIMESTAMP with TIME zone,
	last_modified_date TIMESTAMP with time zone,
	type SMALLINT default 0
);


COMMIT;
