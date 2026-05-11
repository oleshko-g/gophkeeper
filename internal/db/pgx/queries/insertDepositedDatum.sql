-- name: InsertDepositedSecret :one
INSERT INTO
  deposited_secrets (id, depositor_pub_key_id, encrypted_data)
VALUES
  ($1, $2, $3)
RETURNING
  *;
