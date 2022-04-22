-- Verify eurus-backend-db-user:withdraw_add_user_gas_used_column on pg

BEGIN;

-- XXX Add verifications here.
select gas_price, user_gas_used from pending_prewithdraws limit 1;
select gas_price, user_gas_used from withdraw_transactions limit 1;

ROLLBACK;
