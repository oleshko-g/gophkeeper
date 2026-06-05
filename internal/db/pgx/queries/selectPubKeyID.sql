-- name: SelectPubKeyByID :one
SELECT
  *
FROM
  depositor_pub_keys
WHERE
  id = $1;
