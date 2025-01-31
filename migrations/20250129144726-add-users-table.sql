-- +migrate Up

CREATE TABLE users (
    id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX unique_user_email ON users (email) WHERE deleted_at IS NULL;

-- +migrate Down

DROP INDEX IF EXISTS unique_user_email;
DROP TABLE users;