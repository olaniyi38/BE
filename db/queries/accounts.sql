-- name: GetAccount :one
    select * from accounts
    where id = $1 LIMIT 1;


-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListAccounts :many
select * from accounts
where name = $1
order by id
limit $2
offset $3;

-- name: CreateAccount :one
insert into accounts (name, balance, currency)
values ($1, $2, $3)
returning *;

-- name: UpdateAccount :one
update accounts
set balance = balance + sqlc.arg(amount)
where id = sqlc.arg(id)
returning *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
delete from accounts
where id = $1;
