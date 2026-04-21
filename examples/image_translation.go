package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/translated/lara-go/lara"
)

/**
 * Complete image translation examples for the Lara Go SDK
 *
 * This example demonstrates:
 * - Basic image translation
 * - Advanced options with memories and glossaries
 * - Extracting and translating text from an image
 */

func main() {
	// All examples use environment variables for credentials, so set them first:
	// export LARA_ACCESS_KEY_ID="your-access-key-id"
	// export LARA_ACCESS_KEY_SECRET="your-access-key-secret"

	accessKeyID := os.Getenv("LARA_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("LARA_ACCESS_KEY_SECRET")

	laraTranslator := lara.NewTranslator(lara.NewCredentials(accessKeyID, accessKeySecret), nil)

	// Replace with your actual image file path
	sampleFilePath := filepath.Join(".", "sample_image.png")

	if _, err := os.Stat(sampleFilePath); os.IsNotExist(err) {
		fmt.Printf("Please create a sample image file at: %s\n", sampleFilePath)
		return
	}

	sourceLang := "en"
	targetLang := "de"

	// Example 1: Basic image translation (image output)
	fmt.Println("=== Basic Image Translation ===")
	fmt.Printf("Translating image: %s from %s to %s\n", filepath.Base(sampleFilePath), sourceLang, targetLang)

	imageBytes, err := laraTranslator.Images.Translate(&sampleFilePath, &sourceLang, targetLang)
	if err != nil {
		log.Printf("Error translating image: %v", err)
		return
	}

	outputPath := filepath.Join(".", "sample_image_translated.png")
	if err := os.WriteFile(outputPath, imageBytes, 0644); err != nil {
		log.Printf("Error saving translated image: %v", err)
		return
	}
	fmt.Println("Image translation completed")
	fmt.Printf("Translated image saved to: %s\n\n", filepath.Base(outputPath))

	// Example 2: Image translation with advanced options
	fmt.Println("=== Image Translation with Advanced Options ===")

	advancedBytes, err := laraTranslator.Images.TranslateWithOptions(&sampleFilePath, &sourceLang, targetLang, &lara.ImageTranslateOptions{
		AdaptTo:     []string{"mem_1A2b3C4d5E6f7G8h9I0jKl"},     // Replace with actual memory IDs
		Glossaries:  []string{"gls_1A2b3C4d5E6f7G8h9I0jKl"},     // Replace with actual glossary IDs
		Style:       lara.TranslationStyleFaithful,
		TextRemoval: lara.TextRemovalInpainting,
	})
	if err != nil {
		log.Printf("Error in advanced translation: %v", err)
		return
	}

	advancedOutputPath := filepath.Join(".", "advanced_image_translated.png")
	if err := os.WriteFile(advancedOutputPath, advancedBytes, 0644); err != nil {
		log.Printf("Error saving translated image: %v", err)
		return
	}
	fmt.Println("Advanced image translation completed")
	fmt.Printf("Translated image saved to: %s\n\n", filepath.Base(advancedOutputPath))

	// Example 3: Extract and translate text from an image
	fmt.Println("=== Extract and Translate Text ===")

	results, err := laraTranslator.Images.TranslateText(&sampleFilePath, &sourceLang, targetLang)
	if err != nil {
		log.Printf("Error extracting and translating text: %v", err)
		return
	}

	fmt.Println("Extract and translate completed")
	fmt.Printf("Found %d text blocks\n", len(results.Paragraphs))

	for i, result := range results.Paragraphs {
		fmt.Printf("\nText Block %d:\n", i+1)
		fmt.Printf("Original: %s\n", result.Text)
		fmt.Printf("Translated: %s\n", result.Translation)
	}
}
