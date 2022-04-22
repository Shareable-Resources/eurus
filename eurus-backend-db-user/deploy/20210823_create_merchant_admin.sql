-- Deploy eurus-backend-db-user:20210823_create_merchant_admin to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS merchant_admins (
	operator_id SERIAL PRIMARY KEY,
	merchant_id BIGINT NOT NULL,
	username VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	password_hash VARCHAR(300),
	status INT NOT NULL,
	created_date TIMESTAMPTZ  NOT NULL,
	last_modified_date TIMESTAMPTZ  NOT NULL
);

CREATE UNIQUE INDEX merchant_admins_idx1 ON merchant_admins (merchant_id, username);

COMMIT;
