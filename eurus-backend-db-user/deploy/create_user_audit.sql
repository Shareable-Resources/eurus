-- Deploy eurus-backend-db-user:create_user_audit to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS user_audits (
    id SERIAL PRIMARY KEY,
    user_id BIGINT,
	login_address VARCHAR(50) NOT NULL,
	wallet_address VARCHAR(50),
	mainnet_wallet_address	 VARCHAR(50),
    owner_wallet_address VARCHAR(50),
	email VARCHAR(50),
	kyc_status SMALLINT DEFAULT(0) NOT NULL,
    status smallint NOT NULL,
	created_date TIMESTAMP WITH TIME ZONE NOT NULL,
	last_modified_date TIMESTAMP WITH TIME ZONE NOT NULL,
	is_metamask_addr BOOLEAN NOT NULL
);

COMMIT;
