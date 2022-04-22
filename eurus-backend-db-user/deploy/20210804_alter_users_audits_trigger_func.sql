-- Deploy eurus-backend-db-user:20210804_alter_users_audits_trigger_func to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE user_audits RENAME COLUMN kyc_status to kyc_level;


CREATE OR REPLACE FUNCTION public.update_user()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$

BEGIN
	INSERT INTO user_audits(user_id,login_address,wallet_address,mainnet_wallet_address,owner_wallet_address, status, email, kyc_level, last_modified_date, created_date, is_metamask_addr)
	VALUES(OLD.id,OLD.login_address,OLD.wallet_address, OLD.mainnet_wallet_address, OLD.owner_wallet_address, OLD.status, OLD.email, OLD.kyc_level,now(),OLD.created_date, OLD.is_metamask_addr);
	RETURN NEW;
END;
$function$
;


COMMIT;
