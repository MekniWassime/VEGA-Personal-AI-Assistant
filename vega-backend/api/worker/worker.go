package worker

import (
	"context"
	"errors"
	"fmt"
	"os"
	db "vega/api/internal/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var workerID pgtype.UUID

func ID() pgtype.UUID {
	if !workerID.Valid {
		panic("worker.Init has not been called")
	}
	return workerID
}

func Init(ctx context.Context, q *db.Queries) error {
	raw := os.Getenv("VEGA_CORE_ID")
	if raw == "" {
		raw = uuid.New().String()
	}

	if err := workerID.Scan(raw); err != nil {
		return fmt.Errorf("invalid VEGA_CORE_ID %q: %w", raw, err)
	}

	_, err := q.GetWorker(ctx, workerID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		if _, err := q.CreateWorker(ctx, workerID); err != nil {
			return err
		}
	}

	return nil
}

