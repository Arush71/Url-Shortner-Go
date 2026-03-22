-- +goose Up
CREATE TABLE URLS(
 id BIGSERIAL PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    counter BIGINT DEFAULT 0 NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE URLS;