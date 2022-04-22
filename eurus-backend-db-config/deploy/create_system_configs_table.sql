-- Deploy config:create_system_configs_table to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS system_configs
(
    id          BIGSERIAL   NOT NULL,
    owner_id    BIGINT      NOT NULL,
    is_service  BOOLEAN     NOT NULL,
    key         VARCHAR     NOT NULL,
    value       VARCHAR     NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (key)
);

COMMIT;
