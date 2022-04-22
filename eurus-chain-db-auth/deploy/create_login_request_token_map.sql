-- Deploy eurus-backend-db-auth:create_login_rquest_token_map to pg

BEGIN;

SET search_path to public;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS login_request_token_maps (
	id SERIAL PRIMARY KEY,
    login_request_token VARCHAR(50),
	token TEXT,
	user_id VARCHAR(1024),
	expired_time TIMESTAMP with TIME zone,
	created_date TIMESTAMP with TIME zone,
	last_modified_date TIMESTAMP with time zone 
);

COMMIT;
