alter table if exists "accounts" drop constraint if exists "unique_account";
alter table if exists "accounts" drop constraint if exists "fk_accounts_users";
DROP TABLE users;