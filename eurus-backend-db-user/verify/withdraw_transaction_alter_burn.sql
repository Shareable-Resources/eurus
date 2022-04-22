-- Verify eurus-backend-db-user:withdraw_transaction_alter_burn on pg

BEGIN;

-- XXX Add verifications here.
SELECT burn_trans_hash,"status",burn_date,request_date FROM withdraw_transactions;
  
ROLLBACK;
