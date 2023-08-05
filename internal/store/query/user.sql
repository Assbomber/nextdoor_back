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

-- name: UpdateBasicUserDetails :exec
UPDATE users
SET name = $1, gender = $2, birth_date=$3
WHERE ID = sqlc.arg('UserID');


-- name: CreateUserLocation :one
INSERT INTO users_locations (
    user_id,
    location, 
    active -- Defines if this location is currently active or not
) VALUES (
    $1, ST_SetSRID(ST_MakePoint(sqlc.arg(longitude)::double precision, sqlc.arg(latitude)::double precision),4326), $2
) RETURNING *;

-- name: InactiveUserLocation :exec
UPDATE users_locations
SET active = false
WHERE active = true AND user_id = $1;