package conversation

import (
	"context"
	"fmt"
	"strings"
	"vega/api/ai"
	db "vega/api/internal/database"
	"vega/api/skills"

	"github.com/jackc/pgx/v5/pgtype"
)

func ProcessConversation(ctx context.Context, q *db.Queries, client ai.ModelAPI, conversationID pgtype.UUID) error {
	messages, err := ConstructContext(ctx, q, conversationID)
	if err != nil {
		return fmt.Errorf("failed to load context: %w", err)
	}

	for {
		resp, err := ai.Complete(client, messages)
		if err != nil {
			return fmt.Errorf("completion error: %w", err)
		}

		resp.Content = strings.TrimSpace(resp.Content)

		if err := persistMessage(ctx, q, conversationID, &messages, *resp); err != nil {
			return fmt.Errorf("failed to persist assistant message: %w", err)
		}

		if strings.HasSuffix(resp.Content, "TASK_COMPLETE") {
			break
		}

		if call, ok := skills.ExtractSkillCall(resp.Content); ok {
			result := skills.ParseAndRun(call)
			if result.Suspend {
				break
			}
			if err := persistMessage(ctx, q, conversationID, &messages, ai.Message{Role: "system", Content: result.Content}); err != nil {
				return fmt.Errorf("failed to persist skill result: %w", err)
			}
		}
		printMessages(messages)
	}
	printMessages(messages)

	return nil
}

func persistMessage(ctx context.Context, q *db.Queries, conversationID pgtype.UUID, messages *[]ai.Message, msg ai.Message) error {
	if err := AppendMessage(ctx, q, conversationID, msg); err != nil {
		return err
	}
	*messages = append(*messages, msg)
	return nil
}

func printMessages(messages []ai.Message) {
	fmt.Println("\n--- Conversation ---")
	for _, msg := range messages {
		fmt.Printf("[%s]: %s\n", msg.Role, msg.Content)
	}
	fmt.Println("--------------------")
}
