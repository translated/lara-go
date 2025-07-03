package lara

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type GlossariesService struct {
	client          *Client
	pollingInterval time.Duration
}

func newGlossariesService(client *Client) *GlossariesService {
	return &GlossariesService{
		client:          client,
		pollingInterval: 2 * time.Second,
	}
}

func (g *GlossariesService) List() ([]Glossary, error) {
	var glossaries []Glossary
	err := g.client.Get("/glossaries", nil, nil, &glossaries)
	if err != nil {
		return nil, fmt.Errorf("failed to list glossaries: %w", err)
	}
	return glossaries, nil
}

func (g *GlossariesService) Create(name string) (*Glossary, error) {
	body := map[string]interface{}{
		"name": name,
	}

	var glossary Glossary
	err := g.client.Post("/glossaries", body, nil, nil, &glossary)
	if err != nil {
		return nil, fmt.Errorf("failed to create glossary: %w", err)
	}
	return &glossary, nil
}

func (g *GlossariesService) Get(id string) (*Glossary, error) {
	var glossary Glossary
	err := g.client.Get(fmt.Sprintf("/glossaries/%s", id), nil, nil, &glossary)
	if err != nil {
		if laraErr, ok := err.(*LaraError); ok && laraErr.Status == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get glossary: %w", err)
	}
	return &glossary, nil
}

func (g *GlossariesService) Delete(id string) (*Glossary, error) {
	var glossary Glossary
	err := g.client.Delete(fmt.Sprintf("/glossaries/%s", id), nil, nil, &glossary)
	if err != nil {
		return nil, fmt.Errorf("failed to delete glossary: %w", err)
	}
	return &glossary, nil
}

func (g *GlossariesService) Update(id, name string) (*Glossary, error) {
	body := map[string]interface{}{
		"name": name,
	}

	var glossary Glossary
	err := g.client.Put(fmt.Sprintf("/glossaries/%s", id), body, nil, nil, &glossary)
	if err != nil {
		return nil, fmt.Errorf("failed to update glossary: %w", err)
	}
	return &glossary, nil
}

func (g *GlossariesService) ImportCsvFromPath(id string, csvPath string) (*GlossaryImport, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	return g.ImportCsv(id, file)
}

func (g *GlossariesService) ImportCsv(id string, csv *os.File) (*GlossaryImport, error) {
	// Auto-detect gzip compression based on filename (like Java SDK)
	fileName := csv.Name()
	isGzipped := strings.HasSuffix(strings.ToLower(fileName), ".gz")

	body := map[string]interface{}{}
	if isGzipped {
		body["compression"] = "gzip"
	}

	files := map[string]io.Reader{
		"csv": csv,
	}

	var glossaryImport GlossaryImport
	err := g.client.Post(fmt.Sprintf("/glossaries/%s/import", id), body, files, nil, &glossaryImport)
	if err != nil {
		return nil, fmt.Errorf("failed to import CSV to glossary: %w", err)
	}
	return &glossaryImport, nil
}

func (g *GlossariesService) GetImportStatus(id string) (*GlossaryImport, error) {
	var glossaryImport GlossaryImport
	err := g.client.Get(fmt.Sprintf("/glossaries/imports/%s", id), nil, nil, &glossaryImport)
	if err != nil {
		return nil, fmt.Errorf("failed to get glossary import status: %w", err)
	}
	return &glossaryImport, nil
}

func (g *GlossariesService) Counts(id string) (*GlossaryCounts, error) {
	var counts GlossaryCounts
	err := g.client.Get(fmt.Sprintf("/glossaries/%s/counts", id), nil, nil, &counts)
	if err != nil {
		return nil, fmt.Errorf("failed to get glossary counts: %w", err)
	}
	return &counts, nil
}

func (g *GlossariesService) Export(id string, contentType string, source *string) ([]byte, error) {
	params := map[string]string{
		"content_type": contentType,
	}
	if source != nil {
		params["source"] = *source
	}

	content, err := g.client.GetRaw(fmt.Sprintf("/glossaries/%s/export", id), params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to export glossary: %w", err)
	}
	return content, nil
}

func (g *GlossariesService) WaitForImport(glossaryImport *GlossaryImport, updateCallback func(*GlossaryImport), maxWaitTime *time.Duration) (*GlossaryImport, error) {
	start := time.Now()
	current := *glossaryImport

	for current.Progress < 1.0 {
		if maxWaitTime != nil && time.Since(start) > *maxWaitTime {
			return &current, fmt.Errorf("timeout waiting for glossary import to complete")
		}

		time.Sleep(g.pollingInterval)

		updated, err := g.GetImportStatus(current.ID)
		if err != nil {
			return &current, fmt.Errorf("failed to get import status: %w", err)
		}

		current = *updated

		if updateCallback != nil {
			updateCallback(&current)
		}
	}

	return &current, nil
}
