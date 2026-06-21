-- name: CreateWorker :one
INSERT INTO workers (id)
VALUES ($1)
RETURNING *;

-- name: GetWorker :one
SELECT * FROM workers
WHERE id = $1;

-- name: ListWorkers :many
SELECT * FROM workers;
