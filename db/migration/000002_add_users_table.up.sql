CREATE TABLE users (
    username VARCHAR(255) PRIMARY KEY,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    password_updated_at TIMESTAMPTZ NOT NULL DEFAULT '1970-01-01 00:00:00',
    full_name VARCHAR(255) NOT NULL
);

ALTER TABLE accounts
ADD CONSTRAINT fk_accounts_users
FOREIGN KEY (name) REFERENCES users (username) ON DELETE CASCADE;

ALTER TABLE accounts
ADD CONSTRAINT unique_account UNIQUE (name, currency);