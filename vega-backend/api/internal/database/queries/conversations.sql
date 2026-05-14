-- name: GetConversation :one
SELECT * FROM conversations
WHERE id = $1;

-- name: ListConversations :many
SELECT * FROM conversations
ORDER BY created_at DESC;
