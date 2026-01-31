package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	baseURL := getEnv("POWERMEM_BASE_URL", "http://localhost:8000")
	apiKey := getEnv("POWERMEM_API_KEY", "")

	client := NewClient(baseURL, apiKey)

	fmt.Println("=== PowerMem Go Client Example ===")

	// 1. Create Memory
	fmt.Println("1. Creating memory...")
	memories, err := client.CreateMemory(CreateMemoryRequest{
		Content: "User likes coffee and goes to Starbucks every morning",
		UserID:  "user-123",
		AgentID: "agent-456",
		RunID:   "run-789",
		Metadata: map[string]interface{}{
			"source":     "conversation",
			"importance": "high",
		},
		Filters: map[string]interface{}{
			"category": "preference",
			"topic":    "beverage",
		},
		Scope:      "user",
		MemoryType: "preference",
		Infer:      true,
	})
	if err != nil {
		log.Fatalf("Failed to create memory: %v", err)
	}

	fmt.Printf("   Created %d memories\n", len(memories))
	for i, mem := range memories {
		fmt.Printf("   Memory %d - ID: %s, Content: %s\n", i+1, mem.MemoryID, mem.Content)
	}

	// 2. List Memories
	fmt.Println("\n2. Listing memories...")
	memoryList, err := client.ListMemories(ListMemoriesRequest{
		UserID:  "user-123",
		AgentID: "agent-456",
		Limit:   10,
		Offset:  0,
	})
	if err != nil {
		log.Fatalf("Failed to list memories: %v", err)
	}

	fmt.Printf("   Found %d memories (total: %d)\n", len(memoryList.Memories), memoryList.Total)
	for i, mem := range memoryList.Memories {
		fmt.Printf("   [%d] ID: %s, UserID: %s, AgentID: %s, Content: %s\n",
			i+1, mem.MemoryID, mem.UserID, mem.AgentID, mem.Content)
	}

	// 3. Search Memories
	fmt.Println("\n3. Searching memories...")
	searchResults, err := client.SearchMemories(SearchMemoryRequest{
		Query:   "What beverages does the user like",
		UserID:  "user-123",
		AgentID: "agent-456",
		Limit:   10,
	})
	if err != nil {
		log.Fatalf("Failed to search memories: %v", err)
	}

	fmt.Printf("   Found %d search results\n", len(searchResults.Results))
	for i, result := range searchResults.Results {
		fmt.Printf("   [%d] ID: %s, Score: %.4f, Content: %s\n",
			i+1, result.MemoryID, result.Score, result.Content)
	}

	// 4. Update Memory (use newly created memory to avoid old data issues)
	memoryID := memories[0].MemoryID
	fmt.Printf("\n4. Updating newly created memory (ID: %s)...\n", memoryID)

	updatedMemory, err := client.UpdateMemory(memoryID, UpdateMemoryRequest{
		Content: "User loves latte coffee and visits Starbucks daily",
		UserID:  "user-123",
		AgentID: "agent-456",
		Metadata: map[string]interface{}{
			"source":     "conversation",
			"importance": "very_high",
			"updated":    true,
		},
	})
	if err != nil {
		log.Printf("   Failed to update memory: %v", err)
	} else {
		fmt.Printf("   Updated memory ID: %s\n", updatedMemory.MemoryID)
		fmt.Printf("   New content: %s\n", updatedMemory.Content)
	}

	// 5. Delete Memory
	fmt.Printf("\n5. Deleting memory (ID: %s)...\n", memoryID)
	err = client.DeleteMemory(memoryID, "user-123", "agent-456")
	if err != nil {
		log.Printf("   Failed to delete memory: %v", err)
	} else {
		fmt.Printf("   Successfully deleted memory ID: %s\n", memoryID)
	}

	fmt.Println("\n=== Example completed ===")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
