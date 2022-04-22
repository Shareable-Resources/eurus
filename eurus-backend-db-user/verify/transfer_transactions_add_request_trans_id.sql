-- Verify eurus-backend-db-user:transfer_transactions_add_request_trans_id on pg

BEGIN;

-- XXX Add verifications here.
Select request_trans_id, is_send from transfer_transactions limit 1;
ROLLBACK;
