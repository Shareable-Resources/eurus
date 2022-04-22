-- Deploy eurus-backend-db-user:users_alter_id to pg

BEGIN;

-- XXX Add DDLs here.
DROP TABLE users;

CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	login_address VARCHAR(50) NOT NULL,
	wallet_address VARCHAR(50),
	mainnet_wallet_address	 VARCHAR(50),
	email VARCHAR(50),
	kyc_status SMALLINT DEFAULT(0) NOT NULL,
	status SMALLINT NOT NULL,
	created_date TIMESTAMP WITH TIME ZONE NOT NULL,
	last_modified_date TIMESTAMP WITH TIME ZONE NOT NULL,
	is_metamask_addr BOOLEAN NOT NULL,
	last_login_time	TIMESTAMP WITH TIME ZONE
);

COMMIT;
