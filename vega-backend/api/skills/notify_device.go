package skills

import (
	"context"
	"encoding/json"
	"fmt"
	db "vega/api/internal/database"

	"github.com/jackc/pgx/v5/pgtype"
)

var notifyDeviceInfo = SkillPromptInfo{
	Name:    "notify_device",
	Summary: "send a notification message to a specific device by its identifier, Note that the user has no direct access to this conversation, you will often need to send them a notification with the result",
	Content: `
Sends a notification to a client device.
Requires the device identifier (obtainable via list_devices) and a message string.
`,
	Example: `{"name": "notify_device", "arguments": {"device_id": "mobile_device_1", "message": "Hello!"}}`,
}

type NotifyDeviceSkill struct{}

func (n *NotifyDeviceSkill) Info() *SkillPromptInfo {
	return &notifyDeviceInfo
}

type notifyDeviceArgs struct {
	DeviceID string `json:"device_id"`
	Message  string `json:"message"`
}

func (n *NotifyDeviceSkill) Run(ctx context.Context, q *db.Queries, conversationID pgtype.UUID, input string) (RunResult, error) {
	var p notifyDeviceArgs
	if err := json.Unmarshal([]byte(input), &p); err != nil {
		return RunResult{}, fmt.Errorf("failed to parse arguments: %w — expected {\"device_id\": \"...\", \"message\": \"...\"}", err)
	}
	if p.DeviceID == "" {
		return RunResult{}, fmt.Errorf("device_id is required")
	}
	if p.Message == "" {
		return RunResult{}, fmt.Errorf("message is required")
	}

	fmt.Printf("[notify_device] device: %s | message: %s\n", p.DeviceID, p.Message)
	return RunResult{Content: fmt.Sprintf("Notification sent to device %q: %q", p.DeviceID, p.Message)}, nil
}
