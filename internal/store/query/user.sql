-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    password,
    last_login
) VALUES (
    $1,$2,$3,$4
) RETURNING *;


-- name: GetUserByEmailOrUsername :one
SELECT *
FROM users
WHERE email = $1 OR username = $2
LIMIT 1;

-- name: UpdateUserPasswordByEmail :exec
UPDATE users
SET password = $2
WHERE email = $1;

-- name: UpdateUserLoginTimeByEmail :one
UPDATE users
SET last_login = $2
WHERE email = $1
RETURNING *;