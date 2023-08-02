-- name: CreateUser :one
INSERT INTO users (
    name,
    email,
    password
) VALUES (
    $1,$2,$3
) RETURNING *;


-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: UpdateUserPasswordByEmail :exec
UPDATE users
SET password = $2
WHERE email = $1;