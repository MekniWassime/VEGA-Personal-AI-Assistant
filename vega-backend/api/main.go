package main

import (
	"fmt"
	"log"
	"strings"
	"vega/api/ai"
	"vega/api/skills"
	"vega/api/system"

	"github.com/joho/godotenv"
)

const userPrompt = "find the brand of my mobile device and send it to my macbook"

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	client := ai.NewOllamaAPI("gemma3:4b")
	fmt.Println("[main] client initialized")

	conversation := []ai.Message{
		{Role: "system", Content: system.SystemPrompt},
		{Role: "user", Content: userPrompt},
	}
	fmt.Println("[main] conversation initialized, starting loop")

	for {
		fmt.Println("[main] sending request to AI...")
		fmt.Println("\n--- Conversation ---")
		for _, msg := range conversation {
			fmt.Printf("[%s]: %s\n", msg.Role, msg.Content)
		}
		fmt.Println("--------------------")
		resp, err := ai.Complete(client, conversation)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Println("[main] received response from AI")

		resp.Content = strings.TrimSpace(resp.Content)
		conversation = append(conversation, *resp)

		if resp.Content == "TASK_COMPLETE" || strings.HasSuffix(resp.Content, "TASK_COMPLETE") {
			fmt.Println("[main] task complete, exiting loop")
			break
		}

		if call, ok := skills.ExtractSkillCall(resp.Content); ok {
			fmt.Printf("[main] skill call detected: %s\n", call)
			skillResult := skills.ParseAndRun(call)
			fmt.Printf("[main] skill result: %s\n", skillResult)
			conversation = append(conversation, ai.Message{
				Role:    "system",
				Content: skillResult,
			})
		} else {
			fmt.Println("[main] no skill call detected, looping")
		}
	}
	fmt.Println("\n--- END OF Conversation ---")
	for _, msg := range conversation {
		fmt.Printf("[%s]: %s\n", msg.Role, msg.Content)
	}
	fmt.Println("--------------------")
}
