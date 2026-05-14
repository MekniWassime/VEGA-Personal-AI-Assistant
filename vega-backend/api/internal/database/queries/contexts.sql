-- name: GetContext :one
SELECT * FROM contexts
WHERE id = $1;

-- name: GetContextByConversation :one
SELECT * FROM contexts
WHERE conversation_id = $1
ORDER BY timestamp DESC
LIMIT 1;
