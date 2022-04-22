-- Deploy eurus-backend-db-user:withdraw_transaction_alter_burn to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE withdraw_transaction ADD COLUMN burn_trans_hash Varchar(255);
ALTER TABLE withdraw_transaction ADD COLUMN "status" smallint;
ALTER TABLE withdraw_transaction ADD COLUMN  burn_date timestamptz;
ALTER TABLE withdraw_transaction ADD COLUMN  request_date timestamptz;

  

COMMIT;
