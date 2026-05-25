package skills

import (
	"context"
	"encoding/json"
	"fmt"
	db "vega/api/internal/database"
)

var skillDefinitionInfo = SkillPromptInfo{
	Name:    "skill_definition",
	Summary: `Before using any skill, call this first to get its full instructions, Do not try to guess skill parameters, always call this atleast one time per skill needed. Usage: {"name": "skill_definition", "arguments": {"skill_name": "<name>"}}. Example: {"name": "skill_definition", "arguments": {"skill_name": "list_devices"}}`,
	Content: `
Returns a long description of a skill.
Use this to understand how to use a skill before calling it.
`,
	Example: `{"name": "skill_definition", "arguments": {"skill_name": "list_devices"}}`,
}

type SkillDefinitionSkill struct{}

func (s *SkillDefinitionSkill) Info() *SkillPromptInfo {
	return &skillDefinitionInfo
}

type skillDefinitionArgs struct {
	SkillName string `json:"skill_name"`
}

func (s *SkillDefinitionSkill) Run(ctx context.Context, q *db.Queries, input string) (RunResult, error) {
	var p skillDefinitionArgs
	if err := json.Unmarshal([]byte(input), &p); err != nil {
		return RunResult{}, fmt.Errorf("failed to parse arguments: %w — expected {\"skill_name\": \"...\"}", err)
	}
	if p.SkillName == "" {
		return RunResult{}, fmt.Errorf("skill_name is required")
	}
	skill, ok := skillsByName[p.SkillName]
	if !ok {
		return RunResult{}, fmt.Errorf("no skill found with name %q", p.SkillName)
	}
	info := skill.Info()
	return RunResult{Content: fmt.Sprintf(
		"[SKILL DEFINITION]\nName: %s\n\nDescription:\n%s\n\nExample usage:\n%s\n\n[END SKILL DEFINITION]\nThis was the skill definition, not the skill call itself. Now call the skill using the example above as a reference.",
		info.Name, info.Content, info.Example,
	)}, nil
}
