package main

import (
	"context"
	"log"
	"os"
	"vega/api/ai"
	"vega/api/conversation"
	db "vega/api/internal/database"
	"vega/api/worker"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

const userPrompt = "figure out the brand of my mobile device and send it to my macbook"

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	if err := worker.Init(ctx, queries); err != nil {
		log.Fatalf("failed to initialize worker: %v", err)
	}

	client := ai.NewOllamaAPI("gemma3:4b")

	conversationID, err := conversation.StartConversation(ctx, queries, db.ConversationTypeTask, ai.Message{
		Role:    "user",
		Content: userPrompt,
	})
	if err != nil {
		log.Fatalf("failed to start conversation: %v", err)
	}

	if err := conversation.ProcessConversation(ctx, queries, client, conversationID); err != nil {
		log.Fatalf("conversation error: %v", err)
	}
}
