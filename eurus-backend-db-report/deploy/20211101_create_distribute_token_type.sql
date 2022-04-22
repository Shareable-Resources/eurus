-- Deploy eurus-backend-db-report:20211101_create_distribute_token_type to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS distributed_token_types (
	id bigint PRIMARY KEY,
	name text NOT NULL,
	created_date timestamptz NOT NULL,
	last_modified_date timestamptz NOT NULL
);

INSERT INTO distributed_token_types (id, name, created_date, last_modified_date) VALUES (1, 'First deposit wallet', now(), now()) ON CONFLICT DO NOTHING;

COMMIT;
