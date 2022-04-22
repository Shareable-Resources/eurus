-- Verify eurus-backend-db-auth:create_login_rquest_token_map on pg

BEGIN;

SET search_path to public;
-- XXX Add verifications here.
SELECT * FROM login_request_token_maps;

ROLLBACK;
