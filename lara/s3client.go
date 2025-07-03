package lara

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type s3UploadFields map[string]string

type S3Client struct {
	httpClient *http.Client
}

func newS3Client() *S3Client {
	return &S3Client{
		httpClient: &http.Client{},
	}
}

func (s *S3Client) Upload(url string, fields s3UploadFields, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for key, value := range fields {
		writer.WriteField(key, value)
	}

	part, err := writer.CreateFormFile("file", fields["key"])
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	writer.Close()

	resp, err := s.httpClient.Post(url, writer.FormDataContentType(), &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("upload failed with status %d", resp.StatusCode)
	}

	return nil
}

func (s *S3Client) Download(url string) (io.ReadCloser, error) {
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	return resp.Body, nil
}
