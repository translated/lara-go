package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/translated/lara-go/lara"
)

/**
 * Complete audio transcript translation examples for the Lara Go SDK
 *
 * Supported audio formats: .wav, .mp3, .opus, .ogg, .webm
 *
 * This example demonstrates the async Audio2Text flow, which returns only the
 * translated transcript (JSON) instead of a dubbed audio file:
 * - Basic transcript translation
 * - Advanced options with memories and glossaries
 * - Step-by-step transcript translation with status monitoring
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

	// Example 1: Basic transcript translation
	fmt.Println("=== Basic Transcript Translation ===")
	fmt.Printf("Transcribing audio: %s from %s to %s\n", filename, sourceLang, targetLang)

	result, err := translator.Audio.TranslateTranscript(&sampleFilePath, &filename, &sourceLang, targetLang)
	if err != nil {
		log.Printf("Error translating transcript: %v\n", err)
		return
	}

	fmt.Println("Transcript translation completed")
	fmt.Printf("Translation: %s\n", result.Translation)
	fmt.Printf("Segments: %d\n\n", len(result.Segments))

	// Example 2: Transcript translation with advanced options
	fmt.Println("=== Transcript Translation with Advanced Options ===")

	result2, err := translator.Audio.TranslateTranscriptWithOptions(&sampleFilePath, &filename, &sourceLang, targetLang, &lara.AudioTranscriptUploadOptions{
		AdaptTo:    []string{"mem_1A2b3C4d5E6f7G8h9I0jKl"}, // Replace with actual memory IDs
		Glossaries: []string{"gls_1A2b3C4d5E6f7G8h9I0jKl"}, // Replace with actual glossary IDs
	})
	if err != nil {
		log.Printf("Error in advanced translation: %v\n", err)
		return
	}

	fmt.Println("Advanced transcript translation completed")
	fmt.Printf("Translation: %s\n\n", result2.Translation)

	// Example 3: Step-by-step transcript translation
	fmt.Println("=== Step-by-Step Transcript Translation ===")

	// Step 1: Upload audio
	fmt.Println("Step 1: Uploading audio...")
	audio, err := translator.Audio.UploadForTranscriptionWithOptions(&sampleFilePath, &filename, &sourceLang, targetLang, &lara.AudioTranscriptUploadOptions{
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
	maxWaitTime := 15 * time.Minute
	start := time.Now()
	for updatedAudio.Status != lara.AudioStatusTranslated && updatedAudio.Status != lara.AudioStatusError {
		if time.Since(start) > maxWaitTime {
			log.Printf("Timed out waiting for the transcript translation\n")
			return
		}
		time.Sleep(2 * time.Second)
		updatedAudio, err = translator.Audio.Status(audio.ID)
		if err != nil {
			log.Printf("Status error: %v\n", err)
			return
		}
	}
	log.Printf("Final status: %s\n", updatedAudio.Status)

	if updatedAudio.Status == lara.AudioStatusError {
		errorMsg := "translation failed"
		if updatedAudio.ErrorReason != nil {
			errorMsg = *updatedAudio.ErrorReason
		}
		log.Printf("Transcript translation failed: %s\n", errorMsg)
		return
	}

	// Step 3: Retrieve translated transcript
	fmt.Println("\nStep 3: Retrieving translated transcript...")
	result3, err := translator.Audio.GetTranslatedTranscript(audio.ID)
	if err != nil {
		log.Printf("Retrieve error: %v\n", err)
		return
	}

	fmt.Printf("Text: %s\n", result3.Text)
	fmt.Printf("Translation: %s\n", result3.Translation)
	for _, segment := range result3.Segments {
		fmt.Printf("[%.0f - %.0f] %s\n", segment.Start, segment.End, segment.Translation)
	}

	fmt.Println("Step-by-step transcript translation completed")
}
