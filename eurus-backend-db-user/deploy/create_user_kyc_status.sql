-- Deploy eurus-backend-db-user:create_user_kyc_status to pg

BEGIN;

SET search_path to public;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS user_kyc_status (
    id SERIAL NOT NULL,
    user_id BIGINT NOT NULL,   --Unique
    kyc_level BIGINT DEFAULT 0 NOT NULL, --Unique
    approval_date TIMESTAMP WITH TIME ZONE,
    operator_id BIGINT,
    kyc_retry_count INT DEFAULT 0 NOT NULL,
    kyc_country_code CHAR(2),
    created_date TIMESTAMP WITH TIME ZONE NOT NULL,
    last_modified_date TIMESTAMP WITH TIME ZONE NOT NULL,
    kyc_status  SMALLINT DEFAULT 0 NOT NULL,
    kyc_doc  SMALLINT NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT uc_user_kyc_status UNIQUE (user_id,kyc_level)
);

COMMENT ON COLUMN user_kyc_status.id IS 'Auto increment primay key representing the id of a kyc status';
COMMENT ON COLUMN user_kyc_status.user_id IS 'Foreign key to user table, representing the user who owned this kyc status';
COMMENT ON COLUMN user_kyc_status.kyc_level IS 'KYC Level, 0(Pending), 1(ID+Selfie/Password+Selfie verified)';
COMMENT ON COLUMN user_kyc_status.approval_date IS 'The date CS admin approve this KYC status';
COMMENT ON COLUMN user_kyc_status.operator_id IS 'The CS admin who do things to this record';
COMMENT ON COLUMN user_kyc_status.kyc_retry_count IS 'How many time the client user submit for kyc approval';
COMMENT ON COLUMN user_kyc_status.kyc_country_code IS 'Which country the client wants KYC verification';
COMMENT ON COLUMN user_kyc_status.created_date IS 'Created Date';
COMMENT ON COLUMN user_kyc_status.last_modified_date IS 'Last modified date';
COMMENT ON COLUMN user_kyc_status.kyc_status IS 'KYC status can be 0(Pending), 1(Waiting for Approval), 2(Waiting for resubmit), 3(Approved), 4(Rejected)';
COMMENT ON COLUMN user_kyc_status.kyc_doc IS 'KYC document can be (0)ID card, 1(Passport)';

COMMIT;
