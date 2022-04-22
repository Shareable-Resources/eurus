-- Revert eurus-backend-db-user:withdraw_transaction from pg

BEGIN;

-- XXX Add DDLs here.
DROP TABLE withdraw_transaction


COMMIT;
