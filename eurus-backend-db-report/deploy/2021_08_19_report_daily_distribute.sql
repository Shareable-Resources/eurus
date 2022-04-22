-- Deploy eurus-backend-db-report:2021_08_19_report_daily_distribute to pg

BEGIN;

-- XXX Add DDLs here.
CREATE OR REPLACE FUNCTION public.report_daily_distribute(var_selected_date date, var_distribute_type integer, var_asset_name character varying, var_chain integer)
 RETURNS TABLE(asset_name character varying, selected_date date, total_amount numeric, chain_id integer)
 LANGUAGE plpgsql
AS $function$
begin
return query	
select
	var_asset_name as asset_name ,
	var_selected_date::date as selected_date,
	COALESCE(sum(amount), 0) as total_amount,
	var_chain as chain_id
from
	distributed_tokens dt
where
	dt.distributed_type = var_distribute_type and 
	dt.created_date::date=var_selected_date and
	dt.asset_name = var_asset_name and 
	dt."chain" = var_chain;

end;

$function$
;



COMMIT;
