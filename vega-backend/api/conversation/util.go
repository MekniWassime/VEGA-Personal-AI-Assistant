package conversation

import (
	"context"
	"vega/api/ai"
	db "vega/api/internal/database"
	"vega/api/system"

	"github.com/jackc/pgx/v5/pgtype"
)

func StartConversation(ctx context.Context, q *db.Queries, conversationType db.ConversationType, firstMessage ai.Message) (pgtype.UUID, error) {
	conv, err := q.CreateConversation(ctx, conversationType)
	if err != nil {
		return pgtype.UUID{}, err
	}

	dbCtx, err := q.CreateContext(ctx, conv.ID)
	if err != nil {
		return pgtype.UUID{}, err
	}

	_, err = q.CreateMessage(ctx, db.CreateMessageParams{
		ContextID: dbCtx.ID,
		Role:      firstMessage.Role,
		Content:   firstMessage.Content,
	})
	if err != nil {
		return pgtype.UUID{}, err
	}

	return conv.ID, nil
}

func AppendMessage(ctx context.Context, q *db.Queries, conversationID pgtype.UUID, message ai.Message) error {
	dbCtx, err := q.GetContextByConversation(ctx, conversationID)
	if err != nil {
		return err
	}

	_, err = q.CreateMessage(ctx, db.CreateMessageParams{
		ContextID: dbCtx.ID,
		Role:      message.Role,
		Content:   message.Content,
	})
	return err
}

func ConstructContext(ctx context.Context, q *db.Queries, conversationID pgtype.UUID) ([]ai.Message, error) {
	messages, err := q.ListMessagesByConversation(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	result := make([]ai.Message, 0, len(messages)+1)
	result = append(result, ai.Message{Role: "system", Content: system.SystemPrompt})

	for _, m := range messages {
		result = append(result, ai.Message{Role: m.Role, Content: m.Content})
	}

	return result, nil
}
