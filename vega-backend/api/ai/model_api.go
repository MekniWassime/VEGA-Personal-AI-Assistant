package ai

type ModelAPI interface {
	Complete(messages []Message) (*Message, error)
}

func Complete(api ModelAPI, messages []Message) (*Message, error) {
	return api.Complete(messages)
}
