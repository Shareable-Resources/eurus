-- Verify eurus-backend-db-user:transfer_transactions_add_column on pg

BEGIN;

-- XXX Add verifications here.
Select from_address, confirm_trans_hash, last_modified_date from transfer_transactions limit 1;

ROLLBACK;
