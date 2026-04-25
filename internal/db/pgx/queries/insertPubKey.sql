-- name: InsertPubKey :exec
INSERT INTO
  depositor_pub_keys (id, pub_key)
VALUES
  ($1, $2);
