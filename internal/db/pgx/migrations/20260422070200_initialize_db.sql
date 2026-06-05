-- +goose Up
CREATE TABLE depositor_pub_keys (id UUID PRIMARY KEY, pub_key TEXT UNIQUE NOT NULL);

COMMENT ON TABLE depositor_pub_keys IS 'depositor_pub_keys are the registered public keys of depositors';
COMMENT ON COLUMN depositor_pub_keys.id IS 'id is a UUID v7 value. It includes the timestamp at which the pub_key was registered';


CREATE TABLE deposited_secrets (
  id UUID PRIMARY KEY,
  depositor_pub_key_id UUID NOT NULL REFERENCES depositor_pub_keys (id) ON DELETE CASCADE,
  encrypted_data BYTEA NOT NULL
);

COMMENT ON COLUMN deposited_secrets.id IS 'id is a UUID v7 value. It includes the timestamp at which the secret was deposited';


-- +goose Down
DROP TABLE IF EXISTS deposited_secrets;


DROP TABLE IF EXISTS depositor_pub_keys;
