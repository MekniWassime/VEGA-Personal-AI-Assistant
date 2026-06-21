package jobprocess

import (
	"context"
	"log"
	"time"
	db "vega/api/internal/database"
)

const defaultLockDuration = 30 * time.Second

// maxDrainedJobs caps how many jobs a single Drain call will process.
const maxDrainedJobs = 20

// Drain processes pending jobs until the queue is empty or maxDrainedJobs is reached.
func Drain(ctx context.Context, q *db.Queries) {
	log.Println("[jobprocess] draining pending jobs...")
	count := 0
	for count < maxDrainedJobs {
		count++
		job, err := Process(ctx, q, defaultLockDuration)
		if err != nil {
			// the failing job is already marked errored, so keep draining the rest
			log.Printf("[jobprocess] job processing error: %v", err)
			continue
		}
		if job == nil {
			break
		}
	}
	log.Printf("[jobprocess] drained %d job(s)", count)
}
