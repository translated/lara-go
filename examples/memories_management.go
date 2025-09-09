package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/translated/lara-go/lara"
)

/**
 * Complete translation memory management examples for the Lara Go SDK
 *
 * This example demonstrates:
 * - Create, list, update, delete memories
 * - Add individual translations
 * - Multiple memory operations
 * - TMX file import with progress monitoring
 * - Translation deletion
 * - Translation with TUID and context
 */

func main() {
	// All examples use environment variables for credentials, so set them first:
	// export LARA_ACCESS_KEY_ID="your-access-key-id"
	// export LARA_ACCESS_KEY_SECRET="your-access-key-secret"

	// Set your credentials here
	accessKeyID := os.Getenv("LARA_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("LARA_ACCESS_KEY_SECRET")

	laraTranslator := lara.NewTranslator(lara.NewCredentials(accessKeyID, accessKeySecret), nil)

	var memory2 *lara.Memory

	// Example 1: Basic memory management
	fmt.Println("=== Basic Memory Management ===")
	memory, err := laraTranslator.Memories.Create("MyDemoMemory")
	if err != nil {
		log.Printf("Error creating memory: %v", err)
		return
	}
	fmt.Printf("✅ Created memory: %s (ID: %s)\n", memory.Name, memory.ID)

	// Get memory details
	retrievedMemory, err := laraTranslator.Memories.Get(memory.ID)
	if err != nil {
		log.Printf("Error retrieving memory: %v", err)
		return
	}
	if retrievedMemory != nil {
		fmt.Printf("📖 Memory: %s (Owner: %s)\n", retrievedMemory.Name, retrievedMemory.OwnerID)
	}

	// Update memory
	updatedMemory, err := laraTranslator.Memories.Update(memory.ID, "UpdatedDemoMemory")
	if err != nil {
		log.Printf("Error updating memory: %v", err)
		return
	}
	fmt.Printf("📝 Updated name: '%s' -> '%s'\n", memory.Name, updatedMemory.Name)

	// List all memories
	memories, err := laraTranslator.Memories.List()
	if err != nil {
		log.Printf("Error listing memories: %v", err)
		return
	}
	fmt.Printf("📝 Total memories: %d\n", len(memories))
	fmt.Println()

	// Example 2: Adding translations
	// Important: To update/overwrite a translation unit you must provide a tuid. Calls without a tuid always create a new unit and will not update existing entries.
	fmt.Println("=== Adding Translations ===")

	// Basic translation addition (with TUID)
	memImport1, err := laraTranslator.Memories.AddTranslationWithTuid(memory.ID, "en-US", "fr-FR", "Hello", "Bonjour", "greeting_001")
	if err != nil {
		log.Printf("Error adding translation: %v", err)
	} else {
		fmt.Printf("✅ Added: 'Hello' -> 'Bonjour' with TUID 'greeting_001' (Import ID: %s)\n", memImport1.ID)
	}

	// Translation with context
	memImport2, err := laraTranslator.Memories.AddTranslationWithContext(
		memory.ID, "en-US", "fr-FR", "How are you?", "Comment allez-vous?", "greeting_002",
		"Good morning", "Have a nice day",
	)
	if err != nil {
		log.Printf("Error adding translation with context: %v", err)
	} else {
		fmt.Printf("✅ Added with context (Import ID: %s)\n", memImport2.ID)
	}
	fmt.Println()

	// Example 3: Multiple memory operations
	fmt.Println("=== Multiple Memory Operations ===")

	// Create second memory for multi-memory operations
	memory2, err = laraTranslator.Memories.Create("SecondDemoMemory")
	if err != nil {
		log.Printf("Error creating second memory: %v", err)
	} else {
		fmt.Printf("✅ Created second memory: %s\n", memory2.Name)

		// Add translation to multiple memories (with TUID)
		memoryIds := []string{memory.ID, memory2.ID}
		multiImportJob, err := laraTranslator.Memories.AddTranslationMultipleWithTuid(memoryIds, "en-US", "it-IT", "Hello World!", "Ciao Mondo!", "greeting_003")
		if err != nil {
			log.Printf("Error adding translation to multiple memories: %v", err)
		} else {
			fmt.Printf("✅ Added translation to multiple memories (Import ID: %s)\n", multiImportJob.ID)
		}
		fmt.Println()
	}

	// Example 4: TMX import functionality
	fmt.Println("=== TMX Import Functionality ===")

	// Replace with your actual TMX file path
	tmxFilePath := "sample_memory.tmx" // Create this file with your TMX content

	if _, err := os.Stat(tmxFilePath); err == nil {
		fmt.Printf("Importing TMX file: %s\n", tmxFilePath)
		tmxImport, err := laraTranslator.Memories.ImportTmxFromPath(memory.ID, tmxFilePath)
		if err != nil {
			log.Printf("Error with TMX import: %v", err)
		} else {
			fmt.Printf("Import started with ID: %s\n", tmxImport.ID)
			fmt.Printf("Initial progress: %.0f%%\n", tmxImport.Progress*100)

			// Wait for import to complete
			maxWaitTime := 300 * time.Second // 5 minutes
			completedImport, err := laraTranslator.Memories.WaitForImport(tmxImport, nil, &maxWaitTime)
			if err != nil {
				log.Printf("Import timeout: The import process took too long to complete.")
			} else {
				fmt.Println("✅ Import completed!")
				fmt.Printf("Final progress: %.0f%%\n", completedImport.Progress*100)
			}
		}
		fmt.Println()
	} else {
		fmt.Printf("TMX file not found: %s\n", tmxFilePath)
	}

	// Example 5: Translation deletion
	fmt.Println("=== Translation Deletion ===")

	// Delete translation unit
	// Important: without TUID, all entries that match the provided fields will be removed
	deleteJob, err := laraTranslator.Memories.DeleteTranslation(
		memory.ID,
		"en-US",
		"fr-FR",
		"Hello",
		"Bonjour",
	)
	if err != nil {
		log.Printf("Error deleting translation: %v", err)
	} else {
		fmt.Printf("🗑️  Deleted translation unit (Job ID: %s)\n", deleteJob.ID)
	}
	fmt.Println()

	// Cleanup
	fmt.Println("=== Cleanup ===")
	_, err = laraTranslator.Memories.Delete(memory.ID)
	if err != nil {
		log.Printf("Error deleting memory: %v", err)
	} else {
		fmt.Printf("🗑️  Deleted memory: %s\n", memory.Name)
	}

	if memory2 != nil {
		_, err = laraTranslator.Memories.Delete(memory2.ID)
		if err != nil {
			log.Printf("Error deleting second memory: %v", err)
		} else {
			fmt.Printf("🗑️  Deleted second memory: %s\n", memory2.Name)
		}
	}

	fmt.Println("\n🎉 Memory management examples completed!")
}
