-- Deploy eurus-backend-db-user:create_pending_sweep_wallets to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS pending_sweep_wallets
(
    id                      BIGSERIAL                   NOT NULL,
    user_id                 BIGINT,
    mainnet_wallet_address  VARCHAR(50)                 NOT NULL,
    asset_name              VARCHAR(100)                NOT NULL,
    created_date            TIMESTAMP WITH TIME ZONE    NOT NULL,
    last_modified_date      TIMESTAMP WITH TIME ZONE    NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (mainnet_wallet_address, asset_name)
);

COMMIT;
