package skills

import (
	"context"
	db "vega/api/internal/database"
	"vega/api/jobqueue"
	"vega/api/worker"
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

func (d *DevicesSkill) Run(ctx context.Context, q *db.Queries, input string) (RunResult, error) {
	if _, err := jobqueue.Enqueue(ctx, q, []byte(`{"name":"list_devices"}`), worker.ID()); err != nil {
		return RunResult{}, err
	}
	return RunResult{Suspend: true}, nil
}
