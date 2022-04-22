-- Deploy eurus-backend-db-user:withdraw_requests to pg

BEGIN;

CREATE TABLE withdraw_requests (
    withdraw_id bigint,
    request_trans_hash Varchar(255),
    service_id integer,
    data    bytea,
    signature bytea,
    created_date timestamptz,
    last_modified_date timestamptz
);

CREATE UNIQUE INDEX idx_request_trans_hash_service_id ON withdraw_requests(request_trans_hash, service_id);

COMMIT;
