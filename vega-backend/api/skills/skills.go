package skills

import (
	"encoding/json"
	"fmt"
	"strings"
)

const skillBasePrompt = `
	You have access to these skills, These are tool calls that allow you to perform actions and interact with devices and external tools,
	You would need to respond only with the tool call with no other text, only the tool call json should be present if you would like to call a tool
	Here are a list of available skills:
`

type SkillPromptInfo struct {
	Name    string `json:"name"`
	Summary string `json:"summary"`
	Content string `json:"content"`
	Example string `json:"example"`
}

type Skill interface {
	Info() *SkillPromptInfo
	Run(input string) (string, error)
}

var registeredSkills = []Skill{
	&SkillDefinitionSkill{},
	&DevicesSkill{},
	&NotifyDeviceSkill{},
}

func buildSkillMap(skills []Skill) map[string]Skill {
	m := make(map[string]Skill, len(skills))
	for _, skill := range skills {
		m[skill.Info().Name] = skill
	}
	return m
}

var skillsByName = buildSkillMap(registeredSkills)

func buildSkillsPrompt() string {
	var sb strings.Builder
	sb.WriteString(skillBasePrompt)
	for _, skill := range registeredSkills {
		info := skill.Info()
		sb.WriteString("\n  - " + info.Name + ": " + info.Summary)
	}
	return sb.String()
}

var SkillsPrompt = buildSkillsPrompt()

func BuildSystemPrompt(basePrompt string) string {
	return basePrompt + "\n" + SkillsPrompt
}

type ParseInput struct {
	Name   string          `json:"name"`
	Params json.RawMessage `json:"params"`
}

// extractSkillCall tries to find a skill call in common AI response formats.
// Returns the normalised call string and true if a skill call was detected.
func ExtractSkillCall(content string) (string, bool) {
	s := strings.TrimSpace(content)

	// strip markdown code fences: ```json ... ``` or ``` ... ```
	for _, fence := range []string{"```json", "```"} {
		if strings.HasPrefix(s, fence) {
			s = strings.TrimPrefix(s, fence)
			s = strings.TrimSuffix(s, "```")
			s = strings.TrimSpace(s)
			break
		}
	}

	// already valid JSON with a name field
	if strings.HasPrefix(s, "{") {
		return s, true
	}

	// plain skill name
	if _, ok := skillsByName[s]; ok {
		return `{"name":"` + s + `"}`, true
	}

	return "", false
}

func ParseAndRun(input string) string {
	skill, parsed, err := ParseAndMatch(input)
	if err != nil {
		return err.Error()
	}
	result, err := (*skill).Run(string(parsed.Params))
	if err != nil {
		return err.Error()
	}
	return result
}

func ParseAndMatch(input string) (*Skill, *ParseInput, error) {
	var parsed ParseInput
	if err := json.Unmarshal([]byte(input), &parsed); err != nil {
		return nil, nil, fmt.Errorf("failed to parse skill call: %w — fix the JSON and try the tool call again", err)
	}

	if skill, ok := skillsByName[parsed.Name]; ok {
		return &skill, &parsed, nil
	}
	names := make([]string, 0, len(skillsByName))
	for name := range skillsByName {
		names = append(names, name)
	}
	return nil, nil, fmt.Errorf("no skill matched %q, registered skills: %s, try the tool call again with a valid skill", parsed.Name, strings.Join(names, ", "))
}
