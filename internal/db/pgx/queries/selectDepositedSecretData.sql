-- name: SelectDepositedSecretData :one
SELECT encrypted_data FROM deposited_secrets WHERE id = $1;
