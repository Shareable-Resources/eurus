-- Deploy eurus-backend-db-report:2021_08_18_report_audit_cal_withdraw_transactions_upsert to pg

BEGIN;

-- XXX Add DDLs here.
CREATE OR REPLACE FUNCTION public.report_audit_cal_withdraw_transactions_upsert(var_selected_date date)
 RETURNS integer
 LANGUAGE plpgsql
AS $function$
	declare row_count integer;
begin
	insert into report_audit (
			 asset_name,
			 selected_date,
			 side_chain_to_mainnet_admin_fee,
			 side_chain_to_mainnet_count,
			 side_chain_to_mainnet_sum,
			 last_modified_date
	)
	select 
		wt.asset_name,
		var_selected_date as selected_date,
		SUM(admin_fee) as side_chain_to_mainnet_admin_fee,
		COUNT(*) as side_chain_to_mainnet_count,
		SUM(amount) as side_chain_to_mainnet_sum,
		now() as last_modified_date
	from 
		withdraw_transactions wt
		inner join assets a on wt.asset_name =a.asset_name 
		where created_date::date=var_selected_date
		group by wt.asset_name
	on
		conflict on
		constraint report_audit_pk 
	do
		update
			set
				side_chain_to_mainnet_admin_fee = EXCLUDED.side_chain_to_mainnet_admin_fee,
				side_chain_to_mainnet_count = EXCLUDED.side_chain_to_mainnet_count,
				side_chain_to_mainnet_sum = EXCLUDED.side_chain_to_mainnet_sum,
				last_modified_date = now();
	GET DIAGNOSTICS row_count = ROW_COUNT;
   	RAISE NOTICE 'Row affected(report_audit_cal_withdraw_transactions_upsert): % at %: ', row_count, now(); 
    return (select row_count);
   
end;

$function$
;

COMMIT;
