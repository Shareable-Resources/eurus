-- Deploy eurus-backend-db-user:20220124_distributed_token_errors to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE public.distributed_token_errors (
	id bigint NOT NULL PRIMARY KEY,
	asset_name varchar(100) NOT NULL,
	amount numeric(78) NOT NULL,
	"chain" int2 NULL,
	distributed_type int4 NOT NULL,
	user_id int8 NULL,
	tx_hash varchar(255) NULL,
	from_address varchar(255) NOT NULL,
	to_address varchar(255) NOT NULL,
	gas_price numeric(78) NULL,
	gas_used int8 NULL,
	gas_fee numeric(78) NULL,
	created_date timestamptz NOT NULL,
	last_modified_date timestamptz NOT NULL,
	trigger_type int4 NOT NULL DEFAULT 0
);

COMMIT;
