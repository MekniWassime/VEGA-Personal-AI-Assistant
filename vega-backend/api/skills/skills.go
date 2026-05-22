package skills

import (
	"encoding/json"
	"fmt"
	"strings"
)

const skillBasePrompt = `
You have access to these skills. These are tool calls that allow you to perform actions and interact with devices and external tools.
To call a skill, wrap the JSON in a <tool_call> tag anywhere in your message: <tool_call>{"name": "skill_name", "arguments": {...}}</tool_call>
The system will notify you of the result. If no notification was received, your tool call format is wrong — try again.
If the response is marked as from "assistant" or "user" then it is not the result of a tool call.
Review the list of available skills carefully:
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

type ParseInput struct {
	Name   string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// ExtractSkillCall tries to find a skill call in common AI response formats.
// Returns the normalised call string and true if a skill call was detected.
func ExtractSkillCall(content string) (string, bool) {
	s := strings.TrimSpace(content)

	// strip <think>...</think> blocks emitted by reasoning models
	if start := strings.Index(s, "<think>"); start != -1 {
		if end := strings.Index(s, "</think>"); end != -1 {
			s = strings.TrimSpace(s[:start] + s[end+len("</think>"):])
		}
	}

	// prefer <tool_call>...</tool_call> tags
	if start := strings.Index(s, "<tool_call>"); start != -1 {
		after := s[start+len("<tool_call>"):]
		if end := strings.Index(after, "</tool_call>"); end != -1 {
			s = strings.TrimSpace(after[:end])
		}
	} else if idx := strings.Index(s, "```"); idx != -1 {
		// fall back to code fence
		after := s[idx:]
		if start := strings.Index(after, "{"); start != -1 {
			after = after[start:]
			if end := strings.Index(after, "```"); end != -1 {
				after = after[:end]
			}
			s = strings.TrimSpace(after)
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

func runSkill(input string) (string, error) {
	skill, parsed, err := ParseAndMatch(input)
	if err != nil {
		return "", err
	}
	return (*skill).Run(string(parsed.Arguments))
}

func ParseAndRun(input string) string {
	result, err := runSkill(input)
	if err != nil {
		return "[TOOL CALL ERROR]\n" + err.Error() + "\n[END TOOL CALL ERROR]\nFix the error and try the tool call again."
	}
	return "[TOOL CALL RESULT]\n" + result + "\n[END TOOL CALL RESULT]\nWith this result, continue your task."
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
