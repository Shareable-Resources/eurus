-- Verify eurus-backend-db-user:20210811_create_merchant_refund_request on pg

BEGIN;

-- XXX Add verifications here.
SELECT * FROM merchant_refund_requests LIMIT 1;

ROLLBACK;
