package main

import (
	"fmt"
	"log"
	"os"

	"github.com/translated/lara-go/lara"
)

/**
 * Complete styleguide management examples for the Lara Go SDK
 *
 * This example demonstrates:
 * - Create, list, get, update, delete styleguides
 * - Update name, content, or both at once
 * - Handling of non-existent styleguides
 */

func main() {
	// All examples use environment variables for credentials, so set them first:
	// export LARA_ACCESS_KEY_ID="your-access-key-id"
	// export LARA_ACCESS_KEY_SECRET="your-access-key-secret"

	// Set your credentials here
	accessKeyID := os.Getenv("LARA_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("LARA_ACCESS_KEY_SECRET")

	laraTranslator := lara.NewTranslator(lara.NewCredentials(accessKeyID, accessKeySecret), nil)

	fmt.Println("📋 Styleguides require a specific subscription plan.")
	fmt.Println("   If you encounter errors, please check your subscription level.\n")

	var styleguideID string

	err := func() error {
		// Example 1: Basic styleguide management
		fmt.Println("=== Basic Styleguide Management ===")
		styleguide, err := laraTranslator.Styleguides.Create("MyDemoStyleguide", "Use a formal tone. Prefer British English spelling. Avoid contractions.")
		if err != nil {
			return fmt.Errorf("Error creating styleguide: %v", err)
		}
		fmt.Printf("✅ Created styleguide: %s (ID: %s)\n", styleguide.Name, styleguide.ID)
		styleguideID = styleguide.ID

		// List all styleguides
		styleguides, err := laraTranslator.Styleguides.List()
		if err != nil {
			return fmt.Errorf("Error listing styleguides: %v", err)
		}
		fmt.Printf("📝 Total styleguides: %d\n", len(styleguides))
		fmt.Println()

		// Example 2: Styleguide operations
		fmt.Println("=== Styleguide Operations ===")
		retrievedStyleguide, err := laraTranslator.Styleguides.Get(styleguide.ID)
		if err != nil {
			log.Printf("Error retrieving styleguide: %v", err)
		} else if retrievedStyleguide != nil {
			fmt.Printf("📖 Styleguide: %s (Owner: %s)\n", retrievedStyleguide.Name, retrievedStyleguide.OwnerID)
			if retrievedStyleguide.Content != nil {
				fmt.Printf("📄 Content: %s\n", *retrievedStyleguide.Content)
			}
		}
		fmt.Println()

		// Example 3: Update styleguide
		fmt.Println("=== Update Styleguide ===")
		// Update only the name
		updatedName := "UpdatedDemoStyleguide"
		renamedStyleguide, err := laraTranslator.Styleguides.Update(styleguide.ID, &updatedName, nil)
		if err != nil {
			log.Printf("Error updating styleguide name: %v", err)
		} else {
			fmt.Printf("📝 Updated name: '%s' -> '%s'\n", styleguide.Name, renamedStyleguide.Name)
		}

		// Update only the content
		updatedContent := "Use a casual tone. Prefer American English spelling. Contractions are welcome."
		updatedStyleguide, err := laraTranslator.Styleguides.Update(styleguide.ID, nil, &updatedContent)
		if err != nil {
			log.Printf("Error updating styleguide content: %v", err)
		} else {
			fmt.Printf("📝 Updated content for styleguide: %s\n", updatedStyleguide.Name)
		}

		// Update both name and content
		finalName := "FinalDemoStyleguide"
		finalContent := "Use clear and concise language. Avoid jargon."
		fullyUpdated, err := laraTranslator.Styleguides.Update(styleguide.ID, &finalName, &finalContent)
		if err != nil {
			log.Printf("Error updating styleguide: %v", err)
		} else {
			fmt.Printf("📝 Updated name and content: %s\n", fullyUpdated.Name)
		}
		fmt.Println()

		// Example 4: Get a non-existent styleguide
		fmt.Println("=== Get Non-Existent Styleguide ===")
		missing, err := laraTranslator.Styleguides.Get("non-existent-id")
		if err != nil {
			log.Printf("Error getting styleguide: %v", err)
		} else if missing == nil {
			fmt.Println("ℹ️  Styleguide not found (returned nil as expected)")
		}
		fmt.Println()

		return nil
	}()

	// Cleanup (equivalent to finally block)
	fmt.Println("=== Cleanup ===")
	if styleguideID != "" {
		deleted, deleteErr := laraTranslator.Styleguides.Delete(styleguideID)
		if deleteErr != nil {
			log.Printf("Error deleting styleguide: %v", deleteErr)
		} else {
			fmt.Printf("🗑️  Deleted styleguide: %s\n", deleted.Name)
		}
	}

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("\n🎉 Styleguide management examples completed!")
}
