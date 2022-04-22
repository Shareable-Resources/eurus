-- Deploy eurus-backend-db-user:20220124_alter_distributed_tokens to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE distributed_tokens ADD status BIGINT NOT NULL DEFAULT(20);
ALTER TABLE distributed_tokens ALTER COLUMN status SET DEFAULT -1 * currval('distributed_tokens_id_seq'::regclass);

DROP INDEX IF EXISTS distributed_token_idx1;
CREATE UNIQUE INDEX IF NOT EXISTS distributed_tokens_idx1 ON distributed_tokens (user_id, distributed_type, status);

COMMIT;
