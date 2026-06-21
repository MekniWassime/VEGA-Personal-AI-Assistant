package jobprocess

import (
	"context"
	"encoding/json"
	"fmt"
	"vega/api/ai"
	"vega/api/conversation"
	db "vega/api/internal/database"

	"github.com/jackc/pgx/v5/pgtype"
)

type ProcessConversationHandler struct{}

func (h *ProcessConversationHandler) Name() string {
	return "append_process_conversation"
}

type processConversationArgs struct {
	ConversationID string `json:"conversation_id"`
	Source         string `json:"source"`
	Message        string `json:"message"`
}

func (h *ProcessConversationHandler) Run(ctx context.Context, q *db.Queries, arguments json.RawMessage) error {
	var args processConversationArgs
	if err := json.Unmarshal(arguments, &args); err != nil {
		return fmt.Errorf("failed to parse arguments: %w — expected {\"conversation_id\": \"...\", \"source\": \"...\", \"message\": \"...\"}", err)
	}

	if args.Source != "assistant" && args.Source != "system" {
		return fmt.Errorf("invalid source %q: must be \"assistant\" or \"system\"", args.Source)
	}
	if args.Message == "" {
		return fmt.Errorf("message is required")
	}

	var conversationID pgtype.UUID
	if err := conversationID.Scan(args.ConversationID); err != nil {
		return fmt.Errorf("invalid conversation_id %q: %w", args.ConversationID, err)
	}

	client := ai.NewOllamaAPI("gemma3:4b")

	return conversation.AppendAndProcess(ctx, q, client, conversationID, ai.Message{
		Role:    args.Source,
		Content: args.Message,
	})
}
