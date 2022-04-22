-- Revert eurus-backend-db-user:20211109_distributed_token_add_index from pg

BEGIN;

-- XXX Add DDLs here.
DROP INDEX distributed_token_idx1;

COMMIT;
