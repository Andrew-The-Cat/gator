-- +goose Up
CREATE TABLE feed_follows (
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id uuid REFERENCES users(id) NOT NULL,
    feed_id uuid REFERENCES feeds(id) NOT NULL,
    UNIQUE(user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;