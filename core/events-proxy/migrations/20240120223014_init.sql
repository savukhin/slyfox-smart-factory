-- +goose Up
CREATE TABLE users (
    id BIGINT,
    username VARCHAR(256),
    aes_key VARCHAR(256),
    totp_key VARCHAR(256),
    active BOOLEAN
);

-- +goose Down
DROP TABLE users;
