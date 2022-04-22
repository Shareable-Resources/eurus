-- Deploy eurus-backend-db-report:2021_08_18_report_audit_cal_bal_change_upsert to pg

BEGIN;

-- XXX Add DDLs here.
CREATE OR REPLACE FUNCTION public.report_audit_cal_bal_change_upsert(var_selected_date date)
 RETURNS integer
 LANGUAGE plpgsql
AS $function$
	declare row_count integer;
	declare totalAffectedRows integer;
begin
	insert into report_audit (
			 asset_name,
			 selected_date,
			 mainnet_hot_balance_change,
			 created_date,
			 last_modified_date
	)
	select 
		p_asset_name as asset_name,
		p_selected_date as selected_date,
		p_balance_change as mainnet_hot_balance_change,
		p_created_date as created_date,
		p_last_modified_date as last_modified_date
	from 
		report_audit_cal_bal_change_select(var_selected_date,90) --Mainnet Hot Wallet (90)
	on
		conflict on
		constraint report_audit_pk 
	do
		update
			set
				mainnet_hot_balance_change = EXCLUDED.mainnet_hot_balance_change,
				last_modified_date = now();
	GET DIAGNOSTICS row_count = ROW_COUNT;
   	RAISE NOTICE 'Row affected(Mainnet Hot): % at %: ', row_count, now(); 
   	totalAffectedRows :=row_count;
   
   
	insert into report_audit (
			 asset_name,
			 selected_date,
			 mainnet_cold_balance_change,
			 created_date,
			 last_modified_date
	)
	select 
		p_asset_name as asset_name,
		p_selected_date as selected_date,
		p_balance_change as mainnet_cold_balance_change,
		p_created_date as created_date,
		p_last_modified_date as last_modified_date
	from 
		report_audit_cal_bal_change_select(var_selected_date, 91) --Mainnet Cold Wallet (91)
	on
		conflict on
		constraint report_audit_pk 
	do
		update
			set
				mainnet_cold_balance_change = EXCLUDED.mainnet_cold_balance_change,
				last_modified_date = now();
	GET DIAGNOSTICS row_count = ROW_COUNT;
	RAISE NOTICE 'Row affected(Mainnet Cold): % at %: ', row_count, now(); 
     	totalAffectedRows :=row_count + totalAffectedRows;
    RAISE NOTICE 'Row affected(report_audit_cal_bal_change_upsert) : % at %', now(),
totalAffectedRows; 
	RETURN (SELECT totalAffectedRows);
end;

$function$
;

COMMIT;
