-- Verify eurus-backend-db-user:users_add_wallet_addr_idx on pg

BEGIN;

select
    1 / COUNT(*)
from
    pg_class t,
    pg_class i,
    pg_index ix,
    pg_attribute a
where
    t.oid = ix.indrelid
    and i.oid = ix.indexrelid
    and a.attrelid = t.oid
    and a.attnum = ANY(ix.indkey)
    and t.relkind = 'r'
    and t.relname like 'users'
    and i.relname = 'users_wallet_addr_idx'
group by
    t.relname,
    i.relname;

ROLLBACK;
