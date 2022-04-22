-- Deploy eurus-backend-db-report:20211201_create_block_timestamp_indices to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS block_timestamp_indices (
    mark_time TIMESTAMP with time zone not null,
    chain_id INT not null,
    block_number INT,
    block_date_time TIMESTAMP with time zone,
    block_timestamp INT,
    created_date TIMESTAMP with time zone NOT NULL,
    last_modified_date TIMESTAMP with time zone,
    CONSTRAINT block_timestamp_indices_pk PRIMARY KEY (mark_time, chain_id)
);

COMMIT;
