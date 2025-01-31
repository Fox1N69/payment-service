-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    amount BIGINT NOT NULL,
    type VARCHAR(50) NOT NULL,
    description TEXT,
    timestamp TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE transactions;
-- +goose StatementEnd
