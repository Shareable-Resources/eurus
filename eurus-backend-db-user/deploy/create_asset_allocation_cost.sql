CREATE SEQUENCE IF NOT EXISTS public.asset_allocation_cost_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;
	

CREATE TABLE IF NOT EXISTS public.asset_allocation_costs (
	id int4 NOT NULL DEFAULT nextval('asset_allocation_cost_id_seq'::regclass),
	trans_hash varchar(255) NULL,
	allocation_type varchar(10) NULL,
	gas_used numeric(78) NULL,
	gas_price numeric(78) NULL,
	created_date timestamptz NULL,
	last_modified_date timestamptz NULL,
	CONSTRAINT asset_allocation_cost_pkey PRIMARY KEY (id)
);
