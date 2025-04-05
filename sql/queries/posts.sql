-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
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

-- name: GetPosts :many
SELECT posts.* FROM posts
INNER JOIN feeds ON feeds.id = posts.feed_id
WHERE feeds.user_id = $2
ORDER BY posts.updated_at ASC LIMIT $1;
