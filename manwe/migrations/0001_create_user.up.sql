CREATE TABLE user (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL,
    age BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
