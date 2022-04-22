-- Deploy eurus-backend-db-report:20211101_report_daily_distribute_change_decimal to pg

BEGIN;

CREATE OR REPLACE FUNCTION public.report_daily_distribute(var_selected_date date)
 RETURNS TABLE(asset_name character varying, distributed_type_name character varying, selected_date date, total_amount numeric(78,18), chain_id integer)
 LANGUAGE plpgsql
AS $function$
begin
return query	
select
	dt.asset_name,
	t.name::varchar(255),
	var_selected_date::date as selected_date,
	COALESCE(  (sum(amount) / power(10, a."decimal") ), 0)::numeric(78, 18) as total_amount,
	dt."chain"::integer as chain_id
from
	distributed_tokens dt
inner join assets as a  on dt.asset_name = a.asset_name
inner join distributed_token_types t on  dt.distributed_type = t.id 
where
	dt.created_date::date=var_selected_date
group by
	dt.distributed_type,
	t.name,
	dt.asset_name,
	dt."chain",
	a."decimal"
order by 
	dt.distributed_type, 
	dt.asset_name;
end;

$function$
;

COMMIT;
