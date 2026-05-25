CREATE TYPE job_state AS ENUM (
    'pending',
    'processing',
    'waiting',
    'processed',
    'errored'
);

CREATE TABLE job_queue (
    id           UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
    content      TEXT      NOT NULL,
    timestamp    TIMESTAMPTZ NOT NULL DEFAULT now(),
    worker_id    UUID      REFERENCES workers(id) ON DELETE SET NULL,
    state        job_state NOT NULL DEFAULT 'pending',
    locked_until TIMESTAMPTZ,
    payload      JSONB
);

CREATE INDEX job_queue_state_idx ON job_queue (state)
WHERE state IN ('pending', 'processing', 'waiting');

---- create above / drop below ----

DROP INDEX job_queue_state_idx;
DROP TABLE job_queue;
DROP TYPE job_state;
