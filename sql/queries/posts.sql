-- name: CreatePost :one
INSERT INTO posts (
    id,
    created_at,
    updated_at,
    title,
    url,
    description,
    published_at,
    feed_id
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.* FROM posts
INNER JOIN feed_follows
ON feed_follows.feed_id = posts.feed_id
WHERE user_id = $1
ORDER BY posts.created_at DESC
LIMIT $2;