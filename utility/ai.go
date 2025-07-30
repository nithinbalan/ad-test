package utility

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// GetAIEndpoint builds the full URL from environment base and relative path
func GetAIEndpoint(path string) (string, error) {
	baseURL := os.Getenv("AI_SERVICE_BASE_URL")
	if baseURL == "" {
		return "", fmt.Errorf("AI_SERVICE_BASE_URL is not set")
	}
	return fmt.Sprintf("%s%s", baseURL, path), nil
}

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

func ParseDocument(filePath string) (*ParsedDocumentResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

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

	apiURL, err := GetAIEndpoint("/api/v1/documents-parsing")
	if err != nil {
		return nil, err
	}

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

	var response DocumentParserResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	sampleQuestionsEnJSON, _ := json.Marshal(response.SampleQuestionsEn)
	sampleQuestionsArJSON, _ := json.Marshal(response.SampleQuestionsAr)
	keyPointsEnJSON, _ := json.Marshal(response.KeyPointsEn)
	keyPointsArJSON, _ := json.Marshal(response.KeyPointsAr)
	faqEnJSON, _ := json.Marshal(response.FAQEn)
	faqArJSON, _ := json.Marshal(response.FAQAr)

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

func CallAIAPI(path string, requestBody interface{}) (map[string]interface{}, error) {
	apiURL, err := GetAIEndpoint(path)
	if err != nil {
		return nil, err
	}

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

// func CallAIChatAPI(path string, payload map[string]interface{}) (map[string]interface{}, error) {
// 	apiURL, err := GetAIEndpoint(path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	jsonData, err := json.Marshal(payload)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to marshal payload: %w", err)
// 	}

// 	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create request: %w", err)
// 	}

// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to make request: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("AI service returned non-200 status code: %d", resp.StatusCode)
// 	}

// 	var result map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		return nil, fmt.Errorf("failed to decode response: %w", err)
// 	}

// 	return result, nil
// }

func StreamAIAPI(ctx context.Context, path string, payload map[string]interface{}) (<-chan string, error) {
	apiURL, err := GetAIEndpoint(path)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("AI service returned non-200 status code: %d", resp.StatusCode)
	}

	streamChan := make(chan string)

	go func() {
		defer resp.Body.Close()
		defer close(streamChan)

		reader := bufio.NewReader(resp.Body)
		for {
			select {
			case <-ctx.Done():
				log.Println("⛔ Context cancelled inside stream reader")
				return
			default:
				line, err := reader.ReadString('\n')
				if err != nil {
					if err != io.EOF {
						log.Printf("❌ Read error: %v\n", err)
						streamChan <- fmt.Sprintf(`{"error": "%s"}`, err.Error())
					}
					return
				}

				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}

				log.Printf("📌 Raw line: %s", line)
				streamChan <- line
			}
		}
	}()

	return streamChan, nil
}

func CallAIAPIObjectResponse(path string) (map[string]interface{}, error) {
	apiURL, err := GetAIEndpoint(path)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned non-200 status: %d - %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding API response: %w", err)
	}

	return result, nil
}

func CallAIAPIArrayResponse(path string) ([]map[string]interface{}, error) {
	apiURL, err := GetAIEndpoint(path)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned non-200 status: %d - %s", resp.StatusCode, string(body))
	}

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding API response: %w", err)
	}

	return result, nil
}
