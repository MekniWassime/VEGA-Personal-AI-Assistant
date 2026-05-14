-- name: GetWorker :one
SELECT * FROM workers
WHERE id = $1;

-- name: ListWorkers :many
SELECT * FROM workers;
