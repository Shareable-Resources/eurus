-- Deploy eurus-backend-db-user:transaction_indices_add_assetNameUserIDIndex to pg

BEGIN;

CREATE INDEX UserIdAssetNameIdx on transaction_indices(user_id,asset_name);

COMMIT;
