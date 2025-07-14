-- +goose Up
CREATE TABLE feeds (
	name TEXT NOT NULL,
	url TEXT UNIQUE NOT NULL, 
	user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;