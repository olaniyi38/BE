-- name: CreateEntry :one
insert into entries (amount, account_id)
values ($1, $2)
returning *;

-- name: GetEntry :one
select * from entries
where id = $1;

-- name: ListEntries :many
select * from entries
where account_id = $1
order by id
limit $2
offset $3;