-- Deploy eurus-backend-db-user:user_schema to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE users (
	id BIGINT PRIMARY KEY,
	login_address VARCHAR(50) NOT NULL,
	wallet_address VARCHAR(50),
	mainnet_wallet_address	 VARCHAR(50),
	email VARCHAR(50),
	kyc_status SMALLINT DEFAULT(0) NOT NULL,
	created_date TIMESTAMP WITH TIME ZONE NOT NULL,
	last_modified_date TIMESTAMP WITH TIME ZONE NOT NULL,
	is_metamask_addr BOOLEAN NOT NULL,
	last_login_time	TIMESTAMP WITH TIME ZONE
);

COMMIT;
