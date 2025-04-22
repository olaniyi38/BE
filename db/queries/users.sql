-- name: CreateUser :one
insert into users (username, password, email, full_name)
values ($1,$2,$3,$4)
returning *;

-- name: GetUser :one
select * from users
where username=$1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $1,
    password_updated_at = now()
WHERE username=$2;

-- name: UpdateUserData :one
UPDATE users
SET username = $1,
    full_name = $2,
    email = $3
WHERE username = $1
RETURNING *;

