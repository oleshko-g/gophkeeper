-- name: DeleteSecretByID :exec
DELETE FROM deposited_secrets
WHERE
  depositor_pub_key_id = $1
  AND id = $2;
