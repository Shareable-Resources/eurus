-- Verify eurus-backend-db-user:transfer_trans_add_gas_price on pg

BEGIN;

-- XXX Add verifications here.
SELECT gas_price from transfer_transactions limit 1;

ROLLBACK;
