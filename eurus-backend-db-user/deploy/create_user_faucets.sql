-- Deploy eurus-backend-db-user:create_user_faucets to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS user_faucets(
    trans_hash Varchar(255),
    key Varchar(20),
    user_id BIGINT,
    status SMALLINT,
    created_date TIMESTAMP WITH TIME ZONE NOT NULL,
    last_modified_date TIMESTAMP WITH TIME ZONE NOT NULL
);

COMMIT;
