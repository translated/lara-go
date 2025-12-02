package main

import (
	"fmt"
	"log"
	"os"

	"github.com/translated/lara-go/lara"
)

/**
 * Language detection examples for the Lara Go SDK
 *
 * This example demonstrates:
 * - Detect language from a single string
 * - Detect language from multiple strings
 * - Detect with hint parameter
 * - Detect with passlist parameter
 */

func main() {
	// All examples use environment variables for credentials, so set them first:
	// export LARA_ACCESS_KEY_ID="your-access-key-id"
	// export LARA_ACCESS_KEY_SECRET="your-access-key-secret"

	// Set your credentials here
	accessKeyID := os.Getenv("LARA_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("LARA_ACCESS_KEY_SECRET")

	laraTranslator := lara.NewTranslator(lara.NewCredentials(accessKeyID, accessKeySecret), nil)

	// Example 1: Basic language detection from a single string
	fmt.Println("=== Basic Language Detection ===")
	result, err := laraTranslator.Detect("Bonjour le monde!", "", nil)
	if err != nil {
		log.Printf("Error with basic detection: %v", err)
		return
	}
	fmt.Printf("Text: Bonjour le monde!\n")
	fmt.Printf("Detected language: %s\n", result.Language)
	fmt.Printf("Content type: %s\n\n", result.ContentType)

	// Example 2: Detect language from multiple strings
	fmt.Println("=== Multiple Strings Detection ===")
	texts := []string{"Hello world", "How are you?", "Goodbye"}
	multiResult, err := laraTranslator.Detect(texts, "", nil)
	if err != nil {
		log.Printf("Error with multiple strings detection: %v", err)
		return
	}
	fmt.Printf("Texts: %v\n", texts)
	fmt.Printf("Detected language: %s\n", multiResult.Language)
	fmt.Printf("Content type: %s\n\n", multiResult.ContentType)

	// Example 3: Detection with hint parameter
	fmt.Println("=== Detection with Hint ===")
	hintResult, err := laraTranslator.Detect("Ciao mondo!", "it-IT", nil)
	if err != nil {
		log.Printf("Error with hint detection: %v", err)
		return
	}
	fmt.Printf("Text: Ciao mondo!\n")
	fmt.Printf("Hint: it-IT\n")
	fmt.Printf("Detected language: %s\n", hintResult.Language)
	fmt.Printf("Content type: %s\n\n", hintResult.ContentType)

	// Example 4: Detection with passlist
	fmt.Println("=== Detection with Passlist ===")
	passlist := []string{"en-US", "fr-FR", "es-ES"}
	passlistResult, err := laraTranslator.Detect("Hola mundo!", "", passlist)
	if err != nil {
		log.Printf("Error with passlist detection: %v", err)
		return
	}
	fmt.Printf("Text: Hola mundo!\n")
	fmt.Printf("Passlist: %v\n", passlist)
	fmt.Printf("Detected language: %s\n", passlistResult.Language)
	fmt.Printf("Content type: %s\n\n", passlistResult.ContentType)

	// Example 5: Detection with both hint and passlist
	fmt.Println("=== Detection with Hint and Passlist ===")
	combinedPasslist := []string{"de-DE", "en-US", "fr-FR"}
	combinedResult, err := laraTranslator.Detect("Guten Tag!", "de-DE", combinedPasslist)
	if err != nil {
		log.Printf("Error with combined detection: %v", err)
		return
	}
	fmt.Printf("Text: Guten Tag!\n")
	fmt.Printf("Hint: de-DE\n")
	fmt.Printf("Passlist: %v\n", combinedPasslist)
	fmt.Printf("Detected language: %s\n", combinedResult.Language)
	fmt.Printf("Content type: %s\n", combinedResult.ContentType)
}
