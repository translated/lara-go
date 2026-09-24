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
		if options.IncludeLayout != nil {
			body["include_layout"] = *options.IncludeLayout
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

// RenderTranslated renders supplied translations onto the original image using generative_fast.
func (s *ImagesService) RenderTranslated(filePath, source *string, target string, paragraphs []ImageParagraph) ([]byte, error) {
	return s.RenderTranslatedWithOptions(filePath, source, target, paragraphs, nil)
}

// RenderTranslatedWithOptions renders supplied translations without translating again.
// Overlay and inpainting require BBox, LinesBBoxes, TextInfo, and Alignment on every
// paragraph. Generative models accept text-only paragraphs or complete layout.
func (s *ImagesService) RenderTranslatedWithOptions(filePath, source *string, target string, paragraphs []ImageParagraph, options *ImageRenderOptions) ([]byte, error) {
	if filePath == nil {
		return nil, fmt.Errorf("image file path is required")
	}
	file, err := os.Open(*filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	// Select only rendering fields; match metadata is not part of this endpoint.
	renderParagraphs := make([]map[string]interface{}, len(paragraphs))
	for i, paragraph := range paragraphs {
		p := map[string]interface{}{"text": paragraph.Text, "translation": paragraph.Translation}
		if paragraph.BBox != nil {
			p["bbox"] = paragraph.BBox
		}
		if paragraph.LinesBBoxes != nil {
			p["lines_bboxes"] = paragraph.LinesBBoxes
		}
		if paragraph.TextInfo != nil {
			p["text_info"] = paragraph.TextInfo
		}
		if paragraph.Alignment != "" {
			p["alignment"] = paragraph.Alignment
		}
		renderParagraphs[i] = p
	}
	body := map[string]interface{}{"target": target, "paragraphs": renderParagraphs}
	if source != nil {
		body["source"] = *source
	}
	var headers map[string]string
	if options != nil {
		if options.Model != "" {
			body["model"] = string(options.Model)
		}
		if options.NoTrace != nil && *options.NoTrace {
			headers = map[string]string{"X-No-Trace": "true"}
		}
	}
	result, err := s.client.PostRaw("/v2/images/render-translated", body, map[string]*os.File{"image": file}, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to render translated image: %w", err)
	}
	return result, nil
}
