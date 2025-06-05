package utility

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type FAQEntry struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type DocumentParserResponse struct {
	DocumentContent   string     `json:"document_content"`
	SampleQuestionsEn []string   `json:"sample_questions_en"`
	SampleQuestionsAr []string   `json:"sample_questions_ar"`
	SummaryEn         string     `json:"summary_en"`
	SummaryAr         string     `json:"summary_ar"`
	KeyPointsEn       []string   `json:"key_points_en"`
	KeyPointsAr       []string   `json:"key_points_ar"`
	FAQEn             []FAQEntry `json:"faq_en"`
	FAQAr             []FAQEntry `json:"faq_ar"`
}

type ParsedDocumentResult struct {
	ParsedContent     string
	SampleQuestionsEn string
	SampleQuestionsAr string
	SummaryEn         string
	SummaryAr         string
	KeyPointsEn       string
	KeyPointsAr       string
	FAQEn             string
	FAQAr             string
}

// ParseDocument sends the document to the AI parsing service and returns structured content.
func ParseDocument(filePath string) (*ParsedDocumentResult, error) {
	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %w", err)
	}
	if _, err = io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("error copying file content: %w", err)
	}
	if err = writer.Close(); err != nil {
		return nil, fmt.Errorf("error closing multipart writer: %w", err)
	}

	// Prepare API call
	baseURL := os.Getenv("AI_SERVICE_BASE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("AI_SERVICE_BASE_URL is not set")
	}
	apiURL := fmt.Sprintf("%s/api/v1/documents-parsing", baseURL)

	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode response
	var response DocumentParserResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	// Marshal array fields
	sampleQuestionsEnJSON, _ := json.Marshal(response.SampleQuestionsEn)
	sampleQuestionsArJSON, _ := json.Marshal(response.SampleQuestionsAr)
	keyPointsEnJSON, _ := json.Marshal(response.KeyPointsEn)
	keyPointsArJSON, _ := json.Marshal(response.KeyPointsAr)
	faqEnJSON, _ := json.Marshal(response.FAQEn)
	faqArJSON, _ := json.Marshal(response.FAQAr)

	// Return result
	return &ParsedDocumentResult{
		ParsedContent:     response.DocumentContent,
		SampleQuestionsEn: string(sampleQuestionsEnJSON),
		SampleQuestionsAr: string(sampleQuestionsArJSON),
		SummaryEn:         response.SummaryEn,
		SummaryAr:         response.SummaryAr,
		KeyPointsEn:       string(keyPointsEnJSON),
		KeyPointsAr:       string(keyPointsArJSON),
		FAQEn:             string(faqEnJSON),
		FAQAr:             string(faqArJSON),
	}, nil
}

// CallAIAPI sends a generic JSON POST request and decodes the response into a map.
func CallAIAPI(apiURL string, requestBody interface{}) (map[string]interface{}, error) {
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return result, nil
}
