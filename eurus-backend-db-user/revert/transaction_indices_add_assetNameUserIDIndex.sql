-- Revert eurus-backend-db-user:transaction_indices_add_assetNameUserIDIndex from pg

BEGIN;

EXPLAIN SELECT * FROM transaction_indices;
DROP INDEX UserIdAssetNameIdx;

COMMIT;
