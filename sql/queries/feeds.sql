-- name: AddFeed :one
INSERT INTO feeds (name, url, user_id) 
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.name, feeds.url, users.name as user_name FROM feeds
FULL JOIN users
ON users.id = feeds.user_id;