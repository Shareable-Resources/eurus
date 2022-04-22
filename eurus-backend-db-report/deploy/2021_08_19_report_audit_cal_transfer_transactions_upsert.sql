-- Deploy eurus-backend-db-report:2021_08_19_report_audit_cal_transfer_transactions_upsert to pg

BEGIN;

-- XXX Add DDLs here.
CREATE OR REPLACE FUNCTION public.report_audit_cal_transfer_transactions_upsert(var_selected_date date, var_chain integer)
 RETURNS integer
 LANGUAGE plpgsql
AS $function$
	declare row_count integer;
begin
	insert into report_audit (
			 asset_name,
			 selected_date,
			 side_chain_to_side_chain_count,
			 side_chain_to_side_chain_sum,
			 last_modified_date
	)
	select 
		tt.asset_name,
		var_selected_date as selected_date,
		COUNT(*) as side_chain_to_side_chain_count,
		SUM(amount) as side_chain_to_side_chain_sum,
		now() as last_modified_date
	from 
		transfer_transactions tt
		inner join assets a on tt.asset_name =a.asset_name 
		where created_date::date=var_selected_date and "chain"=var_chain
		group by tt.asset_name
	on
		conflict on
		constraint report_audit_pk 
	do
		update
			set
				side_chain_to_side_chain_count = EXCLUDED.side_chain_to_side_chain_count,
				side_chain_to_side_chain_sum = EXCLUDED.side_chain_to_side_chain_sum,
				last_modified_date = now();
	GET DIAGNOSTICS row_count = ROW_COUNT;
   	RAISE NOTICE 'Row affected(report_audit_cal_transfer_transactions_upsert): % at %: ', row_count, now(); 
    return (select row_count);
   
end;

$function$
;

COMMIT;

