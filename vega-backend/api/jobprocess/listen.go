package jobprocess

import (
	"context"
	"fmt"
	"log"
	db "vega/api/internal/database"

	"github.com/jackc/pgx/v5"
)

const jobChannel = "job_created"

// Listen subscribes to job-created notifications on a dedicated connection and
// processes a job each time one arrives. It blocks until ctx is cancelled.
func Listen(ctx context.Context, conn *pgx.Conn, q *db.Queries) error {
	if _, err := conn.Exec(ctx, "LISTEN "+jobChannel); err != nil {
		return fmt.Errorf("failed to listen on %q: %w", jobChannel, err)
	}
	log.Printf("[jobprocess] listening on %q", jobChannel)

	for {
		notification, err := conn.WaitForNotification(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil // context cancelled — clean shutdown
			}
			return fmt.Errorf("error waiting for notification: %w", err)
		}
		log.Printf("[jobprocess] notified: job %s created", notification.Payload)

		if _, err := Process(ctx, q, defaultLockDuration); err != nil {
			log.Printf("[jobprocess] processing error: %v", err)
		}
	}
}
