-- name: CreateConversation :one
INSERT INTO conversations (type)
VALUES ($1)
RETURNING *;

-- name: GetConversation :one
SELECT * FROM conversations
WHERE id = $1;

-- name: ListConversations :many
SELECT * FROM conversations
ORDER BY created_at DESC;
