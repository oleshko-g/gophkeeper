-- ++goose Up
CREATE TABLE depositor_pub_keys (id UUID PRIMARY KEY, pub_key TEXT NOT NULL,);


CREATE TABLE depositor_refresh_tokens (
  id UUID PRIMARY KEY,
  token TEXT NOT NULL,
  depositor_pub_key_id UUID NOT NULL,
  issued_at TIMESTAMP WITH TIME ZONE NOT NULL,
  revoked_at TIMESTAMP WITH TIME ZONE
);


CREATE TABLE depositor_apps (
  id UUID PRIMARY KEY,
  app_name TEXT NOT NULL,
  depositor_pub_key_id UUID NOT NULL,
);


CREATE TABLE deposited_data (
  id UUID PRIMARY KEY,
  depositor_pub_key_id UUID NOT NULL,
  ciphered_data JSONB NOT NULL,
);


CREATE TABLE depositor_app_sessions (
  id UUID PRIMARY KEY,
  depositor_app_id UUID NOT NULL,
  started_at TIMESTAMP WITH TIME ZONE NOT NULL,
  expired_at TIMESTAMP WITH TIME ZONE,
);


-- ++goose Down
DROP TABLE IF EXISTS depositor_app_sessions;


DROP TABLE IF EXISTS deposited_data;


DROP TABLE IF EXISTS depositor_apps;


DROP TABLE IF EXISTS depositor_refresh_tokens;


DROP TABLE IF EXISTS depositor_pub_keys;
