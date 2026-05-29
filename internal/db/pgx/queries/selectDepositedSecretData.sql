-- name: SelectDepositedSecretData :one
SELECT
  depositor_pub_key_id,
  encrypted_data
FROM
  deposited_secrets
WHERE
  id = $1;
