-- name: CreateUser :one
INSERT INTO users (
    id, created_at, updated_at, email, hashed_password
) VALUES (
    gen_random_uuid(), now(), now(), $1, $2
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;


-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET email = $2, hashed_password = $3, updated_at=now()
WHERE id = $1
RETURNING *;

-- name: UpdateUserChirpyRed :exec
UPDATE users
SET is_chirpy_red = True
WHERE id = $1;
