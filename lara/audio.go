package lara

import (
	"fmt"
	"io"
	"time"
)

type AudioTranslator struct {
	client   *Client
	s3Client *S3Client
}

func newAudioTranslator(client *Client, s3Client *S3Client) *AudioTranslator {
	return &AudioTranslator{
		client:   client,
		s3Client: s3Client,
	}
}

// Upload uploads an audio file and creates a translation job
func (a *AudioTranslator) Upload(filePath, filename, source *string, target string) (*Audio, error) {
	return a.UploadWithOptions(filePath, filename, source, target, nil)
}

// UploadWithOptions uploads an audio file with advanced options
func (a *AudioTranslator) UploadWithOptions(filePath, filename, source *string, target string, options *AudioUploadOptions) (*Audio, error) {
	params := map[string]string{
		"filename": *filename,
	}

	var uploadResponse struct {
		URL    string         `json:"url"`
		Fields s3UploadFields `json:"fields"`
	}
	err := a.client.Get("/v2/audio/upload-url", params, nil, &uploadResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to get upload URL: %w", err)
	}

	err = a.s3Client.Upload(uploadResponse.URL, uploadResponse.Fields, *filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to S3: %w", err)
	}

	body := map[string]interface{}{
		"s3key":  uploadResponse.Fields["key"],
		"target": target,
	}

	if source != nil {
		body["source"] = *source
	}

	if options != nil {
		if len(options.AdaptTo) > 0 {
			body["adapt_to"] = options.AdaptTo
		}
		if len(options.Glossaries) > 0 {
			body["glossaries"] = options.Glossaries
		}
		if options.Style != "" {
			body["style"] = options.Style
		}
		if options.VoiceCloning != nil {
			body["voice_cloning"] = *options.VoiceCloning
		}
		if options.VoiceGender != "" {
			body["voice_gender"] = options.VoiceGender
		}
	}

	var headers map[string]string
	if options != nil && options.NoTrace != nil && *options.NoTrace {
		headers = map[string]string{"X-No-Trace": "true"}
	}

	var audio Audio
	err = a.client.Post("/v2/audio/translate", body, nil, headers, &audio)
	if err != nil {
		return nil, fmt.Errorf("failed to create audio translation: %w", err)
	}

	return &audio, nil
}

// Status retrieves the current status of an audio translation job
func (a *AudioTranslator) Status(id string) (*Audio, error) {
	var audio Audio
	err := a.client.Get(fmt.Sprintf("/v2/audio/%s", id), nil, nil, &audio)
	if err != nil {
		return nil, fmt.Errorf("failed to get audio status: %w", err)
	}
	return &audio, nil
}

// Download retrieves the translated audio file
func (a *AudioTranslator) Download(id string) (io.ReadCloser, error) {
	var downloadResponse struct {
		URL string `json:"url"`
	}
	err := a.client.Get(fmt.Sprintf("/v2/audio/%s/download-url", id), nil, nil, &downloadResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to get download URL: %w", err)
	}

	return a.s3Client.Download(downloadResponse.URL)
}

// Translate performs a complete translation workflow: upload, wait, and download
func (a *AudioTranslator) Translate(filePath, filename, source *string, target string) (io.ReadCloser, error) {
	return a.TranslateWithOptions(filePath, filename, source, target, nil)
}

// TranslateWithOptions performs a complete translation workflow with options
func (a *AudioTranslator) TranslateWithOptions(filePath, filename, source *string, target string, options *AudioUploadOptions) (io.ReadCloser, error) {
	audio, err := a.UploadWithOptions(filePath, filename, source, target, options)
	if err != nil {
		return nil, err
	}

	audio, err = a.waitForCompletion(audio)
	if err != nil {
		return nil, err
	}

	if audio.Status == AudioStatusError {
		errorMsg := "translation failed"
		if audio.ErrorReason != nil {
			errorMsg = *audio.ErrorReason
		}
		return nil, fmt.Errorf("audio translation failed: %s", errorMsg)
	}

	return a.Download(audio.ID)
}

