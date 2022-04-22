-- Deploy eurus-backend-db-user:create_users_trigger to pg

BEGIN;

-- XXX Add DDLs here.
CREATE OR REPLACE FUNCTION update_user() RETURNS TRIGGER LANGUAGE PLPGSQL AS $$

BEGIN
	INSERT INTO user_audits(user_id,login_address,wallet_address,mainnet_wallet_address,owner_wallet_address, status, email, kyc_status, last_modified_date, created_date, is_metamask_addr)
	VALUES(OLD.id,OLD.login_address,OLD.wallet_address, OLD.mainnet_wallet_address, OLD.owner_wallet_address, OLD.status, OLD.email, OLD.kyc_status,now(),OLD.created_date, OLD.is_metamask_addr);
	RETURN NEW;
END;
$$;



CREATE TRIGGER user_trigger
AFTER UPDATE OF login_address,wallet_address,mainnet_wallet_address,owner_wallet_address,status,email,kyc_status ON users
FOR EACH ROW 
EXECUTE PROCEDURE update_user();



COMMIT;
