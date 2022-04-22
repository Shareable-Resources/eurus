-- Revert eurus-backend-db-user:create_user_audit from pg

BEGIN;

-- XXX Add DDLs here.
DROP TABLE user_audits;

COMMIT;
