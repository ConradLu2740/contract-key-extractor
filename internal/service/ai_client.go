package service

import (
	"bytes"
	"contract-key-extractor/internal/config"
	"contract-key-extractor/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type AIServiceClient struct {
	baseURL    string
	httpClient *http.Client
	logger     *zap.Logger
}

func NewAIServiceClient(cfg *config.AIServiceConfig, logger *zap.Logger) *AIServiceClient {
	return &AIServiceClient{
		baseURL: fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port),
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
		logger: logger,
	}
}

func (c *AIServiceClient) ExtractContractInfo(doc *model.ParsedDocument) (*model.AIExtractionResponse, error) {
	reqBody := model.AIExtractionRequest{
		DocumentText: doc.Content,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := c.baseURL + "/api/v1/extract"
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI service returned error: %s - %s", resp.Status, string(respBody))
	}

	var result model.AIExtractionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (c *AIServiceClient) PerformOCR(fileData []byte) (string, error) {
	url := c.baseURL + "/api/v1/ocr/raw"
	req, err := http.NewRequest("POST", url, bytes.NewReader(fileData))
	if err != nil {
		return "", fmt.Errorf("failed to create OCR request: %w", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send OCR request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OCR service returned error: %s - %s", resp.Status, string(respBody))
	}

	var result struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode OCR response: %w", err)
	}

	return result.Text, nil
}

func (c *AIServiceClient) PerformPDFOCR(pdfData []byte) (string, error) {
	url := c.baseURL + "/api/v1/ocr/pdf"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "document.pdf")
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = part.Write(pdfData)
	if err != nil {
		return "", fmt.Errorf("failed to write pdf data: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("failed to create PDF OCR request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send PDF OCR request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("PDF OCR service returned error: %s - %s", resp.Status, string(respBody))
	}

	var result struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode PDF OCR response: %w", err)
	}

	return result.Text, nil
}

func (c *AIServiceClient) HealthCheck() error {
	url := c.baseURL + "/health"
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("AI service health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("AI service unhealthy: %s", resp.Status)
	}

	return nil
}
