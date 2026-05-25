package jobqueue

import (
	"context"
	"errors"
	"fmt"
	db "vega/api/internal/database"

	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrInvalidTransition = errors.New("invalid state transition")
	ErrNoJob             = errors.New("no job available")
)

func Enqueue(ctx context.Context, q *db.Queries, content string, workerID pgtype.UUID) (db.JobQueue, error) {
	return q.Enqueue(ctx, db.EnqueueParams{
		Content:  content,
		WorkerID: workerID,
	})
}

func Dequeue(ctx context.Context, q *db.Queries, lockDuration string) (db.JobQueue, error) {
	interval := pgtype.Interval{}
	if err := interval.Scan(lockDuration); err != nil {
		return db.JobQueue{}, fmt.Errorf("invalid lock duration %q: %w", lockDuration, err)
	}

	job, err := q.Dequeue(ctx, interval)
	if err != nil {
		return db.JobQueue{}, ErrNoJob
	}
	return job, nil
}

func SetDone(ctx context.Context, q *db.Queries, job db.JobQueue) (db.JobQueue, error) {
	if job.State != db.JobStateProcessing {
		return db.JobQueue{}, fmt.Errorf("%w: %s -> processed", ErrInvalidTransition, job.State)
	}
	return q.SetDone(ctx, job.ID)
}

func SetErrored(ctx context.Context, q *db.Queries, job db.JobQueue) (db.JobQueue, error) {
	if job.State != db.JobStateProcessing {
		return db.JobQueue{}, fmt.Errorf("%w: %s -> errored", ErrInvalidTransition, job.State)
	}
	return q.SetErrored(ctx, job.ID)
}

func SetWaiting(ctx context.Context, q *db.Queries, job db.JobQueue) (db.JobQueue, error) {
	if job.State != db.JobStateProcessing {
		return db.JobQueue{}, fmt.Errorf("%w: %s -> waiting", ErrInvalidTransition, job.State)
	}
	return q.SetWaiting(ctx, job.ID)
}

func Claim(ctx context.Context, q *db.Queries, job db.JobQueue, lockDuration string) (db.JobQueue, error) {
	if job.State != db.JobStateWaiting {
		return db.JobQueue{}, fmt.Errorf("%w: %s -> processing", ErrInvalidTransition, job.State)
	}

	interval := pgtype.Interval{}
	if err := interval.Scan(lockDuration); err != nil {
		return db.JobQueue{}, fmt.Errorf("invalid lock duration %q: %w", lockDuration, err)
	}

	return q.ClaimJob(ctx, db.ClaimJobParams{ID: job.ID, Column2: interval})
}
