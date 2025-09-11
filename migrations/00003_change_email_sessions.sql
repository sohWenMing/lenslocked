-- +goose Up
-- +goose StatementBegin
CREATE TABLE forgot_password_tokens(
    id SERIAL PRIMARY KEY,
    user_id int,
    token UUID,
    expires_on TIMESTAMPTZ NOT NULL DEFAULT(now() + INTERVAL '15 minutes'),
    CONSTRAINT fk_user FOREIGN KEY (user_id) references users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE forgot_password_tokens;
-- +goose StatementEnd
