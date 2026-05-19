package lara

import (
	"fmt"
	"os"
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
	err := m.client.Get("/v2/memories", nil, nil, &memories)
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
	err := m.client.Post("/v2/memories", body, nil, nil, &memory)
	if err != nil {
		return nil, fmt.Errorf("failed to create memory: %w", err)
	}

	return &memory, nil
}

func (m *MemoriesService) Get(id string) (*Memory, error) {
	var memory Memory
	err := m.client.Get(fmt.Sprintf("/v2/memories/%s", id), nil, nil, &memory)
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
	err := m.client.Delete(fmt.Sprintf("/v2/memories/%s", id), nil, nil, &memory)
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
	err := m.client.Put(fmt.Sprintf("/v2/memories/%s", id), body, nil, nil, &memory)
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
	err := m.client.Post("/v2/memories/connect", body, nil, nil, &memories)
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
	return m.ImportTmxFromPathWithCallback(id, tmxPath, false, "")
}

func (m *MemoriesService) ImportTmxFromPathWithCallback(id, tmxPath string, gzip bool, callbackUrl string) (*MemoryImport, error) {
	file, err := os.Open(tmxPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return m.ImportTmxWithCallback(id, file, gzip, callbackUrl)
}

func (m *MemoriesService) ImportTmx(id string, tmx *os.File) (*MemoryImport, error) {
	return m.ImportTmxWithCallback(id, tmx, false, "")
}

func (m *MemoriesService) ImportTmxWithCallback(id string, tmx *os.File, gzip bool, callbackUrl string) (*MemoryImport, error) {
	body := map[string]string{}
	if gzip {
		body["compression"] = "gzip"
	}
	if callbackUrl != "" {
		body["callback_url"] = callbackUrl
	}

	files := map[string]*os.File{
		"tmx": tmx,
	}

	var memoryImport MemoryImport
	err := m.client.Post(fmt.Sprintf("/v2/memories/%s/import", id), body, files, nil, &memoryImport)
	if err != nil {
		return nil, fmt.Errorf("failed to import TMX: %w", err)
	}

	return &memoryImport, nil
}

func (m *MemoriesService) GetImportStatus(id string) (*MemoryImport, error) {
	var memoryImport MemoryImport
	err := m.client.Get(fmt.Sprintf("/v2/memories/imports/%s", id), nil, nil, &memoryImport)
	if err != nil {
		return nil, fmt.Errorf("failed to get import status: %w", err)
	}

	return &memoryImport, nil
}

func (m *MemoriesService) ExportAsync(id, callbackUrl string) (*MemoryExport, error) {
	return m.ExportAsyncWithFormat(id, callbackUrl, "")
}

func (m *MemoriesService) ExportAsyncWithFormat(id, callbackUrl string, format MemoryExportFormat) (*MemoryExport, error) {
	params := map[string]string{"callback_url": callbackUrl}
	if format != "" {
		params["format"] = string(format)
	}

	var memoryExport MemoryExport
	err := m.client.Get(fmt.Sprintf("/v2/memories/%s/export/async", id), params, nil, &memoryExport)
	if err != nil {
		return nil, fmt.Errorf("failed to start memory export: %w", err)
	}

	return &memoryExport, nil
}

func (m *MemoriesService) AddTranslation(id, source, target, sentence, translation string) (*MemoryImport, error) {
	return m.AddTranslationWithContextAndHeaders(id, source, target, sentence, translation, "", "", "", nil)
}

func (m *MemoriesService) AddTranslationWithHeaders(id, source, target, sentence, translation string, headers map[string]string) (*MemoryImport, error) {
	return m.AddTranslationWithContextAndHeaders(id, source, target, sentence, translation, "", "", "", headers)
}

func (m *MemoriesService) AddTranslationWithTuid(id, source, target, sentence, translation, tuid string) (*MemoryImport, error) {
	return m.AddTranslationWithContextAndHeaders(id, source, target, sentence, translation, tuid, "", "", nil)
}

func (m *MemoriesService) AddTranslationWithTuidAndHeaders(id, source, target, sentence, translation, tuid string, headers map[string]string) (*MemoryImport, error) {
	return m.AddTranslationWithContextAndHeaders(id, source, target, sentence, translation, tuid, "", "", headers)
}

func (m *MemoriesService) AddTranslationWithContext(id, source, target, sentence, translation, tuid, sentenceBefore, sentenceAfter string) (*MemoryImport, error) {
	return m.AddTranslationWithContextAndHeaders(id, source, target, sentence, translation, tuid, sentenceBefore, sentenceAfter, nil)
}

func (m *MemoriesService) AddTranslationWithContextAndHeaders(id, source, target, sentence, translation, tuid, sentenceBefore, sentenceAfter string, headers map[string]string) (*MemoryImport, error) {
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
	err := m.client.Put(fmt.Sprintf("/v2/memories/%s/content", id), body, nil, headers, &memoryImport)
	if err != nil {
		return nil, fmt.Errorf("failed to add translation: %w", err)
	}

	return &memoryImport, nil
}

func (m *MemoriesService) AddTranslationMultiple(ids []string, source, target, sentence, translation string) (*MemoryImport, error) {
	return m.AddTranslationMultipleWithContextAndHeaders(ids, source, target, sentence, translation, "", "", "", nil)
}

func (m *MemoriesService) AddTranslationMultipleWithHeaders(ids []string, source, target, sentence, translation string, headers map[string]string) (*MemoryImport, error) {
	return m.AddTranslationMultipleWithContextAndHeaders(ids, source, target, sentence, translation, "", "", "", headers)
}

func (m *MemoriesService) AddTranslationMultipleWithTuid(ids []string, source, target, sentence, translation, tuid string) (*MemoryImport, error) {
	return m.AddTranslationMultipleWithContextAndHeaders(ids, source, target, sentence, translation, tuid, "", "", nil)
}

func (m *MemoriesService) AddTranslationMultipleWithTuidAndHeaders(ids []string, source, target, sentence, translation, tuid string, headers map[string]string) (*MemoryImport, error) {
	return m.AddTranslationMultipleWithContextAndHeaders(ids, source, target, sentence, translation, tuid, "", "", headers)
}

func (m *MemoriesService) AddTranslationMultipleWithContext(ids []string, source, target, sentence, translation, tuid, sentenceBefore, sentenceAfter string) (*MemoryImport, error) {
	return m.AddTranslationMultipleWithContextAndHeaders(ids, source, target, sentence, translation, tuid, sentenceBefore, sentenceAfter, nil)
}

func (m *MemoriesService) AddTranslationMultipleWithContextAndHeaders(ids []string, source, target, sentence, translation, tuid, sentenceBefore, sentenceAfter string, headers map[string]string) (*MemoryImport, error) {
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
	err := m.client.Put("/v2/memories/content", body, nil, headers, &memoryImport)
	if err != nil {
		return nil, fmt.Errorf("failed to add multiple translations: %w", err)
	}

	return &memoryImport, nil
}

func (m *MemoriesService) DeleteTranslation(id, source, target, sentence, translation string) (*MemoryImport, error) {
	return m.DeleteTranslationWithContext(id, source, target, sentence, translation, "", "", "")
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
	err := m.client.Delete(fmt.Sprintf("/v2/memories/%s/content", id), body, nil, &memoryImport)
	if err != nil {
		return nil, fmt.Errorf("failed to delete translation: %w", err)
	}

	return &memoryImport, nil
}

func (m *MemoriesService) DeleteTranslationMultiple(ids []string, source, target, sentence, translation string) (*MemoryImport, error) {
	return m.DeleteTranslationMultipleWithContext(ids, source, target, sentence, translation, "", "", "")
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
	err := m.client.Delete("/v2/memories/content", body, nil, &memoryImport)
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
