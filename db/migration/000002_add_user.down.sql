
drop index if exists "idx_owner";
alter table if exists account drop constraint if exists "account_owner_fkey";
DROP TABLE if exists "user";
