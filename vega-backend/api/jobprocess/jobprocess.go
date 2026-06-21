package jobprocess

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	db "vega/api/internal/database"
	"vega/api/jobqueue"
)

type JobHandler interface {
	Name() string
	Run(ctx context.Context, q *db.Queries, arguments json.RawMessage) error
}

type jobPayload struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

var registeredHandlers = []JobHandler{
	&ListDevicesHandler{},
	&ProcessConversationHandler{},
}

var handlersByName = buildHandlerMap(registeredHandlers)

func buildHandlerMap(handlers []JobHandler) map[string]JobHandler {
	m := make(map[string]JobHandler, len(handlers))
	for _, h := range handlers {
		m[h.Name()] = h
	}
	return m
}

func Process(ctx context.Context, q *db.Queries, lockDuration time.Duration) (*db.JobQueue, error) {
	job, err := jobqueue.Dequeue(ctx, q, lockDuration)
	if err != nil {
		if errors.Is(err, jobqueue.ErrNoJob) {
			return nil, nil
		}
		return nil, err
	}

	var payload jobPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		_, _ = jobqueue.SetErrored(ctx, q, job)
		return nil, fmt.Errorf("failed to parse job payload: %w", err)
	}

	handler, ok := handlersByName[payload.Name]
	if !ok {
		_, _ = jobqueue.SetErrored(ctx, q, job)
		names := make([]string, 0, len(handlersByName))
		for name := range handlersByName {
			names = append(names, name)
		}
		return nil, fmt.Errorf("no handler matched %q, registered handlers: %s", payload.Name, strings.Join(names, ", "))
	}

	if err := handler.Run(ctx, q, payload.Arguments); err != nil {
		_, _ = jobqueue.SetErrored(ctx, q, job)
		return nil, fmt.Errorf("handler %q failed: %w", payload.Name, err)
	}

	done, err := jobqueue.SetDone(ctx, q, job)
	if err != nil {
		return nil, err
	}
	return &done, nil
}
