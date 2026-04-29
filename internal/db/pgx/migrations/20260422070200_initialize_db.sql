-- +goose Up
CREATE TABLE depositor_pub_keys (id UUID PRIMARY KEY, pub_key TEXT UNIQUE NOT NULL);


CREATE TABLE depositor_refresh_tokens (
  token UUID PRIMARY KEY,
  depositor_pub_key_id UUID NOT NULL REFERENCES depositor_pub_keys (id),
  revoked_at TIMESTAMP WITH TIME ZONE
);


COMMENT ON TABLE depositor_refresh_tokens IS 'Refresh tokens for depositors';


COMMENT ON COLUMN depositor_refresh_tokens.token IS 'token is a UUID v7 value. It includes the timestamp at which the token was issued';


CREATE TABLE depositor_apps (
  id UUID PRIMARY KEY,
  app_name TEXT NOT NULL,
  depositor_pub_key_id UUID NOT NULL,
  CONSTRAINT app_name_depositor_pub_key_id UNIQUE (app_name, depositor_pub_key_id)
);


CREATE TABLE deposited_data (
  id UUID PRIMARY KEY,
  depositor_pub_key_id UUID NOT NULL REFERENCES depositor_pub_keys (id) ON DELETE CASCADE,
  encrypted_data JSONB NOT NULL
);


CREATE TABLE depositor_app_sessions (
  id UUID PRIMARY KEY,
  depositor_app_id UUID NOT NULL REFERENCES depositor_apps (id) ON DELETE CASCADE
);


COMMENT ON COLUMN depositor_app_sessions.id IS 'token is a UUID v7 value. It includes the timestamp at which was started';


-- +goose Down
DROP TABLE IF EXISTS depositor_app_sessions;


DROP TABLE IF EXISTS deposited_data;


DROP TABLE IF EXISTS depositor_apps;


DROP TABLE IF EXISTS depositor_refresh_tokens;


DROP TABLE IF EXISTS depositor_pub_keys;
