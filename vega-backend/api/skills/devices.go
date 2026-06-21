package skills

import (
	"context"
	"encoding/json"
	db "vega/api/internal/database"
	"vega/api/jobqueue"
	"vega/api/worker"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var promptInfo = SkillPromptInfo{
	Name:    "list_devices",
	Summary: "list the user's device identifiers and information — you will need these identifiers to execute skills on the user's devices",
	Content: `
Returns a list of the user's device identifiers and a short description of each.
Use these identifiers when executing skills that target a specific device.
`,
	Example: `{"name": "list_devices"}`,
}

type DevicesSkill struct{}

func (d *DevicesSkill) Info() *SkillPromptInfo {
	return &promptInfo
}

type devicesSchema struct{}

func (d *DevicesSkill) Run(ctx context.Context, q *db.Queries, conversationID pgtype.UUID, input string) (RunResult, error) {
	payload, err := json.Marshal(map[string]any{
		"name": "list_devices",
		"arguments": map[string]string{
			"conversation_id": uuid.UUID(conversationID.Bytes).String(),
		},
	})
	if err != nil {
		return RunResult{}, err
	}

	if _, err := jobqueue.Enqueue(ctx, q, payload, worker.ID()); err != nil {
		return RunResult{}, err
	}
	return RunResult{Suspend: true}, nil
}
