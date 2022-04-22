-- Deploy eurus-backend-db-user:withdraw_transaction to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE withdraw_transaction (
    id integer PRIMARY KEY,
    customer_id bigint,
    customer_type smallint,
    asset_name Varchar(20),
    amount numeric(78),
    approval_wallet_address Varchar(100),
    request_trans_id bigint,
    request_trans_hash Varchar(255),
    review_date timestamptz,
    reviewed_by Varchar(255),
    review_trans_hash Varchar(255),
    innet_from_address Varchar(255),
    mainnet_from_address Varchar(255),
    mainnet_to_address Varchar(255),
    mainnet_trans_hash Varchar(255),
    mainnet_trans_date timestamptz,
    burn_trans_id bigint
);

-- CREATE UNIQUE INDEX idx_approval_wallet_address_request_trans_id ON withdraw_transaction(approval_wallet_address,request_trans_id);

COMMIT;
