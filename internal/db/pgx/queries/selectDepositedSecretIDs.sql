-- name: SelectDepositedSecretIDs :many
SELECT id FROM deposited_secrets WHERE depositor_pub_key_id = $1;
