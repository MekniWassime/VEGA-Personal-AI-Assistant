-- name: CreateMessage :one
INSERT INTO messages (context_id, role, content, worker_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetMessage :one
SELECT * FROM messages
WHERE id = $1;

-- name: ListMessagesByContext :many
SELECT * FROM messages
WHERE context_id = $1
ORDER BY timestamp ASC;

-- name: ListMessagesByConversation :many
SELECT m.* FROM messages m
JOIN contexts c ON m.context_id = c.id
WHERE c.id = (
    SELECT id FROM contexts
    WHERE contexts.conversation_id = $1
    ORDER BY timestamp DESC
    LIMIT 1
)
ORDER BY m.timestamp ASC;
