-- Deploy eurus-backend-db-user:20211109_distributed_token_add_index to pg

BEGIN;

-- XXX Add DDLs here.
CREATE UNIQUE INDEX IF NOT EXISTS distributed_token_idx1 ON distributed_tokens (user_id, distributed_type);

COMMIT;
