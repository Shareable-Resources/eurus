-- Deploy eurus-backend-db-report:2021_08_18_report_audit_insert_asset to pg

BEGIN;

-- XXX Add DDLs here.
create or replace
function public.report_audit_insert_asset(var_selected_date date)
 returns integer
 language plpgsql
as $function$
	declare row_count integer;

begin

insert
	into
	report_audit (
		asset_name,
		mainnet_hot_balance_change,
		mainnet_cold_balance_change,
		side_chain_to_mainnet_admin_fee,
		side_chain_to_mainnet_count,
		side_chain_to_mainnet_sum,
		mainnet_to_side_chain_count,
		mainnet_to_side_chain_sum,
		side_chain_to_side_chain_count,
		side_chain_to_side_chain_sum,
		selected_date,
		created_date,
		last_modified_date
	)
	select 
		asset_name,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		var_selected_date as selected_date,
		now() as created_date,
		now() as last_modified_date
from 
	    assets
	on
		conflict on
		constraint report_audit_pk 
	do
		update
set
				last_modified_date = now();

get diagnostics row_count = ROW_COUNT;

raise notice 'Row affected(report_audit_insert_asset): % at %: ',
row_count,
now();

return (
select
	row_count);
end;

$function$
;

COMMIT;
