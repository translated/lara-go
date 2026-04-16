package lara

import (
	"encoding/json"
	"fmt"
)

type Translator struct {
	client     *Client
	Documents  *DocumentsService
	Memories   *MemoriesService
	Glossaries *GlossariesService
}

// NewTranslator creates a new Translator with any supported authentication method.
// Accepts *Credentials (deprecated), *AccessKey, *UserCredentials, or *AuthToken.
func NewTranslator(auth interface{}, options *TranslatorOptions) *Translator {
	if auth == nil {
		auth = NewCredentials("", "")
	}

	serverURL := "https://api.laratranslate.com"
	if options != nil && options.ServerURL != "" {
		serverURL = options.ServerURL
	}

	client := newClient(auth, serverURL)

	s3Client := newS3Client()

	return &Translator{
		client:     client,
		Documents:  &DocumentsService{client: client, s3Client: s3Client},
		Memories:   newMemoriesService(client),
		Glossaries: newGlossariesService(client),
	}
}

type TranslateOptions struct {
	AdaptTo      []string
	Glossaries   []string
	Instructions []string
	ContentType  string
	Multiline    *bool
	TimeoutMs    int
	Priority     string
	UseCache     *bool
	CacheTTL     *int
	SourceHint   string
	NoTrace      *bool
	Verbose      *bool
	Style        TranslationStyle
	Reasoning    *bool
	Headers      map[string]interface{}
	Callback     func(*TextResult) error
}

type DetectOptions struct {
	Hint     string
	Passlist []string
}

type Translation struct {
	String     *string
	Strings    []string
	TextBlocks []TextBlock
}

// Auto-called by Go's json package during unmarshaling
func (t *Translation) UnmarshalJSON(data []byte) error {
	var singleString string
	if err := json.Unmarshal(data, &singleString); err == nil {
		t.String = &singleString
		return nil
	}

	var stringSlice []string
	if err := json.Unmarshal(data, &stringSlice); err == nil {
		t.Strings = stringSlice
		return nil
	}

	var textBlocks []TextBlock
	if err := json.Unmarshal(data, &textBlocks); err == nil {
		t.TextBlocks = textBlocks
		return nil
	}
	return fmt.Errorf("translation: unsupported data type")
}

type TextResult struct {
	ContentType       string              `json:"content_type"`
	SourceLanguage    string              `json:"source_language"`
	Translation       Translation         `json:"translation"` // Custom UnmarshalJSON method is automatically called when unmarshaling this field
	AdaptedTo         []string            `json:"adapted_to,omitempty"`
	Glossaries        []string            `json:"glossaries,omitempty"`
	AdaptedToMatches  [][]NGMemoryMatch   `json:"adapted_to_matches,omitempty"`
	GlossariesMatches [][]NGGlossaryMatch `json:"glossaries_matches,omitempty"`
}

func (t *Translator) Translate(text interface{}, source string, target string, opts TranslateOptions) (*TextResult, error) {
	body := make(map[string]interface{})
	// Accept string, []string, or []TextBlock for text
	switch v := text.(type) {
	case string:
		body["q"] = v
	case []string:
		body["q"] = v
	case []TextBlock:
		body["q"] = v
	default:
		return nil, fmt.Errorf("text must be string, []string, or []TextBlock")
	}

	body["target"] = target

	if source != "" {
		body["source"] = source
	}
	if opts.SourceHint != "" {
		body["source_hint"] = opts.SourceHint
	}
	if len(opts.AdaptTo) > 0 {
		body["adapt_to"] = opts.AdaptTo
	}
	if len(opts.Glossaries) > 0 {
		body["glossaries"] = opts.Glossaries
	}
	if len(opts.Instructions) > 0 {
		body["instructions"] = opts.Instructions
	}
	if opts.ContentType != "" {
		body["content_type"] = opts.ContentType
	}
	if opts.Multiline != nil {
		body["multiline"] = *opts.Multiline
	}
	if opts.TimeoutMs > 0 {
		body["timeout"] = opts.TimeoutMs
	}
	if opts.Priority != "" {
		body["priority"] = opts.Priority
	}
	if opts.UseCache != nil {
		body["use_cache"] = *opts.UseCache
	}
	if opts.CacheTTL != nil {
		body["cache_ttl"] = *opts.CacheTTL
	}
	if opts.Verbose != nil {
		body["verbose"] = *opts.Verbose
	}
	if opts.Style != "" {
		body["style"] = opts.Style
	}
	if opts.Reasoning != nil {
		body["reasoning"] = *opts.Reasoning
	}

	headers := make(map[string]string)
	if opts.Headers != nil {
		for name, value := range opts.Headers {
			if value != nil {
				headers[name] = fmt.Sprint(value)
			}
		}
	}

	if opts.NoTrace != nil && *opts.NoTrace {
		headers["X-No-Trace"] = "true"
	}

	// Always use streaming; callback only invoked when reasoning is enabled
	var lastResult *TextResult
	err := t.client.PostAndGetStream("/translate", body, headers, func(contentBytes []byte) error {
		var result TextResult
		if err := json.Unmarshal(contentBytes, &result); err != nil {
			return fmt.Errorf("failed to unmarshal partial result: %w", err)
		}
		lastResult = &result

		// Only invoke callback if reasoning is enabled
		if opts.Callback != nil && opts.Reasoning != nil && *opts.Reasoning {
			if err := opts.Callback(&result); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to translate text: %w", err)
	}

	if lastResult == nil {
		return nil, fmt.Errorf("no translation result received")
	}

	return lastResult, nil
}

func (t *Translator) Languages() ([]string, error) {
	var languages []string
	err := t.client.Get("/v2/languages", nil, nil, &languages)
	if err != nil {
		return nil, fmt.Errorf("failed to get languages: %w", err)
	}

	return languages, nil
}

func (t *Translator) Detect(text interface{}, hint string, passlist []string) (*DetectResult, error) {
	body := map[string]interface{}{
		"q": text,
	}

	if hint != "" {
		body["hint"] = hint
	}

	if len(passlist) > 0 {
		body["passlist"] = passlist
	}

	var result DetectResult
	err := t.client.Post("/v2/detect", body, nil, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to detect language: %w", err)
	}

	return &result, nil
}
