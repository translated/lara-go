package main

import (
	"fmt"
	"io"
	"os"

	"github.com/translated/lara-go/lara"
)

/**
 * Complete document translation examples for the Lara Go SDK
 *
 * This example demonstrates:
 * - Basic document translation
 * - Advanced options with memories and glossaries
 * - Step-by-step document translation with status monitoring
 */

func main() {
	// All examples use environment variables for credentials, so set them first:
	// export LARA_ACCESS_KEY_ID="your-access-key-id"
	// export LARA_ACCESS_KEY_SECRET="your-access-key-secret"

	// Set your credentials here
	accessKeyID := os.Getenv("LARA_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("LARA_ACCESS_KEY_SECRET")

	translator := lara.NewTranslator(lara.NewCredentials(accessKeyID, accessKeySecret), nil)

	// Replace with your actual document file path
	sampleFilePath := "sample_document.docx" // Create this file with your content
	sampleFileName := "sample_document.docx"

	if _, err := os.Stat(sampleFilePath); os.IsNotExist(err) {
		fmt.Printf("Please create a sample document file at: %s\n", sampleFilePath)
		fmt.Println("Add some sample text content to translate.\n")
		return
	}

	// Example 1: Basic document translation
	fmt.Println("=== Basic Document Translation ===")
	sourceLang := "en-US"
	targetLang := "de-DE"

	fmt.Printf("Translating document: %s from %s to %s\n", sampleFileName, sourceLang, targetLang)

	reader, err := translator.Documents.Translate(&sampleFilePath, &sampleFileName, &sourceLang, targetLang)
	if err != nil {
		fmt.Printf("Error translating document: %v\n\n", err)
		return
	}
	defer reader.Close()

	// Save translated document - replace with your desired output path
	outputPath := "sample_document_translated.docx"
	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n\n", err)
		return
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, reader)
	if err != nil {
		fmt.Printf("Error saving translated document: %v\n\n", err)
		return
	}

	fmt.Println("✅ Document translation completed")
	fmt.Printf("📄 Translated file saved to: %s\n\n", outputPath)

	// Example 2: Document translation with advanced options
	fmt.Println("=== Document Translation with Advanced Options ===")
	noTrace := true
	translationOptions := &lara.DocumentTranslateOptions{
		DocumentUploadOptions: lara.DocumentUploadOptions{
			AdaptTo:    []string{"mem_1A2b3C4d5E6f7G8h9I0jKl"}, // Replace with actual memory IDs
			Glossaries: []string{"gls_1A2b3C4d5E6f7G8h9I0jKl"}, // Replace with actual glossary IDs
			Style:      lara.TranslationStyleFluid,
			NoTrace:    &noTrace,
		},
	}

	advancedReader, err := translator.Documents.TranslateWithOptions(&sampleFilePath, &sampleFileName, &sourceLang, targetLang, translationOptions)
	if err != nil {
		fmt.Printf("Error in advanced translation: %v", err)
	} else {
		defer advancedReader.Close()

		// Save translated document - replace with your desired output path
		outputPath2 := "advanced_document_translated.docx"
		outputFile2, err := os.Create(outputPath2)
		if err != nil {
			fmt.Printf("Error creating advanced output file: %v\n", err)
		} else {
			defer outputFile2.Close()
			_, err = io.Copy(outputFile2, advancedReader)
			if err != nil {
				fmt.Printf("Error saving advanced translated document: %v\n", err)
			} else {
				fmt.Println("✅ Advanced document translation completed")
				fmt.Printf("📄 Translated file saved to: %s", outputPath2)
			}
		}
	}
	fmt.Println()

	// Example 3: Step-by-step document translation
	fmt.Println("=== Step-by-Step Document Translation ===")

	// Upload document
	fmt.Println("Step 1: Uploading document...")
	uploadOptions := &lara.DocumentUploadOptions{
		AdaptTo:    []string{"mem_1A2b3C4d5E6f7G8h9I0jKl"}, // Replace with actual memory IDs
		Glossaries: []string{"gls_1A2b3C4d5E6f7G8h9I0jKl"}, // Replace with actual glossary IDs
		Style:      lara.TranslationStyleFluid,
	}

	document, err := translator.Documents.UploadWithOptions(&sampleFilePath, &sampleFileName, &sourceLang, targetLang, uploadOptions)
	if err != nil {
		fmt.Printf("Error in step-by-step process: %v", err)
		return
	}
	fmt.Printf("Document uploaded with ID: %s\n", document.ID)
	fmt.Printf("Initial status: %s\n", document.Status)

	// Check status
	fmt.Println("\nStep 2: Checking status...")
	updatedDocument, err := translator.Documents.Status(document.ID)
	if err != nil {
		fmt.Printf("Error in step-by-step process: %v", err)
		return
	}
	fmt.Printf("Current status: %s\n", updatedDocument.Status)

	// Download translated document
	fmt.Println("\nStep 3: Downloading would happen after translation completes...")
	downloadedContent, err := translator.Documents.Download(document.ID)
	if err != nil {
		fmt.Printf("Error in step-by-step process: %v", err)
		return
	}
	defer downloadedContent.Close()
	fmt.Println("✅ Step-by-step translation completed")
}
