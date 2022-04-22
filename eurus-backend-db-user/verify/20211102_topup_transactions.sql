-- Verify eurus-backend-db-user:20211102_topup_transactions on pg

BEGIN;

select * from top_up_transactions limit 1;

ROLLBACK;
