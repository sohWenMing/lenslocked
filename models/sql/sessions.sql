CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE,
    token_hash TEXT UNIQUE NOT NULL,
    expires_on TIMESTAMP,
    is_expired BOOLEAN
)
