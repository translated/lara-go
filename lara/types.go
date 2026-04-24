package lara

import (
	"encoding/json"
	"fmt"
	"time"
)

type Memory struct {
	ID                 string     `json:"id"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	Name               string     `json:"name"`
	ExternalID         *string    `json:"external_id,omitempty"`
	Secret             *string    `json:"secret,omitempty"`
	OwnerID            string     `json:"owner_id"`
	CollaboratorsCount int        `json:"collaborators_count"`
	SharedAt           *time.Time `json:"shared_at,omitempty"`
}

type Import struct {
	ID       string  `json:"id"`
	Begin    int     `json:"begin"`
	End      int     `json:"end"`
	Channel  int     `json:"channel"`
	Size     int     `json:"size"`
	Progress float64 `json:"progress"`
}

type MemoryImport = Import
type GlossaryImport = Import

type Glossary struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	OwnerID   string    `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GlossaryCounts struct {
	Unidirectional   map[string]int `json:"unidirectional"`
	Multidirectional int            `json:"multidirectional"`
}

type GlossaryFileFormat string

const (
	GlossaryFileFormatCsvTableUni   GlossaryFileFormat = "csv/table-uni"
	GlossaryFileFormatCsvTableMulti GlossaryFileFormat = "csv/table-multi"
)

type GlossaryTerm struct {
	Language string `json:"language"`
	Value    string `json:"value"`
}

type DocumentStatus string

const (
	DocumentStatusInitialized DocumentStatus = "initialized"
	DocumentStatusAnalyzing   DocumentStatus = "analyzing"
	DocumentStatusPaused      DocumentStatus = "paused"
	DocumentStatusReady       DocumentStatus = "ready"
	DocumentStatusTranslating DocumentStatus = "translating"
	DocumentStatusTranslated  DocumentStatus = "translated"
	DocumentStatusError       DocumentStatus = "error"
)

type Document struct {
	ID              string           `json:"id"`
	Status          DocumentStatus   `json:"status"`
	Source          *string          `json:"source,omitempty"`
	Target          string           `json:"target"`
	Filename        string           `json:"filename"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	Options         *DocumentOptions `json:"options,omitempty"`
	TranslatedChars *int             `json:"translated_chars,omitempty"`
	TotalChars      *int             `json:"total_chars,omitempty"`
	ErrorReason     *string          `json:"error_reason,omitempty"`
}

type DocumentOptions struct {
	AdaptTo    []string         `json:"adapt_to,omitempty"`
	Glossaries []string         `json:"glossaries,omitempty"`
	NoTrace    *bool            `json:"no_trace,omitempty"`
	Style      TranslationStyle `json:"style,omitempty"`
}

type DocumentDownloadOptions struct {
	OutputFormat string
}

type DocumentUploadOptions struct {
	DocumentOptions
	ExtractionParams DocumentExtractionParams `json:"extraction_params,omitempty"`
	Password         *string                  `json:"password,omitempty"`
}

type DocumentTranslateOptions struct {
	DocumentUploadOptions
	DocumentDownloadOptions
}

type DocxExtractionParams struct {
	ExtractComments *bool `json:"extract_comments,omitempty"`
	AcceptRevisions *bool `json:"accept_revisions,omitempty"`
}

// Used to handle future extraction parameter types
type DocumentExtractionParams interface {
	extractionParams()
}

func (DocxExtractionParams) extractionParams() {}

type TextBlock struct {
	Text         string `json:"text"`
	Translatable bool   `json:"translatable"`
}

type NGMemoryMatch struct {
	Memory      string    `json:"memory"`
	TUID        *string   `json:"tuid,omitempty"`
	Language    [2]string `json:"language"`
	Sentence    string    `json:"sentence"`
	Translation string    `json:"translation"`
	Score       float64   `json:"score"`
}

type NGGlossaryMatch struct {
	Glossary    string    `json:"glossary"`
	Language    [2]string `json:"language"`
	Term        string    `json:"term"`
	Translation string    `json:"translation"`
}

// NGMemoryMatchGroups is [][]NGMemoryMatch that handles both the flat form
// (single-string translation input) and the nested form (array input).
type NGMemoryMatchGroups [][]NGMemoryMatch

func (g *NGMemoryMatchGroups) UnmarshalJSON(data []byte) error {
	var nested [][]NGMemoryMatch
	if err := json.Unmarshal(data, &nested); err == nil {
		*g = nested
		return nil
	}
	var flat []NGMemoryMatch
	if err := json.Unmarshal(data, &flat); err == nil {
		*g = [][]NGMemoryMatch{flat}
		return nil
	}
	return fmt.Errorf("adapted_to_matches: unsupported data type")
}

// NGGlossaryMatchGroups is [][]NGGlossaryMatch that handles both the flat form
// (single-string translation input) and the nested form (array input).
type NGGlossaryMatchGroups [][]NGGlossaryMatch

func (g *NGGlossaryMatchGroups) UnmarshalJSON(data []byte) error {
	var nested [][]NGGlossaryMatch
	if err := json.Unmarshal(data, &nested); err == nil {
		*g = nested
		return nil
	}
	var flat []NGGlossaryMatch
	if err := json.Unmarshal(data, &flat); err == nil {
		*g = [][]NGGlossaryMatch{flat}
		return nil
	}
	return fmt.Errorf("glossaries_matches: unsupported data type")
}

type Styleguide struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Content   *string   `json:"content,omitempty"`
	OwnerID   string    `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StyleguideChange struct {
	ID                  *string `json:"id,omitempty"`
	OriginalTranslation string  `json:"original_translation"`
	RefinedTranslation  string  `json:"refined_translation"`
	Explanation         string  `json:"explanation"`
}

type StyleguideResults struct {
	OriginalTranslation Translation        `json:"original_translation"`
	Changes             []StyleguideChange `json:"changes"`
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
	Style     TranslationStyle
	Reasoning *bool
	Metadata  interface{}
	StyleguideID                  string
	StyleguideReasoning           *bool
	StyleguideExplanationLanguage string
	Headers   map[string]interface{}
	Callback  func(*TextResult) error
}

type Translation struct {
	String     *string
	Strings    []string
	TextBlocks []TextBlock
}

type TextResult struct {
	ContentType       string              `json:"content_type"`
	SourceLanguage    string              `json:"source_language"`
	Translation       Translation         `json:"translation"`
	AdaptedTo         []string            `json:"adapted_to,omitempty"`
	Glossaries        []string            `json:"glossaries,omitempty"`
	AdaptedToMatches  NGMemoryMatchGroups   `json:"adapted_to_matches,omitempty"`
	GlossariesMatches NGGlossaryMatchGroups `json:"glossaries_matches,omitempty"`
	StyleguideResults *StyleguideResults  `json:"styleguide_results,omitempty"`
}

type TranslationStyle string

const (
	TranslationStyleFaithful TranslationStyle = "faithful"
	TranslationStyleFluid    TranslationStyle = "fluid"
	TranslationStyleCreative TranslationStyle = "creative"
)

type DetectPrediction struct {
	Language   string  `json:"language"`
	Confidence float64 `json:"confidence"`
}

type DetectResult struct {
	Language    string             `json:"language"`
	ContentType string             `json:"content_type"`
	Predictions []DetectPrediction `json:"predictions"`
}

// VoiceGender represents the gender of the voice for audio translation
type VoiceGender string

const (
	VoiceGenderMale   VoiceGender = "male"
	VoiceGenderFemale VoiceGender = "female"
)

// AudioStatus is an alias for DocumentStatus (same statuses used)
type AudioStatus = DocumentStatus

// Audio status constants (aliases for DocumentStatus constants)
const (
	AudioStatusInitialized = DocumentStatusInitialized
	AudioStatusAnalyzing   = DocumentStatusAnalyzing
	AudioStatusPaused      = DocumentStatusPaused
	AudioStatusReady       = DocumentStatusReady
	AudioStatusTranslating = DocumentStatusTranslating
	AudioStatusTranslated  = DocumentStatusTranslated
	AudioStatusError       = DocumentStatusError
)

// Audio represents an audio translation job
type Audio struct {
	ID                string        `json:"id"`
	Status            AudioStatus   `json:"status"`
	Source            *string       `json:"source,omitempty"`
	Target            string        `json:"target"`
	Filename          string        `json:"filename"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
	Options           *AudioOptions `json:"options,omitempty"`
	TranslatedSeconds *float64      `json:"translated_seconds,omitempty"`
	TotalSeconds      *float64      `json:"total_seconds,omitempty"`
	ErrorReason       *string       `json:"error_reason,omitempty"`
}

type AudioOptions struct {
	AdaptTo     []string         `json:"adapt_to,omitempty"`
	Glossaries  []string         `json:"glossaries,omitempty"`
	Style       TranslationStyle `json:"style,omitempty"`
	NoTrace     *bool            `json:"no_trace,omitempty"`
	VoiceGender VoiceGender      `json:"voice_gender,omitempty"`
}

type AudioUploadOptions struct {
	AdaptTo     []string
	Glossaries  []string
	Style       TranslationStyle
	NoTrace     *bool
	VoiceGender VoiceGender
}

type Profanity struct {
	Text           string  `json:"text"`
	StartCharIndex int     `json:"start_char_index"`
	EndCharIndex   int     `json:"end_char_index"`
	Score          float64 `json:"score"`
}

type ProfanityDetectResult struct {
	MaskedText  string      `json:"masked_text"`
	Profanities []Profanity `json:"profanities"`
}

// TextRemoval represents the text removal mode for image translation
type TextRemoval string

const (
	TextRemovalOverlay    TextRemoval = "overlay"
	TextRemovalInpainting TextRemoval = "inpainting"
)

type ImageTranslateOptions struct {
	AdaptTo     []string
	Glossaries  []string
	Style       TranslationStyle
	TextRemoval TextRemoval
	NoTrace     *bool
}

type ImageTextTranslateOptions struct {
	AdaptTo    []string
	Glossaries []string
	Style      TranslationStyle
	NoTrace    *bool
	Verbose    *bool
}

type ImageParagraph struct {
	Text              string            `json:"text"`
	Translation       string            `json:"translation"`
	AdaptedToMatches  []NGMemoryMatch   `json:"adapted_to_matches,omitempty"`
	GlossariesMatches []NGGlossaryMatch `json:"glossaries_matches,omitempty"`
}

type ImageTextResult struct {
	SourceLanguage string           `json:"source_language"`
	AdaptedTo      []string         `json:"adapted_to,omitempty"`
	Glossaries     []string         `json:"glossaries,omitempty"`
	Paragraphs     []ImageParagraph `json:"paragraphs"`
}
