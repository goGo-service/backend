CREATE TABLE user_tokens
(
    id            SERIAL PRIMARY KEY,
    user_id       int       NOT NULL,
    session_id    UUID      NOT NULL,
    refresh_token TEXT      NOT NULL,
    expire_at     TIMESTAMP NOT NULL,
    created_at    TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id)
)