// UploadForTranscription uploads an audio file and creates a transcript translation job
func (a *AudioTranslator) UploadForTranscription(filePath, filename, source *string, target string) (*Audio, error) {
	return a.UploadForTranscriptionWithOptions(filePath, filename, source, target, nil)
}

// UploadForTranscriptionWithOptions uploads an audio file for transcript translation with advanced options
func (a *AudioTranslator) UploadForTranscriptionWithOptions(filePath, filename, source *string, target string, options *AudioTranscriptUploadOptions) (*Audio, error) {
	params := map[string]string{
		"filename": *filename,
	}

	var uploadResponse struct {
		URL    string         `json:"url"`
		Fields s3UploadFields `json:"fields"`
	}
	err := a.client.Get("/v2/audio/upload-url", params, nil, &uploadResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to get upload URL: %w", err)
	}

	err = a.s3Client.Upload(uploadResponse.URL, uploadResponse.Fields, *filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to S3: %w", err)
	}

	body := map[string]interface{}{
		"s3key":  uploadResponse.Fields["key"],
		"target": target,
	}

	if source != nil {
		body["source"] = *source
	}

	if options != nil {
		if len(options.AdaptTo) > 0 {
			body["adapt_to"] = options.AdaptTo
		}
		if len(options.Glossaries) > 0 {
			body["glossaries"] = options.Glossaries
		}
		if options.Style != "" {
			body["style"] = options.Style
		}
	}

	var headers map[string]string
	if options != nil && options.NoTrace != nil && *options.NoTrace {
		headers = map[string]string{"X-No-Trace": "true"}
	}

	var audio Audio
	err = a.client.Post("/v2/audio/translate-transcript", body, nil, headers, &audio)
	if err != nil {
		return nil, fmt.Errorf("failed to create transcript translation: %w", err)
	}

	return &audio, nil
}

// GetTranslatedTranscript retrieves the translated transcript for a completed job
func (a *AudioTranslator) GetTranslatedTranscript(id string) (*AudioTextResult, error) {
	var result AudioTextResult
	err := a.client.Get(fmt.Sprintf("/v2/audio/%s/translated-transcript", id), nil, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get translated transcript: %w", err)
	}
	return &result, nil
}

// TranslateTranscript performs a complete transcript workflow: upload, wait, and retrieve
func (a *AudioTranslator) TranslateTranscript(filePath, filename, source *string, target string) (*AudioTextResult, error) {
	return a.TranslateTranscriptWithOptions(filePath, filename, source, target, nil)
}

// TranslateTranscriptWithOptions performs a complete transcript workflow with options
func (a *AudioTranslator) TranslateTranscriptWithOptions(filePath, filename, source *string, target string, options *AudioTranscriptUploadOptions) (*AudioTextResult, error) {
	audio, err := a.UploadForTranscriptionWithOptions(filePath, filename, source, target, options)
	if err != nil {
		return nil, err
	}

	audio, err = a.waitForCompletion(audio)
	if err != nil {
		return nil, err
	}

	if audio.Status == AudioStatusError {
		errorMsg := "translation failed"
		if audio.ErrorReason != nil {
			errorMsg = *audio.ErrorReason
		}
		return nil, fmt.Errorf("audio transcript translation failed: %s", errorMsg)
	}

	return a.GetTranslatedTranscript(audio.ID)
}

// waitForCompletion polls the API until the translation is complete
func (a *AudioTranslator) waitForCompletion(audio *Audio) (*Audio, error) {
	pollingInterval := 2 * time.Second
	maxWaitTime := 15 * time.Minute

	start := time.Now()
	current := audio

	for current.Status != AudioStatusTranslated && current.Status != AudioStatusError {
		if time.Since(start) > maxWaitTime {
			return current, fmt.Errorf("timeout waiting for translation to complete")
		}

		time.Sleep(pollingInterval)

		updated, err := a.Status(current.ID)
		if err != nil {
			return current, err
		}

		current = updated
	}

	return current, nil
}
