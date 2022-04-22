-- Revert eurus-backend-db-user:create_verification_table from pg

BEGIN;

DROP TABLE verifications;

COMMIT;
