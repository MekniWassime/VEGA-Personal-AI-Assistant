-- name: CreateContext :one
INSERT INTO contexts (conversation_id)
VALUES ($1)
RETURNING *;

-- name: GetContext :one
SELECT * FROM contexts
WHERE id = $1;

-- name: GetContextByConversation :one
SELECT * FROM contexts
WHERE conversation_id = $1
ORDER BY timestamp DESC
LIMIT 1;
