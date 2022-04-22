-- Deploy eurus-backend-db-report:2021_08_17_create_report_audit_table to pg

BEGIN;

-- XXX Add DDLs here.
create table report_audit (
    asset_name VARCHAR(100) NOT NULL,
	mainnet_hot_balance_change NUMERIC(78),
	mainnet_cold_balance_change NUMERIC(78),
	side_chain_to_mainnet_admin_fee NUMERIC(85),
	side_chain_to_mainnet_count INT,
	side_chain_to_mainnet_sum NUMERIC(85),
	mainnet_to_side_chain_count INT,
	mainnet_to_side_chain_sum NUMERIC(85),
	side_chain_to_side_chain_count INT,
	side_chain_to_side_chain_sum NUMERIC(85),
	selected_date DATE NOT NULL,
	created_date TIMESTAMP with time zone,
	last_modified_date TIMESTAMP with time zone,
CONSTRAINT report_audit_pk PRIMARY KEY(asset_name, selected_date)

);
comment on column report_audit.asset_name is 'The asset (crypto currency) name';
comment on column report_audit.mainnet_hot_balance_change is 'The difference between a specific asset in hot wallet address in mainnet between ([selected_date]- 1 Day) and ([selected_date]- 2 Day)';
comment on column report_audit.mainnet_cold_balance_change is 'The difference between a specific asset in cold wallet address in mainnet between ([selected_date]- 1 Day) and ([selected_date]- 2 Day)';
comment on column report_audit.side_chain_to_mainnet_admin_fee is 'Calculated based on table [withdraw_transactions].[admin_fee], the value should be the sum of [admin_fee] group by [asset_name] and the [created_date] should be the selected date';
comment on column report_audit.side_chain_to_mainnet_count is  'Calculated based on table [withdraw_transactions], number of occurence of the selected date';
comment on column report_audit.side_chain_to_mainnet_sum is  'Calculated based on table [withdraw_transactions], total amount in of selected date ';


comment on column report_audit.mainnet_to_side_chain_count is  'Calculated based on table [deposit_transactions], number of occurence of the selected date';
comment on column report_audit.mainnet_to_side_chain_sum is  'Calculated based on table [deposit_transactions], total amount in of selected date ';

comment on column report_audit.side_chain_to_side_chain_count is  'Calculated based on table [transfer_transactions], number of occurence of the selected date';
comment on column report_audit.side_chain_to_side_chain_sum is  'Calculated based on table [transfer_transactions], total amount in of selected date ';

comment on column report_audit.selected_date is 'The report selected date';
comment on column report_audit.created_date is 'Created date of this report';
comment on column report_audit.last_modified_date is 'Last modified date of this report';

COMMIT;
