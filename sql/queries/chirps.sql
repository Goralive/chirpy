-- name: CreateChirp :one

INSERT INTO chirps (
    id, created_at, updated_at, body, user_id
) VALUES (get_random_uuid(), now(), now(), $1, $2)
RETURNING *;
