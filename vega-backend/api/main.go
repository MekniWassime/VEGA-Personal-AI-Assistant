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

const userPrompt = "Give me the brand of my personal mobile phone, send this result to my macbook device"

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	client := ai.NewOllamaAPI("qwen3-vl:8b")
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
				Content: "Here is the result of the tool call you just executed:\n" + skillResult + "\nWith this result, continue your task based on the response.",
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
