-- Deploy eurus-backend-db-report:20211203_alter_wallet_balances_created_date to pg

BEGIN;

alter table wallet_balances  alter column created_date set data type timestamptz; 

-- update wallet_balances  set created_date  = created_date  + interval '8 hours' where DATE_PART('hour',created_date) = 0;

ALTER TABLE wallet_balances add mark_date date;

update wallet_balances set mark_date = date(created_date) where mark_date is null;


ALTER TABLE wallet_balances ALTER COLUMN mark_date SET NOT NULL;

COMMIT;
