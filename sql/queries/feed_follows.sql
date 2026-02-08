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
) SELECT
    fo.id,
    fo.created_at,
    fo.updated_at,
    fo.user_id,
    fo.feed_id,
    f.name AS feed_name,
    u.name AS user_name
FROM inserted_feed_follow as fo
JOIN feeds AS f ON fo.feed_id = f.id
JOIN users AS u ON fo.user_id = u.id;
