package lara

import (
	"fmt"
	"io"
	"time"
)

type DocumentsService struct {
	client   *Client
	s3Client *S3Client
}

type DocumentTranslateOptions struct {
	DocumentUploadOptions
	DocumentDownloadOptions
}

func (d *DocumentsService) Upload(filePath, filename, source *string, target string) (*Document, error) {
	options := &DocumentUploadOptions{}
	return d.UploadWithOptions(filePath, filename, source, target, options)
}

func (d *DocumentsService) UploadWithOptions(filePath, filename, source *string, target string, options *DocumentUploadOptions) (*Document, error) {
	params := map[string]string{
		"filename": *filename,
	}

	var uploadResponse struct {
		URL    string         `json:"url"`
		Fields s3UploadFields `json:"fields"`
	}
	err := d.client.Get("/documents/upload-url", params, nil, &uploadResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to get upload URL: %w", err)
	}

	err = d.s3Client.Upload(uploadResponse.URL, uploadResponse.Fields, *filePath)
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
		headers = map[string]string{
			"X-No-Trace": "true",
		}
	}

	var document Document
	err = d.client.Post("/documents", body, nil, headers, &document)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	return &document, nil
}

func (d *DocumentsService) Status(id string) (*Document, error) {
	var document Document
	err := d.client.Get(fmt.Sprintf("/documents/%s", id), nil, nil, &document)
	if err != nil {
		return nil, fmt.Errorf("failed to get document status: %w", err)
	}

	return &document, nil
}

func (d *DocumentsService) Download(id string) (io.ReadCloser, error) {
	return d.DownloadWithOptions(id, nil)
}

func (d *DocumentsService) DownloadWithOptions(id string, options *DocumentDownloadOptions) (io.ReadCloser, error) {
	params := map[string]string{}
	if options != nil && options.OutputFormat != "" {
		params["output_format"] = options.OutputFormat
	}

	var downloadResponse struct {
		URL string `json:"url"`
	}
	err := d.client.Get(fmt.Sprintf("/documents/%s/download-url", id), params, nil, &downloadResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to get download URL: %w", err)
	}

	return d.s3Client.Download(downloadResponse.URL)
}

func (d *DocumentsService) Translate(filePath, filename, source *string, target string) (io.ReadCloser, error) {
	return d.TranslateWithOptions(filePath, filename, source, target, nil)
}

func (d *DocumentsService) TranslateWithOptions(filePath, filename, source *string, target string, options *DocumentTranslateOptions) (io.ReadCloser, error) {
	uploadOptions := &DocumentUploadOptions{}

	if options != nil {
		uploadOptions.AdaptTo = options.AdaptTo
		uploadOptions.Glossaries = options.Glossaries
		uploadOptions.NoTrace = options.NoTrace
		uploadOptions.Style = options.Style
	}

	document, err := d.UploadWithOptions(filePath, filename, source, target, uploadOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to upload document: %w", err)
	}

	document, err = d.waitForTranslation(document)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for translation: %w", err)
	}

	if document.Status == DocumentStatusError {
		errorMsg := "translation failed"
		if document.ErrorReason != nil {
			errorMsg = *document.ErrorReason
		}
		return nil, fmt.Errorf("document translation failed: %s", errorMsg)
	}

	downloadOptions := &DocumentDownloadOptions{}
	if options != nil && options.OutputFormat != "" {
		downloadOptions.OutputFormat = options.OutputFormat
	}

	return d.DownloadWithOptions(document.ID, downloadOptions)
}

func (d *DocumentsService) waitForTranslation(document *Document) (*Document, error) {
	pollingInterval := 2 * time.Second
	var maxWaitTime time.Duration

	start := time.Now()
	current := document

	for current.Status != DocumentStatusTranslated && current.Status != DocumentStatusError {
		if maxWaitTime > 0 && time.Since(start) > maxWaitTime {
			return current, fmt.Errorf("timeout waiting for translation to complete")
		}

		time.Sleep(pollingInterval)

		updated, err := d.Status(current.ID)
		if err != nil {
			return current, err
		}

		current = updated
	}

	return current, nil
}
