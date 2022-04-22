-- Revert eurus-backend-db-auth:create_login_rquest_token_map from pg

BEGIN;

-- XXX Add DDLs here.
SET search_path to public;
DROP TABLE login_request_token_maps;

COMMIT;
