-- Revert config:add_config_maps_pk from pg

BEGIN;

ALTER TABLE "public".config_maps DROP CONSTRAINT config_maps_pkey,
ADD CONSTRAINT config_maps_pkey PRIMARY KEY(id, key);

COMMIT;
