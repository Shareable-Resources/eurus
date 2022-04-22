-- Verify eurus-backend-db-user:withdraw_requests on pg

BEGIN;

SELECT * FROM withdraw_requests LIMIT 1;

ROLLBACK;
