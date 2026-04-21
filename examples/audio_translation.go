package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/translated/lara-go/lara"
)

/**
 * Complete audio translation examples for the Lara Go SDK
 *
 * Supported audio formats: .wav, .mp3, .opus, .ogg, .webm
 *
 * This example demonstrates:
 * - Basic audio translation
 * - Advanced options with memories and glossaries
 * - Step-by-step audio translation with status monitoring
 */

func main() {
	// All examples use environment variables for credentials, so set them first:
	// export LARA_ACCESS_KEY_ID="your-access-key-id"
	// export LARA_ACCESS_KEY_SECRET="your-access-key-secret"

	// Set your credentials here
	accessKeyID := os.Getenv("LARA_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("LARA_ACCESS_KEY_SECRET")

	credentials := lara.NewAccessKey(accessKeyID, accessKeySecret)
	translator := lara.NewTranslator(credentials, nil)

	// Replace with your actual audio file path
	sampleFilePath := "sample_audio.mp3" // Supported: .wav, .mp3, .opus, .ogg, .webm
	filename := "sample_audio.mp3"

	if _, err := os.Stat(sampleFilePath); os.IsNotExist(err) {
		fmt.Printf("Please create a sample audio file at: %s\n", sampleFilePath)
		return
	}

	sourceLang := "en-US"
	targetLang := "de-DE"

	// Example 1: Basic audio translation
	fmt.Println("=== Basic Audio Translation ===")
	fmt.Printf("Translating audio: %s from %s to %s\n", filename, sourceLang, targetLang)

	reader, err := translator.Audio.Translate(&sampleFilePath, &filename, &sourceLang, targetLang)
	if err != nil {
		log.Printf("Error translating audio: %v\n", err)
		return
	}
	defer reader.Close()

	outputPath := "sample_audio_translated.mp3"
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Printf("Error creating output file: %v\n", err)
		return
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, reader)
	if err != nil {
		log.Printf("Error saving translated audio: %v\n", err)
		return
	}

	fmt.Println("Audio translation completed")
	fmt.Printf("Translated file saved to: %s\n\n", outputPath)

	// Example 2: Audio translation with advanced options
	fmt.Println("=== Audio Translation with Advanced Options ===")

	reader2, err := translator.Audio.TranslateWithOptions(&sampleFilePath, &filename, &sourceLang, targetLang, &lara.AudioUploadOptions{
		AdaptTo:    []string{"mem_1A2b3C4d5E6f7G8h9I0jKl"}, // Replace with actual memory IDs
		Glossaries: []string{"gls_1A2b3C4d5E6f7G8h9I0jKl"}, // Replace with actual glossary IDs
	})
	if err != nil {
		log.Printf("Error in advanced translation: %v\n", err)
		return
	}
	defer reader2.Close()

	outputPath2 := "advanced_audio_translated.mp3"
	outputFile2, err := os.Create(outputPath2)
	if err != nil {
		log.Printf("Error creating output file: %v\n", err)
		return
	}
	defer outputFile2.Close()

	_, err = io.Copy(outputFile2, reader2)
	if err != nil {
		log.Printf("Error saving translated audio: %v\n", err)
		return
	}

	fmt.Println("Advanced Audio translation completed")
	fmt.Printf("Translated file saved to: %s\n\n", outputPath2)

	// Example 3: Step-by-step audio translation
	fmt.Println("=== Step-by-Step Audio Translation ===")

	// Step 1: Upload audio
	fmt.Println("Step 1: Uploading audio...")
	audio, err := translator.Audio.UploadWithOptions(&sampleFilePath, &filename, &sourceLang, targetLang, &lara.AudioUploadOptions{
		AdaptTo:    []string{"mem_1A2b3C4d5E6f7G8h9I0jKl"}, // Replace with actual memory IDs
		Glossaries: []string{"gls_1A2b3C4d5E6f7G8h9I0jKl"}, // Replace with actual glossary IDs
	})
	if err != nil {
		log.Printf("Upload error: %v\n", err)
		return
	}
	fmt.Printf("Audio uploaded with ID: %s\n", audio.ID)
	fmt.Printf("Initial status: %s\n", audio.Status)

	// Step 2: Check status
	fmt.Println("\nStep 2: Checking status...")
	updatedAudio, err := translator.Audio.Status(audio.ID)
	if err != nil {
		log.Printf("Status error: %v\n", err)
		return
	}
	fmt.Printf("Current status: %s\n", updatedAudio.Status)

	// Wait for completion
	for updatedAudio.Status != lara.AudioStatusTranslated && updatedAudio.Status != lara.AudioStatusError {
		time.Sleep(2 * time.Second)
		updatedAudio, err = translator.Audio.Status(audio.ID)
		if err != nil {
			log.Printf("Status error: %v\n", err)
			return
		}
	}
	log.Printf("Final status: %s\n", updatedAudio.Status)

	if updatedAudio.Status == lara.AudioStatusError {
		log.Printf("Audio translation failed: %s\n", *updatedAudio.ErrorReason)
		return
	}

	// Step 3: Download translated audio
	fmt.Println("\nStep 3: Downloading translated audio...")
	reader3, err := translator.Audio.Download(audio.ID)
	if err != nil {
		log.Printf("Download error: %v\n", err)
		return
	}
	defer reader3.Close()

	outputPath3 := "step_audio_translated.mp3"
	outputFile3, err := os.Create(outputPath3)
	if err != nil {
		log.Printf("Error creating output file: %v\n", err)
		return
	}
	defer outputFile3.Close()

	_, err = io.Copy(outputFile3, reader3)
	if err != nil {
		log.Printf("Error saving translated audio: %v\n", err)
		return
	}

	fmt.Println("Step-by-step translation completed")
}
