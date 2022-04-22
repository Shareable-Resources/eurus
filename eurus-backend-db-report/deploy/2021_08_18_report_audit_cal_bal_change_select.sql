-- Deploy eurus-backend-db-report:2021_08_18_report_audit_cal_bal_change_select to pg

BEGIN;

-- XXX Add DDLs here.
CREATE OR REPLACE FUNCTION public.report_audit_cal_bal_change_select(var_selected_date date, var_wallet_type integer)
 RETURNS TABLE(p_asset_name character varying, p_selected_date date, p_created_date timestamp with time zone, p_last_modified_date timestamp with time zone, p_daybe4yesterday_balance numeric, p_yesterday_balance numeric, p_balance_change numeric, p_wallet_type integer)
 LANGUAGE plpgsql
AS $function$
begin
return query	 
select
			an.asset_name as p_asset_name, 
			var_selected_date as p_selected_date,
			now() as p_created_date,
			now() as p_last_modified_date,
		    t2.totalBalance as p_daybe4yesterday_balance,
		    t1.totalBalance as p_yesterday_balance,
		    t1.totalBalance -t2.totalBalance as p_balance_change,
			var_wallet_type as p_wallet_type
from
		assets as an
	-- yesterday
inner join (
	select
				asset_name,
				SUM(balance) as totalBalance
	from
			wallet_balances
	where
			(created_date::date = var_selected_date::date - 1 )
		and wallet_type = var_wallet_type
	group by
				asset_name 
			) 
			as t1
			on
		an.asset_name = t1.asset_name
	-- dayBe4Yesterday
inner join (
	select
			asset_name,
			SUM(balance) as totalBalance
	from
			wallet_balances
	where
			(created_date::date = var_selected_date::date - 2 )
			and wallet_type = var_wallet_type
		group by
				asset_name 
			) 
			as t2
			on
		an.asset_name = t2.asset_name;

end;

$function$
;


COMMIT;
