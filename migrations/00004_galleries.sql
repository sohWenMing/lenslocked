-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE galleries (
    id SERIAL PRIMARY KEY,
    user_id int,
    title TEXT,
    CONSTRAINT fk_user FOREIGN KEY (user_id) references users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE galleries;
-- +goose StatementEnd
