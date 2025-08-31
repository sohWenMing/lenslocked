-- +goose Up
-- +goose StatementBegin
CREATE TABLE forgot_password_tokens(
    id SERIAL PRIMARY KEY,
    token UUID,
    expires_on TIMESTAMPTZ NOT NULL DEFAULT(now() + INTERVAL '15 minutes')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE forgot_password_tokens;
-- +goose StatementEnd
