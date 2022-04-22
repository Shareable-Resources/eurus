-- Deploy eurus-backend-db-report:2021-08-30_report_audit_cal_deposit_transactions_upsert_alter_where to pg

BEGIN;

-- XXX Add DDLs here.
CREATE OR REPLACE FUNCTION public.report_audit_cal_deposit_transactions_upsert(var_selected_date date)
 RETURNS integer
 LANGUAGE plpgsql
AS $function$
	declare row_count integer;
begin
	insert into report_audit (
			 asset_name,
			 selected_date,
			 mainnet_to_side_chain_count,
			 mainnet_to_side_chain_sum,
			 last_modified_date
	)
	select 
		dt.asset_name,
		var_selected_date as selected_date,
		COUNT(*) as mainnet_to_side_chain_count,
		SUM(amount) as mainnet_to_side_chain_sum,
		now() as last_modified_date
	from 
		deposit_transactions dt
		inner join assets a on dt.asset_name =a.asset_name 
		where created_date::date=var_selected_date and dt.status=40
		group by dt.asset_name
	on
		conflict on
		constraint report_audit_pk 
	do
		update
			set
				mainnet_to_side_chain_sum = EXCLUDED.mainnet_to_side_chain_sum,
				mainnet_to_side_chain_count = EXCLUDED.mainnet_to_side_chain_count,
				last_modified_date = now();
	GET DIAGNOSTICS row_count = ROW_COUNT;
   	RAISE NOTICE 'Row affected(report_audit_cal_deposit_transactions_upsert): % at %: ', row_count, now(); 
    return (select row_count);
   
end;

$function$
;

COMMIT;
