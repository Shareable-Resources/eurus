-- Deploy eurus-backend-db-report:2021_08_20_alter_wallet_balances_wallet_type to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE wallet_balances ALTER COLUMN wallet_type TYPE SMALLINT USING (wallet_type::SMALLINT);
COMMENT ON COLUMN wallet_balances.wallet_type IS '[90]->Mainnet hot wallet address, obtained from EurusInternalConfig/[91]->Mainnet cold wallet address, obtained from table [users].[mainnet_wallet_address], [94]-> user wallet address, obtained from [users].[wallet_address], Other types please check wallet_background_indexer const([ApprovalObs](3),[Deposit](1),[WithdrawObs](2),[ConfigServer](8))';

COMMIT;
