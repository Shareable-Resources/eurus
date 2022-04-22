-- Deploy eurus-backend-db-user:20220114_alter_user_audits to pg

BEGIN;

ALTER TABLE user_audits ADD COLUMN mnemonic TEXT;

DROP TRIGGER IF EXISTS user_trigger ON public.users;

CREATE  TRIGGER user_trigger after
update
    of login_address,
    wallet_address,
    mainnet_wallet_address,
    owner_wallet_address,
    status,
    email,
    kyc_level,
    mnemonic on
    public.users for each row execute function update_user();

COMMIT;
