-- Deploy eurus-backend-db-user:pending_prewithdraw to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE pending_prewithdraws (
    id serial PRIMARY KEY,
    customer_id bigint,
    customer_type smallint,
    innet_from_address Varchar(100),
    mainnet_to_address Varchar(100),
    approval_wallet_address Varchar(100),
    request_trans_id bigint,
    request_trans_hash Varchar(255),
    asset_name Varchar(20),
    amount numeric(78),
    status smallint,
    created_date timestamptz,
    last_modified_date timestamptz
);

CREATE UNIQUE INDEX idx_approval_wallet_address_request_trans_id ON pending_prewithdraws(approval_wallet_address, request_trans_id);




COMMIT;
