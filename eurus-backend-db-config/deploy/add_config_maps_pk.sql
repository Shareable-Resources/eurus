-- Deploy config:add_config_maps_pk to pg
BEGIN;

ALTER TABLE "public".config_maps 
DROP CONSTRAINT config_maps_pkey,
ADD CONSTRAINT config_maps_pkey PRIMARY KEY(id, is_service, key);

COMMIT;
