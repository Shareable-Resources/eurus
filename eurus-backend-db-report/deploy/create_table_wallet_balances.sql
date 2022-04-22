-- Deploy eurus-backend-db-user:create_table_wallet_balances to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS wallet_balances (
    wallet_type VARCHAR(2) NOT NULL,
    wallet_address VARCHAR(50) NOT NULL,
    asset_name VARCHAR(100) NOT NULL,
    balance INT,
    user_id SERIAL,
    created_date TIMESTAMP with time zone NOT NULL,
    PRIMARY KEY (wallet_address, asset_name, created_date)
);

COMMENT ON COLUMN wallet_balances.wallet_type IS '[MH]->Mainnet hot wallet address, obtained from EurusInternalConfig/[MC]->Mainnet cold wallet address, obtained from table [users].[mainnet_wallet_address], [S]-> user wallet address, obtained from [users].[wallet_address]';
COMMENT ON COLUMN wallet_balances.wallet_address IS 'Wallet address, can be hot wallet or cold wallet address. or user side chain wallet address';
COMMENT ON COLUMN wallet_balances.asset_name IS 'Asset name, i.e. BTC, ETH, USDT';
COMMENT ON COLUMN wallet_balances.balance IS 'The balance of this records base on wallet address and asset';
COMMENT ON COLUMN wallet_balances.user_id IS 'User id of the user';
COMMENT ON COLUMN wallet_balances.created_date IS 'Created Date';

COMMIT;

