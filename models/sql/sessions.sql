CREATE TABLE sessions (
    id SERIAL PRIMARY KEY NOT NULL,
    user_id INT,
    token_hash TEXT UNIQUE NOT NULL,
    expires_on TIMESTAMPTZ NOT NULL DEFAULT(now() + INTERVAL '15 minutes'),
    is_expired BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)
