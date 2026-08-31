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
 * - Glossary export (sync and async)
 * - Glossary terms count
 * - Sharing a glossary with the account or a group (add, rename, list, revoke)
 * - Import status checking
 * - Add or replace glossary entries
 * - Delete glossary entries
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

		// Example 4: CSV import with a callback URL (async notification when import completes)
		fmt.Println("=== CSV Import with Callback URL ===")
		if _, err := os.Stat(csvFilePath); err == nil {
			callbackUrl := "https://your-server.example.com/lara/import-callback" // Replace with your endpoint
			// Note: the callback variants require an explicit content type, so pass one even for the default format.
			importWithCallback, err := laraTranslator.Glossaries.ImportCsvFromPathWithFormatAndCallback(glossary.ID, csvFilePath, "csv/table-uni", callbackUrl)
			if err != nil {
				log.Printf("Error starting CSV import with callback: %v", err)
			} else {
				fmt.Printf("Import started with ID: %s (callback: %s)\n", importWithCallback.ID, callbackUrl)
			}
			fmt.Println()
		} else {
			fmt.Printf("CSV file not found: %s\n", csvFilePath)
		}

		// Example 5: Export functionality
		fmt.Println("=== Export Functionality ===")

		// Export as CSV table unidirectional format
		fmt.Println("📤 Exporting as CSV table unidirectional...")
		source := "en-US"
		csvUniData, err := laraTranslator.Glossaries.Export(glossary.ID, "csv/table-uni", &source)
		if err != nil {
			fmt.Printf("Error with export: %v\n", err)
		} else {
			fmt.Printf("✅ CSV unidirectional export successful (%d bytes)\n", len(csvUniData))

			// Save sample exports to files - replace with your desired output paths
			exportFilePath := filepath.Join(".", "exported_glossary.csv") // Replace with actual path
			err = os.WriteFile(exportFilePath, csvUniData, 0644)
			if err != nil {
				fmt.Printf("Error saving export file: %v\n", err)
			} else {
				fmt.Printf("💾 Sample export saved to: %s\n", filepath.Base(exportFilePath))
			}

			// Async export - returns a job ID; the result is delivered to your callback URL when ready
			fmt.Println("📤 Starting async export...")
			exportJob, err := laraTranslator.Glossaries.ExportAsync(
				glossary.ID,
				"https://your-server.example.com/lara/export-callback", // Replace with your actual callback URL
				"csv/table-uni",
				&source,
			)
			if err != nil {
				fmt.Printf("Error starting async export: %v\n", err)
			} else {
				fmt.Printf("✅ Async export started (Job ID: %s)\n", exportJob.JobID)
				fmt.Println("   The export result will be delivered to your callback URL when ready.")
			}
		}
		fmt.Println()

		// Example 6: Glossary Terms Count
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

		// Example 7: Add or replace glossary entries
		fmt.Println("=== Add or Replace Glossary Entries ===")

		terms := []lara.GlossaryTerm{
			{Language: "en-US", Value: "computer"},
			{Language: "it-IT", Value: "computer"},
		}
		addResult, err := laraTranslator.Glossaries.AddOrReplaceEntry(glossary.ID, terms, nil)
		if err != nil {
			log.Printf("Error adding/replacing entry: %v", err)
		} else {
			fmt.Printf("✅ Entry added/replaced (import ID: %s)\n", addResult.ID)
		}

		customGUID := "custom-guid-123"
		termsWithGUID := []lara.GlossaryTerm{
			{Language: "en-US", Value: "keyboard"},
			{Language: "it-IT", Value: "tastiera"},
		}
		addWithGUIDResult, err := laraTranslator.Glossaries.AddOrReplaceEntry(glossary.ID, termsWithGUID, &customGUID)
		if err != nil {
			log.Printf("Error adding entry with GUID: %v", err)
		} else {
			fmt.Printf("✅ Entry added with GUID (import ID: %s)\n", addWithGUIDResult.ID)
		}

		updatedTerms := []lara.GlossaryTerm{
			{Language: "en-US", Value: "keyboard"},
			{Language: "it-IT", Value: "tastiera"},
			{Language: "fr-FR", Value: "clavier"},
		}
		replaceResult, err := laraTranslator.Glossaries.AddOrReplaceEntry(glossary.ID, updatedTerms, &customGUID)
		if err != nil {
			log.Printf("Error replacing entry: %v", err)
		} else {
			fmt.Printf("✅ Entry replaced with updated terms (import ID: %s)\n", replaceResult.ID)
		}
		fmt.Println()

		// Example 8: Delete glossary entries
		fmt.Println("=== Delete Glossary Entries ===")

		deleteByGUIDResult, err := laraTranslator.Glossaries.DeleteEntry(glossary.ID, nil, &customGUID)
		if err != nil {
			log.Printf("Error deleting entry by GUID: %v", err)
		} else {
			fmt.Printf("✅ Entry deleted by GUID (import ID: %s)\n", deleteByGUIDResult.ID)
		}

		term := &lara.GlossaryTerm{Language: "en-US", Value: "computer"}
		deleteByTermResult, err := laraTranslator.Glossaries.DeleteEntry(glossary.ID, term, nil)
		if err != nil {
			log.Printf("Error deleting entry by term: %v", err)
		} else {
			fmt.Printf("✅ Entry deleted by term: %s -> \"%s\" (import ID: %s)\n", term.Language, term.Value, deleteByTermResult.ID)
		}
		fmt.Println()

		// Example 9: Glossary sharing
		// Sharing requires a multi-user account and the appropriate role (account owner for
		// account-wide shares, owner/admin for group shares). Each call returns the shared
		// glossary, whose Name reflects the shared copy's name and SharedAt the share time.
		fmt.Println("=== Glossary Sharing ===")

		// Share with the whole account/team; plain AddAccountShare reuses the glossary's own name
		teamShare, err := laraTranslator.Glossaries.AddAccountShareWithName(glossary.ID, "Shared with the team")
		if err != nil {
			log.Printf("Error sharing glossary with the account: %v", err)
		} else {
			fmt.Printf("🤝 Shared with the account as: '%s' (shared at %s)\n", teamShare.Name, teamShare.SharedAt)
		}

		// Rename the account/team share
		renamedTeamShare, err := laraTranslator.Glossaries.RenameAccountShare(glossary.ID, "Team glossary")
		if err != nil {
			log.Printf("Error renaming the account share: %v", err)
		} else {
			fmt.Printf("📝 Renamed account share to: '%s'\n", renamedTeamShare.Name)
		}

		// List every share visible to the caller: the account share, group shares and user shares
		shares, err := laraTranslator.Glossaries.GetShares(glossary.ID)
		if err != nil {
			log.Printf("Error listing glossary shares: %v", err)
		} else {
			if shares.Account != nil {
				fmt.Printf("👥 Account share '%s' (%s)\n", shares.Account.ShareName, shares.Account.Permissions)
			}
			for _, group := range shares.Groups {
				fmt.Printf("👥 Group %s: '%s' (%s)\n", group.Name, group.ShareName, group.Permissions)
			}
			for _, user := range shares.Users {
				fmt.Printf("👤 User %s: '%s' (%s)\n", user.Name, user.ShareName, user.Permissions)
			}
		}

		// Revoke the account/team share
		if _, err = laraTranslator.Glossaries.RevokeAccountShare(glossary.ID); err != nil {
			log.Printf("Error revoking the account share: %v", err)
		} else {
			fmt.Println("🚫 Revoked the account share")
		}

		// Group shares work the same way, addressed by a group ID (grp_...)
		if groupID := os.Getenv("LARA_GROUP_ID"); groupID != "" {
			groupShare, groupErr := laraTranslator.Glossaries.AddGroupShareWithName(glossary.ID, groupID, "Shared with the group")
			if groupErr != nil {
				log.Printf("Error sharing glossary with the group: %v", groupErr)
			} else {
				fmt.Printf("🤝 Shared with group %s as: '%s'\n", groupID, groupShare.Name)
			}

			if _, groupErr = laraTranslator.Glossaries.RenameGroupShare(glossary.ID, groupID, "Marketing group"); groupErr != nil {
				log.Printf("Error renaming the group share: %v", groupErr)
			} else {
				fmt.Println("📝 Renamed the group share")
			}

			if _, groupErr = laraTranslator.Glossaries.RevokeGroupShare(glossary.ID, groupID); groupErr != nil {
				log.Printf("Error revoking the group share: %v", groupErr)
			} else {
				fmt.Println("🚫 Revoked the group share")
			}
		} else {
			fmt.Println("Set LARA_GROUP_ID to try the group sharing methods.")
		}
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
