-- Deploy eurus-backend-db-user:create_verification_table to pg

BEGIN;

CREATE TABLE IF NOT EXISTS verifications (
    user_id BIGINT,
    code VARCHAR(50),
    created_date TIMESTAMP WITH TIME ZONE NOT NULL,
    last_modified_date TIMESTAMP WITH TIME ZONE NOT NULL,
    expired_time TIMESTAMP WITH TIME ZONE,
    count SMALLINT DEFAULT 0,
    type SMALLINT DEFAULT 0
);

COMMIT;
