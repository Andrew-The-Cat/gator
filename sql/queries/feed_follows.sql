-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds
ON feeds.id = inserted_feed_follow.feed_id
INNER JOIN users
ON users.id = inserted_feed_follow.user_id;

-- name: FeedFollowsReset :exec
DELETE FROM feed_follows *;

-- name: GetFeedFollowsForUser :many
SELECT
feeds.name AS feed_name,
feeds.url AS feed_url
FROM feed_follows
INNER JOIN users
ON users.id = feed_follows.user_id
INNER JOIN feeds
ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollowForUser :exec
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1
AND feed_follows.feed_id = 
    (SELECT id FROM feeds
    WHERE url = $2
    );