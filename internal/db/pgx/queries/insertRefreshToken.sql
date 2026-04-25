-- name: InsertRefreshToken :exec
INSERT INTO
  depositor_refresh_tokens (
    id,
    token,
    depositor_pub_key_id,
    issued_at,
    revoked_at
  )
VALUES
  ($1, $2, $3, $4, $5);
