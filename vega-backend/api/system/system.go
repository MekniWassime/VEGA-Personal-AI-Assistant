package system

import "vega/api/skills"

const systemBasePrompt = `
	You are an AI agent, you will be given instructions by the user which you would need to analyse, understand, complete and verify the result

	If the task has been completed, respond only with "TASK_COMPLETE", if you intend this message to be the last, append "TASK_COMPLETE" to the strict end of the message
	Before declaring the task completed, always verify that your work aligns with the task and explain why it accomplishes it, if a problem is detected try to fix it
	tool execusions on their own do not count as a task completed, you should write a summery completion message before declaring the task is completed
`

var SystemPrompt = skills.BuildSystemPrompt(systemBasePrompt)
