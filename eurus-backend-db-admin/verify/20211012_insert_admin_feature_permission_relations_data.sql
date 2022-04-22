-- Verify eurus-backend-db-admin:20211012_insert_admin_feature_permission_relations_data on pg

BEGIN;

-- XXX Add verifications here.
select 1/count(*) from admin_feature_permission_relations where feature_id >= 1 AND feature_id <= 37;

ROLLBACK;
