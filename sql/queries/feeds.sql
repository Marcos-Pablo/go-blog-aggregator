-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT 
    f.id,
    f.created_at,
    f.updated_at,
    f.name,
    f.url,
    f.user_id,
    u.name AS user_name
FROM feeds AS f
JOIN users AS u ON f.user_id = u.id;

-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET updated_at = (NOW() AT TIME ZONE 'UTC'), 
last_fetched_at = (NOW() AT TIME ZONE 'UTC')
WHERE id = $1
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT * from feeds
ORDER BY last_fetched_at ASC
LIMIT 1;
