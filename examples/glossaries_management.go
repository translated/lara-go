package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/translated/lara-go/lara"
)

/**
 * Complete glossary management examples for the Lara Go SDK
 *
 * This example demonstrates:
 * - Create, list, update, delete glossaries
 * - CSV import with status monitoring
 * - Glossary export
 * - Glossary terms count
 * - Import status checking
 */

func main() {
	// All examples use environment variables for credentials, so set them first:
	// export LARA_ACCESS_KEY_ID="your-access-key-id"
	// export LARA_ACCESS_KEY_SECRET="your-access-key-secret"

	// Set your credentials here
	accessKeyID := os.Getenv("LARA_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("LARA_ACCESS_KEY_SECRET")

	laraTranslator := lara.NewTranslator(lara.NewCredentials(accessKeyID, accessKeySecret), nil)

	fmt.Println("🗒️  Glossaries require a specific subscription plan.")
	fmt.Println("   If you encounter errors, please check your subscription level.\n")

	var glossaryID string

	err := func() error {
		// Example 1: Basic glossary management
		fmt.Println("=== Basic Glossary Management ===")
		glossary, err := laraTranslator.Glossaries.Create("MyDemoGlossary")
		if err != nil {
			return fmt.Errorf("Error creating glossary: %v", err)
		}
		fmt.Printf("✅ Created glossary: %s (ID: %s)\n", glossary.Name, glossary.ID)
		glossaryID = glossary.ID

		// List all glossaries
		glossaries, err := laraTranslator.Glossaries.List()
		if err != nil {
			return fmt.Errorf("Error listing glossaries: %v", err)
		}
		fmt.Printf("📝 Total glossaries: %d\n", len(glossaries))
		fmt.Println()

		// Example 2: Glossary operations
		fmt.Println("=== Glossary Operations ===")
		// Get glossary details
		retrievedGlossary, err := laraTranslator.Glossaries.Get(glossary.ID)
		if err != nil {
			log.Printf("Error retrieving glossary: %v", err)
		} else if retrievedGlossary != nil {
			fmt.Printf("📖 Glossary: %s (Owner: %s)\n", retrievedGlossary.Name, retrievedGlossary.OwnerID)
		}

		// Get glossary terms count
		counts, err := laraTranslator.Glossaries.Counts(glossary.ID)
		if err != nil {
			log.Printf("Error getting glossary counts: %v", err)
		} else {
			if counts.Unidirectional != nil {
				for lang, count := range counts.Unidirectional {
					fmt.Printf("   %s: %d entries\n", lang, count)
				}
			}
		}

		// Update glossary
		updatedGlossary, err := laraTranslator.Glossaries.Update(glossary.ID, "UpdatedDemoGlossary")
		if err != nil {
			log.Printf("Error updating glossary: %v", err)
		} else {
			fmt.Printf("📝 Updated name: '%s' -> '%s'\n", glossary.Name, updatedGlossary.Name)
		}

		// Example 3: CSV import functionality
		fmt.Println("=== CSV Import Functionality ===")

		// Replace with your actual CSV file path
		csvFilePath := filepath.Join(".", "sample_glossary.csv") // Create this file with your glossary data

		if _, err := os.Stat(csvFilePath); err == nil {
			fmt.Printf("Importing CSV file: %s\n", filepath.Base(csvFilePath))
			csvImport, err := laraTranslator.Glossaries.ImportCsvFromPath(glossary.ID, csvFilePath)
			if err != nil {
				log.Printf("Error with CSV import: %v", err)
			} else {
				fmt.Printf("Import started with ID: %s\n", csvImport.ID)
				fmt.Printf("Initial progress: %.0f%%\n", csvImport.Progress*100)

				// Check import status manually
				fmt.Println("Checking import status...")
				importStatus, err := laraTranslator.Glossaries.GetImportStatus(csvImport.ID)
				if err != nil {
					log.Printf("Error checking import status: %v", err)
				} else {
					fmt.Printf("Current progress: %.0f%%\n", importStatus.Progress*100)
				}

				// Wait for import to complete
				maxWaitTime := 300 * time.Second // 5 minutes
				completedImport, err := laraTranslator.Glossaries.WaitForImport(csvImport, nil, &maxWaitTime)
				if err != nil {
					log.Printf("Import timeout: The import process took too long to complete.")
				} else {
					fmt.Println("✅ Import completed!")
					fmt.Printf("Final progress: %.0f%%\n", completedImport.Progress*100)
				}
			}
			fmt.Println()
		} else {
			fmt.Printf("CSV file not found: %s\n", csvFilePath)
		}

		// Example 4: Export functionality
		fmt.Println("=== Export Functionality ===")

		// Export as CSV table unidirectional format
		fmt.Println("📤 Exporting as CSV table unidirectional...")
		source := "en-US"
		csvUniData, err := laraTranslator.Glossaries.Export(glossary.ID, "csv/table-uni", &source)
		if err != nil {
			return fmt.Errorf("error exporting as CSV table unidirectional: %v", err)
		}
		fmt.Printf("✅ CSV unidirectional export successful (%d bytes)\n", len(csvUniData))

		// Save sample exports to files - replace with your desired output paths
		exportFilePath := filepath.Join(".", "exported_glossary.csv") // Replace with actual path
		err = os.WriteFile(exportFilePath, csvUniData, 0644)
		if err != nil {
			return fmt.Errorf("error saving export file: %v", err)
		}
		fmt.Printf("💾 Sample export saved to: %s\n", filepath.Base(exportFilePath))
		fmt.Println()

		// Example 5: Glossary Terms Count
		fmt.Println("=== Glossary Terms Count ===")

		// Get detailed counts
		detailedCounts, err := laraTranslator.Glossaries.Counts(glossary.ID)
		if err != nil {
			return fmt.Errorf("error getting glossary terms count: %v", err)
		}

		fmt.Println("📊 Detailed glossary terms count:")

		if detailedCounts.Unidirectional != nil && len(detailedCounts.Unidirectional) > 0 {
			fmt.Println("   Unidirectional entries by language pair:")
			for langPair, count := range detailedCounts.Unidirectional {
				fmt.Printf("     %s: %d terms\n", langPair, count)
			}
		} else {
			fmt.Println("   No unidirectional entries found")
		}

		totalEntries := 0
		if detailedCounts.Unidirectional != nil {
			for _, count := range detailedCounts.Unidirectional {
				totalEntries += count
			}
		}
		fmt.Printf("   Total entries: %d\n", totalEntries)
		fmt.Println()

		return nil
	}()

	// Cleanup (equivalent to finally block)
	fmt.Println("=== Cleanup ===")
	if glossaryID != "" {
		_, deleteErr := laraTranslator.Glossaries.Delete(glossaryID)
		if deleteErr != nil {
			log.Printf("Error deleting glossary: %v", deleteErr)
		} else {
			fmt.Printf("🗑️  Deleted glossary: %s\n", "MyDemoGlossary")

			// Clean up export files - replace with actual cleanup if needed
			exportFilePath := filepath.Join(".", "exported_glossary.csv")
			if _, statErr := os.Stat(exportFilePath); statErr == nil {
				removeErr := os.Remove(exportFilePath)
				if removeErr != nil {
					log.Printf("Error cleaning up export file: %v", removeErr)
				} else {
					fmt.Println("🗑️  Cleaned up export file")
				}
			}
		}
	}

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("\n🎉 Glossary management examples completed!")
}
