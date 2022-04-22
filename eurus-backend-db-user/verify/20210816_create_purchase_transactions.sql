-- Verify eurus-backend-db-user:20210816_create_purchase_transactions on pg

BEGIN;

-- XXX Add verifications here.
select * from purchase_transactions limit 1;

ROLLBACK;
