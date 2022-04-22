-- Deploy eurus-backend-db-user:create_user_devices_table to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS user_devices
(
    customer_id         BIGINT                      NOT NULL,
    customer_type       SMALLINT                    NOT NULL,
    device_id           VARCHAR                     NOT NULL,
    pub_key             VARCHAR                     NOT NULL,
    created_date        TIMESTAMP WITH TIME ZONE    NOT NULL,
    last_modified_date  TIMESTAMP WITH TIME ZONE    NOT NULL,
    PRIMARY KEY (customer_id, customer_type, device_id)
);

COMMIT;
