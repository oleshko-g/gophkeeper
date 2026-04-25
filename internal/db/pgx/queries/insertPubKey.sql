-- name: InsertPubKey :one
INSERT INTO public_keys (id, key)
VALUES ($1, $2)
RETURNING *;
