-- Verify eurus-backend-db-user:transfer_trans_add_gas_used on pg

BEGIN;

-- XXX Add verifications here.
SELECT trans_gas_used, user_gas_used FROM transfer_transactions limit 1;

ROLLBACK;
