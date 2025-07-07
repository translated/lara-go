package lara

import "time"

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

type DocumentDownloadOptions struct {
	OutputFormat string
}

type DocumentUploadOptions struct {
	AdaptTo    []string
	Glossaries []string
	NoTrace    *bool
}

type DocumentOptions struct {
	DocumentUploadOptions
}

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

type TranslatorOptions struct {
	ServerURL string
}
