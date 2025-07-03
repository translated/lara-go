package lara

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type MemoriesService struct {
	client *Client
}

func newMemoriesService(client *Client) *MemoriesService {
	return &MemoriesService{
		client: client,
	}
}

func (m *MemoriesService) List() ([]Memory, error) {
	var memories []Memory
	err := m.client.Get("/memories", nil, nil, &memories)
	if err != nil {
		return nil, fmt.Errorf("failed to list memories: %w", err)
	}
	return memories, nil
}

func (m *MemoriesService) Create(name string) (*Memory, error) {
	return m.CreateWithExternalID(name, "")
}

func (m *MemoriesService) CreateWithExternalID(name string, externalID string) (*Memory, error) {
	body := map[string]interface{}{
		"name": name,
	}
	if externalID != "" {
		body["external_id"] = externalID
	}

	var memory Memory
	err := m.client.Post("/memories", body, nil, nil, &memory)
	if err != nil {
		return nil, fmt.Errorf("failed to create memory: %w", err)
	}

	return &memory, nil
}

func (m *MemoriesService) Get(id string) (*Memory, error) {
	var memory Memory
	err := m.client.Get(fmt.Sprintf("/memories/%s", id), nil, nil, &memory)
	if err != nil {
		if laraErr, ok := err.(*LaraError); ok && laraErr.Status == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get memory: %w", err)
	}

	return &memory, nil
}

func (m *MemoriesService) Delete(id string) (*Memory, error) {
	var memory Memory
	err := m.client.Delete(fmt.Sprintf("/memories/%s", id), nil, nil, &memory)
	if err != nil {
		return nil, fmt.Errorf("failed to delete memory: %w", err)
	}

	return &memory, nil
}

func (m *MemoriesService) Update(id, name string) (*Memory, error) {
	body := map[string]interface{}{
		"name": name,
	}

	var memory Memory
	err := m.client.Put(fmt.Sprintf("/memories/%s", id), body, nil, nil, &memory)
	if err != nil {
		return nil, fmt.Errorf("failed to update memory: %w", err)
	}

	return &memory, nil
}

func (m *MemoriesService) ConnectMultiple(ids []string) ([]Memory, error) {
	body := map[string]interface{}{
		"ids": ids,
	}

	var memories []Memory
	err := m.client.Post("/memories/connect", body, nil, nil, &memories)
	if err != nil {
		return nil, fmt.Errorf("failed to connect memories: %w", err)
	}

	return memories, nil
}

func (m *MemoriesService) Connect(id string) (*Memory, error) {
	memories, err := m.ConnectMultiple([]string{id})
	if err != nil {
		return nil, err
	}
	if len(memories) == 0 {
		return nil, fmt.Errorf("no memory returned for id %s", id)
	}
	return &memories[0], nil
}

