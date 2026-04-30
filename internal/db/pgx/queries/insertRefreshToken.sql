-- name: InsertRefreshToken :exec
INSERT INTO
  depositor_refresh_tokens (token, depositor_pub_key_id, revoked_at)
VALUES
  ($1, $2, $3);
