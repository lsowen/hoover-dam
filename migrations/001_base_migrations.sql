CREATE TABLE hoover_user (
  id SERIAL PRIMARY KEY
  , creation_date TIMESTAMP WITH TIME ZONE NOT NULL
  , email TEXT
  , external_id TEXT
  , friendly_name TEXT
  , "source" TEXT
  , username TEXT NOT NULL UNIQUE
);

CREATE TABLE hoover_credential (
  id SERIAL PRIMARY KEY
  , creation_date TIMESTAMP WITH TIME ZONE NOT NULL
  , access_key_id TEXT NOT NULL UNIQUE
  , secret_access_key TEXT NOT NULL
  , user_id INTEGER NOT NULL REFERENCES hoover_user
);

CREATE TABLE hoover_group (
  id SERIAL PRIMARY KEY
  , creation_date TIMESTAMP WITH TIME ZONE NOT NULL
  , name TEXT NOT NULL UNIQUE
);

CREATE TABLE hoover_user_group (
  id SERIAL NOT NULL
  , user_id INTEGER NOT NULL references hoover_user
  , group_id INTEGER NOT NULL references hoover_group
  , PRIMARY KEY(user_id, group_id)
);


CREATE TABLE hoover_policy (
  id SERIAL PRIMARY KEY
  , creation_date TIMESTAMP WITH TIME ZONE NOT NULL
  , name TEXT NOT NULL UNIQUE
  , policy JSONB NOT NULL
);

CREATE TABLE hoover_user_policy (
  id SERIAL NOT NULL
  , user_id INTEGER NOT NULL references hoover_user
  , policy_id INTEGER NOT NULL references hoover_policy
  , PRIMARY KEY(user_id, policy_id)
);

CREATE TABLE hoover_group_policy (
  id SERIAL NOT NULL
  , group_id INTEGER NOT NULL references hoover_group
  , policy_id INTEGER NOT NULL references hoover_policy
  , PRIMARY KEY(group_id, policy_id)
);

---- create above / drop below ----

DROP TABLE hoover_group_policy;
DROP TABLE hoover_user_policy;
DROP TABLE hoover_policy;
DROP TABLE hoover_user_group;
DROP TABLE hoover_group;
DROP TABLE hoover_credential;
DROP TABLE hoover_user;