func (m *MemoriesService) ImportTmxFromPath(id, tmxPath string) (*MemoryImport, error) {
	file, err := os.Open(tmxPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return m.ImportTmx(id, file)
}

func (m *MemoriesService) ImportTmx(id string, tmx *os.File) (*MemoryImport, error) {
	// Auto-detect gzip compression based on filename (like Java SDK)
	fileName := tmx.Name()
	isGzipped := strings.HasSuffix(strings.ToLower(fileName), ".gz")

	body := map[string]string{}
	if isGzipped {
		body["compression"] = "gzip"
	}

	files := map[string]io.Reader{
		"tmx": tmx,
	}

	var memoryImport MemoryImport
	err := m.client.Post(fmt.Sprintf("/memories/%s/import", id), body, files, nil, &memoryImport)
	if err != nil {
		return nil, fmt.Errorf("failed to import TMX: %w", err)
	}

	return &memoryImport, nil
}

func (m *MemoriesService) GetImportStatus(id string) (*MemoryImport, error) {
	var memoryImport MemoryImport
	err := m.client.Get(fmt.Sprintf("/memories/imports/%s", id), nil, nil, &memoryImport)
	if err != nil {
		return nil, fmt.Errorf("failed to get import status: %w", err)
	}

	return &memoryImport, nil
}

func (m *MemoriesService) AddTranslation(id, source, target, sentence, translation string) (*MemoryImport, error) {
	return m.AddTranslationWithTuid(id, source, target, sentence, translation, "")
}

func (m *MemoriesService) AddTranslationWithTuid(id, source, target, sentence, translation, tuid string) (*MemoryImport, error) {
	return m.AddTranslationWithContext(id, source, target, sentence, translation, tuid, "", "")
}

func (m *MemoriesService) AddTranslationWithContext(id, source, target, sentence, translation, tuid, sentenceBefore, sentenceAfter string) (*MemoryImport, error) {
	body := map[string]interface{}{
		"source":      source,
		"target":      target,
		"sentence":    sentence,
		"translation": translation,
	}
	if tuid != "" {
		body["tuid"] = tuid
	}
	if sentenceBefore != "" {
		body["sentence_before"] = sentenceBefore
	}
	if sentenceAfter != "" {
		body["sentence_after"] = sentenceAfter
	}

	var memoryImport MemoryImport
	err := m.client.Put(fmt.Sprintf("/memories/%s/content", id), body, nil, nil, &memoryImport)
	if err != nil {
		return nil, fmt.Errorf("failed to add translation: %w", err)
	}

	return &memoryImport, nil
}

func (m *MemoriesService) AddTranslationMultiple(ids []string, source, target, sentence, translation string) (*MemoryImport, error) {
	return m.AddTranslationMultipleWithTuid(ids, source, target, sentence, translation, "")
}

func (m *MemoriesService) AddTranslationMultipleWithTuid(ids []string, source, target, sentence, translation, tuid string) (*MemoryImport, error) {
	return m.AddTranslationMultipleWithContext(ids, source, target, sentence, translation, tuid, "", "")
}

func (m *MemoriesService) AddTranslationMultipleWithContext(ids []string, source, target, sentence, translation, tuid, sentenceBefore, sentenceAfter string) (*MemoryImport, error) {
	body := map[string]interface{}{
		"ids":         ids,
		"source":      source,
		"target":      target,
		"sentence":    sentence,
		"translation": translation,
	}
	if tuid != "" {
		body["tuid"] = tuid
	}
	if sentenceBefore != "" {
		body["sentence_before"] = sentenceBefore
	}
	if sentenceAfter != "" {
		body["sentence_after"] = sentenceAfter
	}

	var memoryImport MemoryImport
	err := m.client.Put("/memories/content", body, nil, nil, &memoryImport)
	if err != nil {
		return nil, fmt.Errorf("failed to add multiple translations: %w", err)
	}

	return &memoryImport, nil
}

func (m *MemoriesService) DeleteTranslation(id, source, target, sentence, translation string) (*MemoryImport, error) {
	return m.DeleteTranslationWithTuid(id, source, target, sentence, translation, "")
}

func (m *MemoriesService) DeleteTranslationWithTuid(id, source, target, sentence, translation, tuid string) (*MemoryImport, error) {
	return m.DeleteTranslationWithContext(id, source, target, sentence, translation, tuid, "", "")
}

func (m *MemoriesService) DeleteTranslationWithContext(id, source, target, sentence, translation, tuid, sentenceBefore, sentenceAfter string) (*MemoryImport, error) {
	body := map[string]interface{}{
		"source":      source,
		"target":      target,
		"sentence":    sentence,
		"translation": translation,
	}
	if tuid != "" {
		body["tuid"] = tuid
	}
	if sentenceBefore != "" {
		body["sentence_before"] = sentenceBefore
	}
	if sentenceAfter != "" {
		body["sentence_after"] = sentenceAfter
	}

	var memoryImport MemoryImport
	err := m.client.Delete(fmt.Sprintf("/memories/%s/content", id), body, nil, &memoryImport)
	if err != nil {
		return nil, fmt.Errorf("failed to delete translation: %w", err)
	}

	return &memoryImport, nil
}

func (m *MemoriesService) DeleteTranslationMultiple(ids []string, source, target, sentence, translation string) (*MemoryImport, error) {
	return m.DeleteTranslationMultipleWithTuid(ids, source, target, sentence, translation, "")
}

func (m *MemoriesService) DeleteTranslationMultipleWithTuid(ids []string, source, target, sentence, translation, tuid string) (*MemoryImport, error) {
	return m.DeleteTranslationMultipleWithContext(ids, source, target, sentence, translation, tuid, "", "")
}

func (m *MemoriesService) DeleteTranslationMultipleWithContext(ids []string, source, target, sentence, translation, tuid, sentenceBefore, sentenceAfter string) (*MemoryImport, error) {
	body := map[string]interface{}{
		"ids":         ids,
		"source":      source,
		"target":      target,
		"sentence":    sentence,
		"translation": translation,
	}
	if tuid != "" {
		body["tuid"] = tuid
	}
	if sentenceBefore != "" {
		body["sentence_before"] = sentenceBefore
	}
	if sentenceAfter != "" {
		body["sentence_after"] = sentenceAfter
	}

	var memoryImport MemoryImport
	err := m.client.Delete("/memories/content", body, nil, &memoryImport)
	if err != nil {
		return nil, fmt.Errorf("failed to delete multiple translations: %w", err)
	}

	return &memoryImport, nil
}

func (m *MemoriesService) WaitForImport(memoryImport *MemoryImport, updateCallback func(*MemoryImport), maxWaitTime *time.Duration) (*MemoryImport, error) {
	pollingInterval := 2 * time.Second
	start := time.Now()
	current := *memoryImport

	for current.Progress < 1.0 {
		if maxWaitTime != nil && time.Since(start) > *maxWaitTime {
			return &current, fmt.Errorf("timeout waiting for import to complete")
		}

		time.Sleep(pollingInterval)

		updated, err := m.GetImportStatus(current.ID)
		if err != nil {
			return &current, err
		}

		current = *updated

		if updateCallback != nil {
			updateCallback(&current)
		}
	}

	return &current, nil
}
