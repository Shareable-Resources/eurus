-- Deploy eurus-backend-db-report:20210902_report_audit_daily_server_usage to pg

BEGIN;
DROP FUNCTION IF EXISTS report_audit_daily_server_usage(date);

CREATE OR REPLACE FUNCTION public.report_audit_daily_server_usage(var_selected_date date)
 RETURNS TABLE(p_wallet_type character varying, p_used numeric(75), p_asset_name character varying,  p_created_date date)
 LANGUAGE plpgsql
AS $function$
begin
return query	 

select 
case (w1.wallet_type) 
when 1 then 'Deposit'
when 2 then 'Withdraw'
when 3 then 'Withdraw approval'
when 4 then 'User service'
when 6 then 'Centralize user approval'
when 8 then 'Config service'
when 92 then 'Centralize user invoker'
when 93 then 'Centralize user smart contract owner'
else  cast(w1.wallet_type as varchar(10))
end,
yesterday_balance - today_balance as used, 
w1.asset_name,
w1.created_date
from 
(select 
wallet_type,
sum(balance) as today_balance, 
asset_name , 
created_date
from 
wallet_balances 
where 
created_date  = var_selected_date
and (wallet_type  < 90 or wallet_type in (92, 93))
and asset_name = 'EUN'
group by 
wallet_type,  asset_name ,  created_date
) as w1 
inner join (
select 
wallet_type,
sum(balance) as yesterday_balance, 
asset_name , 
created_date 
from 
wallet_balances w2 
where 
created_date  = var_selected_date::date - 1
and (wallet_type  < 90 or wallet_type in (92, 93))
and asset_name = 'EUN'
group by 
wallet_type, w2.asset_name , w2.created_date) w2 
on w1.wallet_type = w2.wallet_type;

end;

$function$
;


COMMIT;
