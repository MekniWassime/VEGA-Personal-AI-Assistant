package jobqueue

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	db "vega/api/internal/database"
	"vega/api/worker"

	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrInvalidTransition = errors.New("invalid state transition")
	ErrNoJob             = errors.New("no job available")
)

// toInterval converts a time.Duration into a pgtype.Interval.
func toInterval(d time.Duration) pgtype.Interval {
	return pgtype.Interval{Microseconds: d.Microseconds(), Valid: true}
}

func Enqueue(ctx context.Context, q *db.Queries, payload []byte, workerID pgtype.UUID) (db.JobQueue, error) {
	job, err := q.Enqueue(ctx, db.EnqueueParams{Payload: payload, WorkerID: workerID})
	if err != nil {
		return db.JobQueue{}, err
	}
	log.Printf("[jobqueue] enqueued | id: %s | worker: %s | payload: %s | timestamp: %s", job.ID, job.WorkerID, job.Payload, job.Timestamp.Time)
	return job, nil
}

func Dequeue(ctx context.Context, q *db.Queries, lockDuration time.Duration) (db.JobQueue, error) {
	job, err := q.Dequeue(ctx, toInterval(lockDuration))
	if err != nil {
		return db.JobQueue{}, ErrNoJob
	}
	log.Printf("[jobqueue] dequeued | id: %s | worker: %s | payload: %s | locked_until: %s", job.ID, job.WorkerID, job.Payload, job.LockedUntil.Time)
	return job, nil
}

func SetDone(ctx context.Context, q *db.Queries, job db.JobQueue) (db.JobQueue, error) {
	if job.State != db.JobStateProcessing {
		return db.JobQueue{}, fmt.Errorf("%w: %s -> processed", ErrInvalidTransition, job.State)
	}
	updated, err := q.SetDone(ctx, job.ID)
	if err != nil {
		return db.JobQueue{}, err
	}
	log.Printf("[jobqueue] id: %s | %s -> processed", job.ID, job.State)
	return updated, nil
}

func SetErrored(ctx context.Context, q *db.Queries, job db.JobQueue) (db.JobQueue, error) {
	if job.State != db.JobStateProcessing {
		return db.JobQueue{}, fmt.Errorf("%w: %s -> errored", ErrInvalidTransition, job.State)
	}
	updated, err := q.SetErrored(ctx, job.ID)
	if err != nil {
		return db.JobQueue{}, err
	}
	log.Printf("[jobqueue] id: %s | %s -> errored", job.ID, job.State)
	return updated, nil
}

func SetWaiting(ctx context.Context, q *db.Queries, job db.JobQueue) (db.JobQueue, error) {
	if job.State != db.JobStateProcessing {
		return db.JobQueue{}, fmt.Errorf("%w: %s -> waiting", ErrInvalidTransition, job.State)
	}
	updated, err := q.SetWaiting(ctx, job.ID)
	if err != nil {
		return db.JobQueue{}, err
	}
	log.Printf("[jobqueue] id: %s | %s -> waiting", job.ID, job.State)
	return updated, nil
}

func Claim(ctx context.Context, q *db.Queries, job db.JobQueue, lockDuration time.Duration) (db.JobQueue, error) {
	if job.State != db.JobStateWaiting {
		return db.JobQueue{}, fmt.Errorf("%w: %s -> processing", ErrInvalidTransition, job.State)
	}

	updated, err := q.ClaimJob(ctx, db.ClaimJobParams{ID: job.ID, Column2: toInterval(lockDuration), WorkerID: worker.ID()})
	if err != nil {
		return db.JobQueue{}, err
	}
	log.Printf("[jobqueue] claimed | id: %s | worker: %s | locked_until: %s", job.ID, worker.ID(), updated.LockedUntil.Time)
	return updated, nil
}
