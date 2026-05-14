package skills

var promptInfo = SkillPromptInfo{
	Name:    "list_devices",
	Summary: "list client device identifiers and information — you will need these identifiers to execute skills on client devices",
	Content: `
Returns a list of client device identifiers and a short description of each.
Use these identifiers when executing skills that target a specific device.
`,
	Example: `{"name": "list_devices"}`,
}

type DevicesSkill struct{}

func (d *DevicesSkill) Info() *SkillPromptInfo {
	return &promptInfo
}

type devicesSchema struct{}

func (d *DevicesSkill) Run(input string) (string, error) {
	var devices = `[
{"id": "mobile_device_1", "description": "my personal samsung mobile phone"},
{"id": "macbook_work_1", "description": "my work macbook pro that I also use for personal stuff"}
]`
	return devices, nil
}
