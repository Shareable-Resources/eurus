-- Deploy config:20211124_create_keep_alive_config to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS keep_alive_config (
	id int PRIMARY KEY,
	last_modified_time bigint NOT NULL
);


INSERT INTO keep_alive_config (id, last_modified_time) VALUES (1, 0) ON CONFLICT DO NOTHING;

COMMIT;
