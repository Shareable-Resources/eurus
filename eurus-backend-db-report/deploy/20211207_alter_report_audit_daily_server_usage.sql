-- Deploy eurus-backend-db-report:20211207_alter_report_audit_daily_server_usage to pg

BEGIN;
drop function if exists public.report_audit_daily_server_usage;

-- XXX Add DDLs here.
CREATE OR REPLACE FUNCTION public.report_audit_daily_server_usage(var_selected_date date)
 RETURNS TABLE(p_wallet_type character varying, p_used numeric, p_asset_name character varying, p_mark_date date)
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
((yesterday_balance - today_balance) / power(10, a."decimal"))::decimal(78, 18) as used, 
w1.asset_name,
w1.mark_date
from 
(select 
wallet_type,
sum(balance) as today_balance, 
asset_name , 
mark_date
from 
wallet_balances 
where 
mark_date  = var_selected_date
and (wallet_type  < 90 or wallet_type in (92, 93))
and asset_name = 'EUN'
group by 
wallet_type,  asset_name ,  mark_date
) as w1 
inner join (
select 
wallet_type,
sum(balance) as yesterday_balance, 
asset_name , 
mark_date 
from 
wallet_balances w2 
where 
mark_date  = var_selected_date::date - 1
and (wallet_type  < 90 or wallet_type in (92, 93))
and asset_name = 'EUN'
group by 
wallet_type, w2.asset_name , w2.mark_date) w2 
on w1.wallet_type = w2.wallet_type
inner join assets as a 
ON w1.asset_name = a.asset_name;

end;

$function$
;



COMMIT;
