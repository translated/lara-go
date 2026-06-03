package lara

import (
	"fmt"
	"os"
)

type ImagesService struct {
	client *Client
}

func newImagesService(client *Client) *ImagesService {
	return &ImagesService{
		client: client,
	}
}

func (s *ImagesService) Translate(filePath, source *string, target string) ([]byte, error) {
	return s.TranslateWithOptions(filePath, source, target, nil)
}

func (s *ImagesService) TranslateWithOptions(filePath, source *string, target string, options *ImageTranslateOptions) ([]byte, error) {
	file, err := os.Open(*filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	body := map[string]interface{}{
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
			body["style"] = string(options.Style)
		}
		model := options.Model
		if model == "" {
			model = options.TextRemoval
		}
		if model != "" {
			body["model"] = string(model)
		}
	}

	var headers map[string]string
	if options != nil && options.NoTrace != nil && *options.NoTrace {
		headers = map[string]string{"X-No-Trace": "true"}
	}

	files := map[string]*os.File{
		"image": file,
	}

	result, err := s.client.PostRaw("/v2/images/translate", body, files, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to translate image: %w", err)
	}

	return result, nil
}

func (s *ImagesService) TranslateText(filePath, source *string, target string) (*ImageTextResult, error) {
	return s.TranslateTextWithOptions(filePath, source, target, nil)
}

func (s *ImagesService) TranslateTextWithOptions(filePath, source *string, target string, options *ImageTextTranslateOptions) (*ImageTextResult, error) {
	file, err := os.Open(*filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	body := map[string]interface{}{
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
			body["style"] = string(options.Style)
		}
		if options.Verbose != nil {
			body["verbose"] = *options.Verbose
		}
	}

	var headers map[string]string
	if options != nil && options.NoTrace != nil && *options.NoTrace {
		headers = map[string]string{"X-No-Trace": "true"}
	}

	files := map[string]*os.File{
		"image": file,
	}

	var result ImageTextResult
	err = s.client.Post("/v2/images/translate-text", body, files, headers, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to translate image text: %w", err)
	}

	return &result, nil
}
