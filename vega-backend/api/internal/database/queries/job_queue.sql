-- name: Enqueue :one
INSERT INTO job_queue (content, worker_id)
VALUES ($1, $2)
RETURNING *;

-- name: Dequeue :one
UPDATE job_queue
SET state = 'processing',
    locked_until = now() + $1::interval
WHERE id = (
    SELECT id FROM job_queue
    WHERE state = 'pending'
    ORDER BY timestamp ASC
    LIMIT 1
    FOR UPDATE SKIP LOCKED
)
RETURNING *;

-- name: SetDone :one
UPDATE job_queue
SET state = 'processed'
WHERE id = $1
RETURNING *;

-- name: SetErrored :one
UPDATE job_queue
SET state = 'errored'
WHERE id = $1
RETURNING *;

-- name: SetWaiting :one
UPDATE job_queue
SET state = 'waiting'
WHERE id = $1
RETURNING *;

-- name: ClaimJob :one
UPDATE job_queue
SET state = 'processing',
    locked_until = now() + $2::interval
WHERE id = $1
RETURNING *;
