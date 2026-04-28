package main

import (
	"fmt"
	"log"
	"os"

	"github.com/translated/lara-go/lara"
)

/**
 * Complete text translation examples for the Lara Go SDK
 *
 * This example demonstrates:
 * - Single string translation
 * - Multiple strings translation
 * - Translation with instructions
 * - TextBlocks translation (mixed translatable/non-translatable content)
 * - Auto-detect source language
 * - Advanced translation options
 * - Profanity detection options
 * - Get available languages
 */

func main() {
	// All examples use environment variables for credentials, so set them first:
	// export LARA_ACCESS_KEY_ID="your-access-key-id"
	// export LARA_ACCESS_KEY_SECRET="your-access-key-secret"

	// Set your credentials here
	accessKeyID := os.Getenv("LARA_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("LARA_ACCESS_KEY_SECRET")

	laraTranslator := lara.NewTranslator(lara.NewCredentials(accessKeyID, accessKeySecret), nil)

	// Example 1: Basic single string translation
	fmt.Println("=== Basic Single String Translation ===")
	result, err := laraTranslator.Translate("Hello, world!", "en-US", "fr-FR", lara.TranslateOptions{})
	if err != nil {
		log.Printf("Translation error: %v", err)
		return
	}
	fmt.Printf("Original: Hello, world!\n")
	fmt.Printf("French: %s\n\n", *result.Translation.String)

	// Example 2: Multiple strings translation
	fmt.Println("=== Multiple Strings Translation ===")
	texts := []string{"Hello", "How are you?", "Goodbye"}
	multiResult, err := laraTranslator.Translate(texts, "en-US", "es-ES", lara.TranslateOptions{})
	if err != nil {
		log.Printf("Error translating multiple texts: %v", err)
		return
	}
	fmt.Printf("Original: %v\n", texts)
	fmt.Printf("Spanish: %v\n\n", multiResult.Translation.Strings)

	// Example 3: TextBlocks translation (mixed translatable/non-translatable content)
	fmt.Println("=== TextBlocks Translation ===")
	textBlocks := []lara.TextBlock{
		{Text: "Adventure novels, mysteries, cookbooks—wait, who packed those?", Translatable: true},
		{Text: "<br>", Translatable: false}, // Non-translatable HTML
		{Text: "Suddenly, it doesn't feel so deserted after all.", Translatable: true},
		{Text: "<div class=\"separator\"></div>", Translatable: false}, // Non-translatable HTML
		{Text: "Every page you turn is a new journey, and the best part?", Translatable: true},
	}

	textBlockResult, err := laraTranslator.Translate(textBlocks, "en-US", "it-IT", lara.TranslateOptions{})
	if err != nil {
		log.Printf("Error with TextBlocks translation: %v", err)
		return
	}
	fmt.Printf("Original TextBlocks: %d blocks\n", len(textBlocks))
	fmt.Printf("Translated blocks: %d\n", len(textBlockResult.Translation.TextBlocks))
	for i, block := range textBlockResult.Translation.TextBlocks {
		fmt.Printf("Block %d: %s\n", i+1, block.Text)
	}
	fmt.Println()

	// Example 4: Translation with instructions
	fmt.Println("=== Translation with Instructions ===")
	instructedResult, err := laraTranslator.Translate(
		"Could you send me the report by tomorrow morning?",
		"en-US",
		"de-DE",
		lara.TranslateOptions{
			Instructions: []string{"Be formal", "Use technical terminology"},
			ContentType:  "text/plain",
		},
	)
	if err != nil {
		log.Printf("Error with instructed translation: %v", err)
		return
	}
	fmt.Printf("Original: Could you send me the report by tomorrow morning?\n")
	fmt.Printf("German (formal): %s\n\n", *instructedResult.Translation.String)

	// Example 5: Auto-detecting source language
	fmt.Println("=== Auto-detect Source Language ===")
	autoResult, err := laraTranslator.Translate("Bonjour le monde!", "", "en-US", lara.TranslateOptions{})
	if err != nil {
		log.Printf("Error with auto-detection: %v", err)
		return
	}
	fmt.Printf("Original: Bonjour le monde!\n")
	fmt.Printf("Detected source: %s\n", autoResult.SourceLanguage)
	fmt.Printf("English: %s\n\n", *autoResult.Translation.String)

	// Example 6: Advanced options with comprehensive settings
	fmt.Println("=== Translation with Advanced Options ===")
	advancedResult, err := laraTranslator.Translate(
		"This is a comprehensive translation example",
		"en-US",
		"it-IT",
		lara.TranslateOptions{
			AdaptTo:      []string{"mem_1A2b3C4d5E6f7G8h9I0jKl", "mem_2XyZ9AbC8dEf7GhI6jKlMn"}, // Replace with actual memory IDs
			Glossaries:   []string{"gls_1A2b3C4d5E6f7G8h9I0jKl", "gls_2XyZ9AbC8dEf7GhI6jKlMn"}, // Replace with actual glossary IDs
			Instructions: []string{"Be professional"},
			Style:        lara.TranslationStyleFluid,
			ContentType:  "text/plain",
			TimeoutMs:    10000,
		},
	)
	if err != nil {
		log.Printf("Error with advanced translation: %v", err)
		return
	}
	fmt.Printf("Original: This is a comprehensive translation example\n")
	fmt.Printf("Italian (with all options): %s\n\n", *advancedResult.Translation.String)

	// Example 7: Profanity detection options
	fmt.Println("=== Translation with Profanities Detection and Handling Options ===")
	profanityText := "Don't be such a tool."
	verbose := true
	detectTargetResult, err := laraTranslator.Translate(profanityText, "en-US", "it-IT", lara.TranslateOptions{
		ProfanitiesDetect:   lara.ProfanitiesDetectSourceTarget,
		ProfanitiesHandling: lara.ProfanitiesHandlingDetect,
		Verbose:             &verbose,
	})
	if err != nil {
		log.Printf("Error with profanity detect (target): %v", err)
		return
	}
	hideResult, err := laraTranslator.Translate(profanityText, "en-US", "it-IT", lara.TranslateOptions{
		ProfanitiesDetect:   lara.ProfanitiesDetectTarget,
		ProfanitiesHandling: lara.ProfanitiesHandlingHide,
		Verbose:             &verbose,
	})
	if err != nil {
		log.Printf("Error with profanity hide: %v", err)
		return
	}
	fmt.Printf("Original: %s\n", profanityText)
	fmt.Printf("Detect handling (target): %s\n", *detectTargetResult.Translation.String)
	fmt.Printf("Hide handling (source+target): %s\n\n", *hideResult.Translation.String)

	// Example 8: Get available languages
	fmt.Println("=== Available Languages ===")
	languages, err := laraTranslator.Languages()
	if err != nil {
		log.Printf("Error getting languages: %v", err)
		return
	}
	fmt.Printf("Supported languages: %v\n", languages)

	// Example 8: Quality estimation for a single sentence pair
	fmt.Println("\n=== Quality Estimation: single sentence ===")
	qeSingle, err := laraTranslator.QualityEstimation(
		"en-US",
		"it-IT",
		"Hello, how are you today?",
		"Ciao, come stai oggi?",
	)
	if err != nil {
		log.Printf("Error in quality estimation: %v", err)
		return
	}
	single := qeSingle.(*lara.QualityEstimationResult)
	fmt.Printf("Score: %v\n", single.Score)

	// Example 9: Quality estimation for a batch of sentence pairs
	fmt.Println("\n=== Quality Estimation: batch ===")
	qeBatch, err := laraTranslator.QualityEstimation(
		"en-US",
		"it-IT",
		[]string{"Good morning.", "The weather is nice."},
		[]string{"Buongiorno.", "Il tempo è bello."},
	)
	if err != nil {
		log.Printf("Error in quality estimation: %v", err)
		return
	}
	batch := qeBatch.([]lara.QualityEstimationResult)
	scores := make([]float64, len(batch))
	for i, r := range batch {
		scores[i] = r.Score
	}
	fmt.Printf("Scores: %v\n", scores)
}
