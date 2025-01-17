CREATE EXTENSION
  IF NOT EXISTS "uuid-ossp";


CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    a_password TEXT NOT NULL,
    first_name TEXT,
    last_name TEXT,

    CONSTRAINT uk_user UNIQUE (username)
);

CREATE INDEX IF NOT EXISTS user_username_idx ON users (username);

CREATE TABLE IF NOT EXISTS tasks (
    id BIGSERIAL PRIMARY KEY,
    identifier UUID NOT NULL DEFAULT uuid_generate_v4()
    title TEXT NOT NULL,
    a_description TEXT NOT NULL,
    a_status TEXT NOT NULL,
    due_date TIMESTAMP NOT NULL,

    CONSTRAINT uk_identifier UNIQUE (identifier)
)

CREATE INDEX IF NOT EXISTS task_identifier_idx ON tasks (identifier);
