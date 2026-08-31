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
 * - Sharing a styleguide with the account or a group (add, rename, list, revoke)
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

		// Example 5: Styleguide sharing
		// Sharing requires a multi-user account and the appropriate role (account owner for
		// account-wide shares, owner/admin for group shares). Each call returns the shared
		// styleguide, whose Name reflects the shared copy's name and SharedAt the share time.
		fmt.Println("=== Styleguide Sharing ===")

		// Share with the whole account/team; plain AddAccountShare reuses the styleguide's own name
		teamShare, err := laraTranslator.Styleguides.AddAccountShareWithName(styleguideID, "Shared with the team")
		if err != nil {
			log.Printf("Error sharing styleguide with the account: %v", err)
		} else {
			fmt.Printf("🤝 Shared with the account as: '%s' (shared at %s)\n", teamShare.Name, teamShare.SharedAt)
		}

		// Rename the account/team share
		renamedTeamShare, err := laraTranslator.Styleguides.RenameAccountShare(styleguideID, "Team styleguide")
		if err != nil {
			log.Printf("Error renaming the account share: %v", err)
		} else {
			fmt.Printf("📝 Renamed account share to: '%s'\n", renamedTeamShare.Name)
		}

		// List every share visible to the caller: the account share, group shares and user shares
		shares, err := laraTranslator.Styleguides.GetShares(styleguideID)
		if err != nil {
			log.Printf("Error listing styleguide shares: %v", err)
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
		if _, err = laraTranslator.Styleguides.RevokeAccountShare(styleguideID); err != nil {
			log.Printf("Error revoking the account share: %v", err)
		} else {
			fmt.Println("🚫 Revoked the account share")
		}

		// Group shares work the same way, addressed by a group ID (grp_...)
		if groupID := os.Getenv("LARA_GROUP_ID"); groupID != "" {
			groupShare, groupErr := laraTranslator.Styleguides.AddGroupShareWithName(styleguideID, groupID, "Shared with the group")
			if groupErr != nil {
				log.Printf("Error sharing styleguide with the group: %v", groupErr)
			} else {
				fmt.Printf("🤝 Shared with group %s as: '%s'\n", groupID, groupShare.Name)
			}

			if _, groupErr = laraTranslator.Styleguides.RenameGroupShare(styleguideID, groupID, "Marketing group"); groupErr != nil {
				log.Printf("Error renaming the group share: %v", groupErr)
			} else {
				fmt.Println("📝 Renamed the group share")
			}

			if _, groupErr = laraTranslator.Styleguides.RevokeGroupShare(styleguideID, groupID); groupErr != nil {
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
