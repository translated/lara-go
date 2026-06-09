# Lara Go SDK

[![Go Version](https://img.shields.io/badge/go-1.19+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

This SDK empowers you to build your own branded translation AI leveraging our translation fine-tuned language model.

All major translation features are accessible, making it easy to integrate and customize for your needs.

## 🌍 **Features:**
- **Text Translation**: Single strings, multiple strings, and complex text blocks
- **Document Translation**: Word, PDF, and other document formats with status monitoring
- **Translation Memory**: Store and reuse translations for consistency
- **Glossaries**: Enforce terminology standards across translations
- **Styleguides**: Define tone, voice, and writing style rules for translations
- **Audio Translation**: Translate audio files with status monitoring
- **Image Translation**: Translate whole images or extract and translate text blocks
- **Language Detection**: Automatic source language identification
- **Advanced Options**: Translation instructions and more

## 📚 Documentation

Lara's SDK full documentation is available at [https://developers.laratranslate.com/](https://developers.laratranslate.com/)

## 🚀 Quick Start

### Installation

```bash
go get github.com/translated/lara-go
```

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/translated/lara-go/lara"
)

func main() {
    // Set your credentials using environment variables (recommended)
    accessKeyID := os.Getenv("LARA_ACCESS_KEY_ID")
    accessKeySecret := os.Getenv("LARA_ACCESS_KEY_SECRET")

    // Create translator instance
    laraTranslator := lara.NewTranslator(lara.NewCredentials(accessKeyID, accessKeySecret), nil)

    // Simple text translation
    result, err := laraTranslator.Translate("Hello, world!", "en-US", "fr-FR", lara.TranslateOptions{})
    if err != nil {
        log.Printf("Translation error: %v", err)
        return
    }

    fmt.Printf("Translation: %s\n", *result.Translation.String)
    // Output: Translation: Bonjour, le monde !
}
```

## 📖 Examples

The `examples/` directory contains comprehensive examples for all SDK features.

**All examples use environment variables for credentials, so set them first:**
```bash
export LARA_ACCESS_KEY_ID="your-access-key-id"
export LARA_ACCESS_KEY_SECRET="your-access-key-secret"
```

### Text Translation
- **[text_translation.go](examples/text_translation.go)** - Complete text translation examples
  - Single string translation
  - Multiple strings translation  
  - Translation with instructions
  - TextBlocks translation (mixed translatable/non-translatable content)
  - Auto-detect source language
  - Advanced translation options
  - Get available languages

```bash
cd examples
go run text_translation.go
```

### Document Translation
- **[document_translation.go](examples/document_translation.go)** - Document translation examples
  - Basic document translation
  - Advanced options with memories and glossaries
  - Step-by-step document translation with status monitoring

```bash
cd examples
go run document_translation.go
```

### Translation Memory Management
- **[memories_management.go](examples/memories_management.go)** - Memory management examples
  - Create, list, update, delete memories
  - Add individual translations
  - Multiple memory operations
  - TMX file import with progress monitoring
  - TMX import with callback URL (async notification)
  - Async memory export with callback URL
  - Translation deletion
  - Translation with TUID and context

```bash
cd examples
go run memories_management.go
```

### Glossary Management
- **[glossaries_management.go](examples/glossaries_management.go)** - Glossary management examples
  - Create, list, update, delete glossaries
  - CSV import with status monitoring
  - Glossary export (sync and async)
  - Glossary terms count
  - Import status checking

```bash
cd examples
go run glossaries_management.go
```

### Styleguide Management
- **[styleguides_management.go](examples/styleguides_management.go)** - Styleguide management examples
  - Create, list, get, update, delete styleguides
  - Update name, content, or both at once
  - Handling of non-existent styleguides

```bash
cd examples
go run styleguides_management.go
```

### Audio Translation
- **[audio_translation.go](examples/audio_translation.go)** - Audio translation examples
  - Basic audio translation
  - Advanced options with memories and glossaries
  - Step-by-step audio translation with status monitoring

```bash
cd examples
go run audio_translation.go
```

### Image Translation
- **[image_translation.go](examples/image_translation.go)** - Image translation examples
  - Basic image translation
  - Advanced options with memories and glossaries
  - Extract and translate text from an image

```bash
cd examples
go run image_translation.go
```

### Language Detection
- **[language_detection.go](examples/language_detection.go)** - Language detection examples
  - Detect language from a single string
  - Detect language from multiple strings
  - Detect with hint parameter
  - Detect with passlist (allowlist) parameter

```bash
cd examples
go run language_detection.go
```

## 🔧 API Reference

### Core Components

### 🔐 Authentication

The SDK supports authentication via access key and secret:

```go
credentials := lara.NewCredentials("your-access-key-id", "your-access-key-secret")
laraTranslator := lara.NewTranslator(credentials, nil)
```

**Environment Variables (Recommended):**
```bash
export LARA_ACCESS_KEY_ID="your-access-key-id"
export LARA_ACCESS_KEY_SECRET="your-access-key-secret"
```

```go
accessKeyID := os.Getenv("LARA_ACCESS_KEY_ID")
accessKeySecret := os.Getenv("LARA_ACCESS_KEY_SECRET")
credentials := lara.NewCredentials(accessKeyID, accessKeySecret)
```

### 🌍 Translator

```go
// Create translator with credentials
laraTranslator := lara.NewTranslator(credentials, nil)
```

#### Text Translation

```go
// Basic translation
result, err := laraTranslator.Translate("Hello", "en-US", "fr-FR", lara.TranslateOptions{})

// Multiple strings
result, err := laraTranslator.Translate([]string{"Hello", "World"}, "en-US", "fr-FR", lara.TranslateOptions{})

// TextBlocks (mixed translatable/non-translatable content)
textBlocks := []lara.TextBlock{
    {Text: "Translatable text", Translatable: true},
    {Text: "<br>", Translatable: false},  // Non-translatable HTML
    {Text: "More translatable text", Translatable: true}
}
result, err := laraTranslator.Translate(textBlocks, "en-US", "fr-FR", lara.TranslateOptions{})

// With advanced options
result, err := laraTranslator.Translate("Hello", "en-US", "fr-FR", lara.TranslateOptions{
    Instructions: []string{"Formal tone"},
    AdaptTo:      []string{"mem_1A2b3C4d5E6f7G8h9I0jKl"},  // Replace with actual memory IDs
    Glossaries:   []string{"gls_1A2b3C4d5E6f7G8h9I0jKl"},  // Replace with actual glossary IDs
    Style:        lara.TranslationStyleFluid,
    TimeoutMs:    10000,
})
```

#### Language Detection

```go
// Basic detection
result, err := laraTranslator.Detect("Bonjour le monde!", "", nil)

// Multiple strings
result, err := laraTranslator.Detect([]string{"Hello", "World"}, "", nil)

// With hint
result, err := laraTranslator.Detect("Ciao mondo!", "it-IT", nil)

// With passlist (allowlist)
result, err := laraTranslator.Detect("Hola mundo!", "", []string{"en-US", "fr-FR", "es-ES"})

// With both hint and passlist
result, err := laraTranslator.Detect("Guten Tag!", "de-DE", []string{"de-DE", "en-US", "fr-FR"})

// Access results
fmt.Printf("Language: %s\n", result.Language)
fmt.Printf("Content Type: %s\n", result.ContentType)
```

#### Quality Estimation

Use `QualityEstimation()` to score how well a translation matches its source. Pass a single sentence/translation pair to get a single result, or two parallel arrays to get one result per pair.

```go
// Single pair
result, err := laraTranslator.QualityEstimation(
    "en-US",
    "it-IT",
    "Hello, how are you today?",
    "Ciao, come stai oggi?",
)
single := result.(*lara.QualityEstimationResult)
fmt.Printf("Score: %v\n", single.Score) // e.g. 0.768

// Batch
result, err = laraTranslator.QualityEstimation(
    "en-US",
    "it-IT",
    []string{"Good morning.", "The weather is nice."},
    []string{"Buongiorno.", "Il tempo è bello."},
)
batch := result.([]lara.QualityEstimationResult)
for _, r := range batch {
    fmt.Printf("Score: %v\n", r.Score) // e.g. 0.751, 0.713
}
```

### 📖 Document Translation
#### Simple document translation

```go
filePath := "/path/to/your/document.txt"  // Replace with actual file path
filename := "document.txt"
source := "en-US"
target := "fr-FR"
reader, err := laraTranslator.Documents.Translate(&filePath, &filename, &source, target)

// With options
options := &lara.DocumentTranslateOptions{
    DocumentUploadOptions: lara.DocumentUploadOptions{
        AdaptTo:    []string{"mem_1A2b3C4d5E6f7G8h9I0jKl"},  // Replace with actual memory IDs
        Glossaries: []string{"gls_1A2b3C4d5E6f7G8h9I0jKl"},  // Replace with actual glossary IDs
        Style:      lara.TranslationStyleFluid,
    },
}
reader, err := laraTranslator.Documents.TranslateWithOptions(&filePath, &filename, &source, target, options)
```
### Document translation with status monitoring
#### Document upload
```go
//Optional: upload options
uploadOptions := &lara.DocumentUploadOptions{
    AdaptTo:    []string{"mem_1A2b3C4d5E6f7G8h9I0jKl"},  // Replace with actual memory IDs
    Glossaries: []string{"gls_1A2b3C4d5E6f7G8h9I0jKl"}  // Replace with actual glossary IDs
}

document, err := laraTranslator.Documents.UploadWithOptions(&filePath, &filename, &source, target, uploadOptions)
```
#### Document translation status monitoring
```go
status, err := laraTranslator.Documents.Status(document.ID)
```
#### Download translated document
```go
reader, err := laraTranslator.Documents.Download(document.ID)
```

### 🎵 Audio Translation
#### Simple audio translation

```go
filePath := "/path/to/your/audio.mp3"
filename := "audio.mp3"
source := "en-US"
target := "fr-FR"
reader, err := laraTranslator.Audio.Translate(&filePath, &filename, &source, target)
```

#### Audio translation with options

```go
options := &lara.AudioUploadOptions{
    AdaptTo:    []string{"mem_1A2b3C4d5E6f7G8h9I0jKl"},  // Replace with actual memory IDs
    Glossaries: []string{"gls_1A2b3C4d5E6f7G8h9I0jKl"},  // Replace with actual glossary IDs
}
reader, err := laraTranslator.Audio.TranslateWithOptions(&filePath, &filename, &source, target, options)
```

#### Step-by-step audio translation

```go
// Upload
audio, err := laraTranslator.Audio.Upload(&filePath, &filename, &source, target)

// Check status
status, err := laraTranslator.Audio.Status(audio.ID)

// Download translated audio
reader, err := laraTranslator.Audio.Download(audio.ID)
```

### 🖼️ Image Translation

```go
filePath := "/path/to/your/image.png"  // Replace with actual file path
source := "en-US"
target := "fr-FR"

// Translate image and receive the translated image bytes
imageBytes, err := laraTranslator.Images.Translate(&filePath, &source, target)

// Extract and translate text blocks from an image
result, err := laraTranslator.Images.TranslateText(&filePath, &source, target)
```

### 🧠 Memory Management

```go
// Create memory
memory, err := laraTranslator.Memories.Create("MyMemory")

// Create memory with external ID (MyMemory integration)
memory, err := laraTranslator.Memories.CreateWithExternalID("Memory from MyMemory", "aabb1122")  // Replace with actual external ID

// Important: To update/overwrite a translation unit you must provide a tuid. Calls without a tuid always create a new unit and will not update existing entries.
// Add translation to single memory
memoryImport, err := laraTranslator.Memories.AddTranslation("mem_1A2b3C4d5E6f7G8h9I0jKl", "en-US", "fr-FR", "Hello", "Bonjour")

// Add translation to multiple memories
memoryImport, err := laraTranslator.Memories.AddTranslationMultiple([]string{"mem_1A2b3C4d5E6f7G8h9I0jKl", "mem_2XyZ9AbC8dEf7GhI6jKlMn"}, "en-US", "fr-FR", "Hello", "Bonjour")

// Add with context
memoryImport, err := laraTranslator.Memories.AddTranslationWithContext(
    "mem_1A2b3C4d5E6f7G8h9I0jKl", "en-US", "fr-FR", "Hello", "Bonjour", "greeting_003",
    "sentenceBefore", "sentenceAfter"
)

// TMX import from file
tmxFilePath := "/path/to/your/memory.tmx"  // Replace with actual TMX file path
memoryImport, err := laraTranslator.Memories.ImportTmxFromPath("mem_1A2b3C4d5E6f7G8h9I0jKl", tmxFilePath)

// TMX import with a callback URL (notified when the import completes)
memoryImport, err = laraTranslator.Memories.ImportTmxFromPathWithCallback(
    "mem_1A2b3C4d5E6f7G8h9I0jKl",
    tmxFilePath,
    "https://your-server.example.com/lara/import-callback",
)

// Async memory export - returns a job ID; the result is delivered to your callback URL when ready
exportJob, err := laraTranslator.Memories.ExportAsync(
    "mem_1A2b3C4d5E6f7G8h9I0jKl",
    "https://your-server.example.com/lara/export-callback",
)
jobID := exportJob.JobID

// Async memory export with explicit format
exportJob, err = laraTranslator.Memories.ExportAsyncWithFormat(
    "mem_1A2b3C4d5E6f7G8h9I0jKl",
    "https://your-server.example.com/lara/export-callback",
    lara.MemoryExportFormatTmx, // or lara.MemoryExportFormatJtm
)

// Delete translation
// Important: if you omit tuid, all entries that match the provided fields will be removed
deleteJob, err := laraTranslator.Memories.DeleteTranslation(
        "mem_1A2b3C4d5E6f7G8h9I0jKl", "en-US", "fr-FR", "Hello", "Bonjour"
)

// Wait for import completion
import "time"
maxWaitTime := 300 * time.Second // 5 minutes
completedImport, err := laraTranslator.Memories.WaitForImport(memoryImport, nil, &maxWaitTime)
```

### 📚 Glossary Management

```go
// Create glossary
glossary, err := laraTranslator.Glossaries.Create("MyGlossary")

// Import CSV from file
csvFilePath := "/path/to/your/glossary.csv"  // Replace with actual CSV file path
glossaryImport, err := laraTranslator.Glossaries.ImportCsvFromPath("gls_1A2b3C4d5E6f7G8h9I0jKl", csvFilePath)

// Check import status
importStatus, err := laraTranslator.Glossaries.GetImportStatus("gls_1A2b3C4d5E6f7G8h9I0jKl")

// Wait for import completion
import "time"
maxWaitTime := 300 * time.Second // 5 minutes
completedImport, err := laraTranslator.Glossaries.WaitForImport(glossaryImport, nil, &maxWaitTime)

// Export glossary
csvData, err := laraTranslator.Glossaries.Export("gls_1A2b3C4d5E6f7G8h9I0jKl", "csv/table-uni", "en-US")

// Async glossary export - returns a job ID; the result is delivered to your callback URL when ready
exportJob, err := laraTranslator.Glossaries.ExportAsync(
    "gls_1A2b3C4d5E6f7G8h9I0jKl",
    "https://your-server.example.com/lara/export-callback",
    "csv/table-uni",
    "en-US",
)
jobID := exportJob.JobID

// Get glossary terms count
counts, err := laraTranslator.Glossaries.Counts("gls_1A2b3C4d5E6f7G8h9I0jKl")
```

### 📋 Styleguide Management

```go
// Create styleguide
styleguide, err := laraTranslator.Styleguides.Create("MyStyleguide", "Always use formal language.")

// List all styleguides
styleguides, err := laraTranslator.Styleguides.List()

// Get a specific styleguide
styleguide, err = laraTranslator.Styleguides.Get("stg_1A2b3C4d5E6f7G8h9I0jKl")

// Update styleguide — pass nil for fields you don't want to change
// Update only the name
name := "UpdatedStyleguide"
styleguide, err = laraTranslator.Styleguides.Update("stg_1A2b3C4d5E6f7G8h9I0jKl", &name, nil)

// Update only the content
content := "Always use informal language."
styleguide, err = laraTranslator.Styleguides.Update("stg_1A2b3C4d5E6f7G8h9I0jKl", nil, &content)

// Update both
styleguide, err = laraTranslator.Styleguides.Update("stg_1A2b3C4d5E6f7G8h9I0jKl", &name, &content)

// Delete styleguide
styleguide, err = laraTranslator.Styleguides.Delete("stg_1A2b3C4d5E6f7G8h9I0jKl")
```

### Translation Options

```go
type TranslateOptions struct {
    AdaptTo      []string                   // Memory IDs to adapt to
    Glossaries   []string                   // Glossary IDs to use
    Instructions []string                   // Translation instructions
    Style        TranslationStyle           // Translation style (fluid, faithful, creative)
    ContentType  string                     // Content type (text/plain, text/html, etc.)
    Multiline    *bool                      // Enable multiline translation
    TimeoutMs    int                        // Request timeout in milliseconds
    SourceHint   string                     // Hint for source language detection
    NoTrace      *bool                      // Disable request tracing
    Verbose      *bool                      // Enable verbose response
}
```

### Language Codes

The SDK supports full language codes (e.g., `en-US`, `fr-FR`, `es-ES`) as well as simple codes (e.g., `en`, `fr`, `es`):

```go
// Full language codes (recommended)
result, err := laraTranslator.Translate("Hello", "en-US", "fr-FR", lara.TranslateOptions{})

// Simple language codes
result, err := laraTranslator.Translate("Hello", "en", "fr", lara.TranslateOptions{})
```

### 🌐 Supported Languages

The SDK supports all languages available in the Lara API. Use the `Languages()` method to get the current list:

```go
languages, err := laraTranslator.Languages()
fmt.Printf("Supported languages: %v\n", languages)
```

## ⚙️ Configuration

### Error Handling

The SDK provides detailed error information:

```go
result, err := laraTranslator.Translate("Hello", "en-US", "fr-FR", lara.TranslateOptions{})
if err != nil {
    if laraErr, ok := err.(*lara.LaraError); ok {
        fmt.Printf("API Error [%d]: %s\n", laraErr.Status, laraErr.Message)
        fmt.Printf("Error type: %s\n", laraErr.Type)
    } else {
        fmt.Printf("SDK Error: %v\n", err)
    }
    return
}
```

## 📋 Requirements

- Go 1.19 or higher
- Valid Lara API credentials

## 🧪 Testing

Run the examples to test your setup:

```bash
# All examples use environment variables for credentials, so set them first:
export LARA_ACCESS_KEY_ID="your-access-key-id"
export LARA_ACCESS_KEY_SECRET="your-access-key-secret"
```

```bash
# Run basic text translation example
cd examples
go run text_translation.go
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Happy translating! 🌍✨