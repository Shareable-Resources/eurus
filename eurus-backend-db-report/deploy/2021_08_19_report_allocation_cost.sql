-- Deploy eurus-backend-db-report:2021_08_19_report_allocation_cost to pg

BEGIN;

-- XXX Add DDLs here.
CREATE OR REPLACE FUNCTION public.report_allocaton_cost(var_selected_date date, var_allocation_type character varying)
 RETURNS TABLE(selected_date date, total_gas_fee numeric, allocation_type character varying)
 LANGUAGE plpgsql
AS $function$
begin
return query	

	select
		var_selected_date as selected_date,
		coalesce (SUM(gas_used * gas_price ),
		0) as total_gas_fee,
		var_allocation_type as allocation_type
from
		asset_allocation_costs aac
where
		created_date::date = var_selected_date
	and aac.allocation_type = var_allocation_type;
	
end;

$function$
;


COMMIT;
