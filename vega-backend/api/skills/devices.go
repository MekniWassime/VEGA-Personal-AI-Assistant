package skills

var promptInfo = SkillPromptInfo{
	Name: "list_devices",
	Summary: `
		list client devices identifiers,
		you would need these identifiers to execute certain skills on client devices
		`,
	Content: `
		To execute your tasks, you will need to interact with client devices,
		this skill will provide you with the list of device identifiers and a short description of each.
	`,
	Example: `{"name": "list_devices"}`,
}

type DevicesSkill struct{}

func (d *DevicesSkill) Info() *SkillPromptInfo {
	return &promptInfo
}

type paramsSchema struct{}

func (d *DevicesSkill) Run(input string) (string, error) {
	var devices = `
		[
		{id: "mobile_device_1", "description": "my personal samsung mobile phone" },
		{id: "macbook_work_1", "description": "my work macbook pro that I also use for personal stuff"}
		]
	`
	return devices, nil
}
