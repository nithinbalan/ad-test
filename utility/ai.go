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

// FAQEntry represents a single question-answer pair in the AI response
type FAQEntry struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// DocumentParserResponse represents the structure of the AI document parsing API response
type DocumentParserResponse struct {
	DocumentContent string     `json:"document_content"`
	SampleQuestions []string   `json:"sample_questions"`
	Summary         string     `json:"summary"`
	KeyPoints       []string   `json:"key_points"`
	FAQ             []FAQEntry `json:"faq"`
}

// ParseDocument sends the document to the AI parsing service and returns structured content.
func ParseDocument(filePath string) (string, string, string, string, string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Create multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("error creating form file: %w", err)
	}
	if _, err = io.Copy(part, file); err != nil {
		return "", "", "", "", "", fmt.Errorf("error copying file content: %w", err)
	}
	if err = writer.Close(); err != nil {
		return "", "", "", "", "", fmt.Errorf("error closing multipart writer: %w", err)
	}

	// Get AI service base URL from environment
	baseURL := os.Getenv("AI_SERVICE_BASE_URL")
	if baseURL == "" {
		return "", "", "", "", "", fmt.Errorf("AI_SERVICE_BASE_URL is not set")
	}
	apiURL := fmt.Sprintf("%s/api/v1/documents-parsing", baseURL)

	// Send request
	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", "", "", "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode response
	var parserResponse DocumentParserResponse
	if err := json.NewDecoder(resp.Body).Decode(&parserResponse); err != nil {
		return "", "", "", "", "", fmt.Errorf("error decoding response: %w", err)
	}

	// Marshal complex fields as JSON strings
	sampleQuestionsJSON, _ := json.Marshal(parserResponse.SampleQuestions)
	keyPointsJSON, _ := json.Marshal(parserResponse.KeyPoints)
	faqJSON, _ := json.Marshal(parserResponse.FAQ)

	return parserResponse.DocumentContent,
		string(sampleQuestionsJSON),
		parserResponse.Summary,
		string(keyPointsJSON),
		string(faqJSON),
		nil
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
