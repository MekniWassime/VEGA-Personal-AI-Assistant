package jobprocess

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	db "vega/api/internal/database"
	"vega/api/jobqueue"
	"vega/api/worker"
)

type ListDevicesHandler struct{}

func (h *ListDevicesHandler) Name() string {
	return "list_devices"
}

type listDevicesArgs struct {
	ConversationID string `json:"conversation_id"`
}

type deviceInfo struct {
	DeviceID    string `json:"device_id"`
	Information string `json:"information"`
}

func (h *ListDevicesHandler) Run(ctx context.Context, q *db.Queries, arguments json.RawMessage) error {
	var args listDevicesArgs
	if err := json.Unmarshal(arguments, &args); err != nil {
		return fmt.Errorf("failed to parse arguments: %w — expected {\"conversation_id\": \"...\"}", err)
	}
	if args.ConversationID == "" {
		return fmt.Errorf("conversation_id is required")
	}

	// mock device list
	devices := []deviceInfo{
		{DeviceID: "mobile_device_1", Information: "my personal samsung mobile phone"},
		{DeviceID: "macbook_work_1", Information: "my work macbook pro that I also use for personal stuff"},
	}
	log.Printf("[list_devices] resolved %d devices", len(devices))

	devicesJSON, err := json.Marshal(devices)
	if err != nil {
		return fmt.Errorf("failed to marshal devices: %w", err)
	}
	message := "Here are the device ids and their information:\n" + string(devicesJSON)

	resultArgs, err := json.Marshal(processConversationArgs{
		ConversationID: args.ConversationID,
		Source:         "system",
		Message:        message,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal result arguments: %w", err)
	}

	resultPayload, err := json.Marshal(jobPayload{
		Name:      "append_process_conversation",
		Arguments: resultArgs,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal result payload: %w", err)
	}

	if _, err := jobqueue.Enqueue(ctx, q, resultPayload, worker.ID()); err != nil {
		return fmt.Errorf("failed to enqueue result job: %w", err)
	}

	return nil
}
