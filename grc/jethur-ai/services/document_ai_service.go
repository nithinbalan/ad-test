package services

import (
	"context"
	"encoding/json"
	"fmt"
	"jethur/grc/documents/services"
	aiDTOs "jethur/grc/jethur-ai/dtos"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"os"
	"strconv"
)

type DocumentAIServiceInterface interface {
	GetAISummary(documentId uint) (*aiDTOs.DocumentAISummaryDTO, error)
	GetDocumentRecommendations(documentId uint, currentUserID uint) (map[string]interface{}, error)
	ProcessChat(ctx context.Context, content, role, documentID, contextID string) (<-chan string, error)
}

type DocumentAIService struct {
	documentService services.DocumentServiceInterface
	userService     organizationServices.UserServiceInterface
}

func NewDocumentAIService(
	documentService services.DocumentServiceInterface,
	userService organizationServices.UserServiceInterface,
) DocumentAIServiceInterface {
	return &DocumentAIService{
		documentService: documentService,
		userService:     userService,
	}
}

func (as *DocumentAIService) GetAISummary(documentId uint) (*aiDTOs.DocumentAISummaryDTO, error) {
	// Get full document details
	document, err := as.documentService.GetFullDetailsByID(documentId, nil)
	if err != nil {
		return nil, err
	}

	if document.LatestDocumentVersion.ID == 0 {
		return nil, fmt.Errorf("no document version found")
	}

	documentVersion, err := as.documentService.GetDocumentVersionByID(document.LatestDocumentVersion.ID)
	if err != nil {
		return nil, err
	}

	// Initialize DTO with empty arrays
	aiSummaryDTO := &aiDTOs.DocumentAISummaryDTO{
		SummaryEn:         documentVersion.SummaryEn,
		SummaryAr:         documentVersion.SummaryAr,
		SampleQuestionsEn: []string{},
		SampleQuestionsAr: []string{},
		KeyPointsEn:       []string{},
		KeyPointsAr:       []string{},
		FAQEn:             []aiDTOs.FAQItemDTO{},
		FAQAr:             []aiDTOs.FAQItemDTO{},
	}

	// Unmarshal EN
	if documentVersion.SampleQuestionsEn != nil {
		if err := json.Unmarshal([]byte(*documentVersion.SampleQuestionsEn), &aiSummaryDTO.SampleQuestionsEn); err != nil {
			return nil, fmt.Errorf("error unmarshalling sample_questions_en: %w", err)
		}
	}
	if documentVersion.KeyPointsEn != nil {
		if err := json.Unmarshal([]byte(*documentVersion.KeyPointsEn), &aiSummaryDTO.KeyPointsEn); err != nil {
			return nil, fmt.Errorf("error unmarshalling key_points_en: %w", err)
		}
	}
	if documentVersion.FAQEn != nil {
		if err := json.Unmarshal([]byte(*documentVersion.FAQEn), &aiSummaryDTO.FAQEn); err != nil {
			return nil, fmt.Errorf("error unmarshalling faq_en: %w", err)
		}
	}

	// Unmarshal AR
	if documentVersion.SampleQuestionsAr != nil {
		if err := json.Unmarshal([]byte(*documentVersion.SampleQuestionsAr), &aiSummaryDTO.SampleQuestionsAr); err != nil {
			return nil, fmt.Errorf("error unmarshalling sample_questions_ar: %w", err)
		}
	}
	if documentVersion.KeyPointsAr != nil {
		if err := json.Unmarshal([]byte(*documentVersion.KeyPointsAr), &aiSummaryDTO.KeyPointsAr); err != nil {
			return nil, fmt.Errorf("error unmarshalling key_points_ar: %w", err)
		}
	}
	if documentVersion.FAQAr != nil {
		if err := json.Unmarshal([]byte(*documentVersion.FAQAr), &aiSummaryDTO.FAQAr); err != nil {
			return nil, fmt.Errorf("error unmarshalling faq_ar: %w", err)
		}
	}

	return aiSummaryDTO, nil
}

func (svc *DocumentAIService) GetDocumentRecommendations(documentId uint, currentUserID uint) (map[string]interface{}, error) {
	userDetails, err := svc.userService.GetCurrentUserDetails(currentUserID)
	if err != nil {
		return nil, err
	}

	designation := ""
	if userDetails.Designation != nil {
		designation = userDetails.Designation.TitleEn
	}

	// Construct request payload
	payload := map[string]interface{}{
		"document_id": documentId,
		"designation": designation,
	}

	baseURL := os.Getenv("AI_SERVICE_BASE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("AI_SERVICE_BASE_URL is not set")
	}
	apiURL := fmt.Sprintf("%s/api/v1/policy/policy-recomendation", baseURL)

	return utility.CallAIAPI(apiURL, payload)
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Response string `json:"response"`
	Error    string `json:"error,omitempty"`
}

func (s *DocumentAIService) ProcessChat(ctx context.Context, content, role, documentID, contextID string) (<-chan string, error) {
	// Convert documentID from string to uint for validation
	docIDInt, err := strconv.ParseUint(documentID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid document ID: %w", err)
	}

	// Get document details
	document, err := s.documentService.GetFullDetailsByID(uint(docIDInt), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	if document.LatestDocumentVersion.ID == 0 {
		return nil, fmt.Errorf("no document version found")
	}

	// Get latest document version
	documentVersion, err := s.documentService.GetDocumentVersionByID(document.LatestDocumentVersion.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document version: %w", err)
	}

	if documentVersion.ParsedContent == nil {
		return nil, fmt.Errorf("document content not available")
	}

	// Create response channel
	responseChan := make(chan string)

	go func() {
		defer close(responseChan)

		// Call AI streaming API
		baseURL := os.Getenv("AI_SERVICE_BASE_URL")
		if baseURL == "" {
			response := ChatResponse{Error: "AI_SERVICE_BASE_URL is not set"}
			jsonResp, _ := json.Marshal(response)
			responseChan <- string(jsonResp)
			return
		}

		apiURL := fmt.Sprintf("%s/api/v1/document-chat", baseURL)

		// Corrected payload to match API expectation
		payload := map[string]interface{}{
			"content":     content,
			"role":        role,
			"document_id": documentID,
			"context_id":  contextID,
		}

		// Call utility for streaming
		streamChan, err := utility.StreamAIAPI(ctx, apiURL, payload)
		if err != nil {
			response := ChatResponse{Error: err.Error()}
			jsonResp, _ := json.Marshal(response)
			responseChan <- string(jsonResp)
			return
		}

		// Forward streamed responses
		for {
			select {
			case <-ctx.Done():
				response := ChatResponse{Error: "request cancelled"}
				jsonResp, _ := json.Marshal(response)
				responseChan <- string(jsonResp)
				return

			case resp, ok := <-streamChan:
				if !ok {
					return
				}
				response := ChatResponse{Response: resp}
				jsonResp, _ := json.Marshal(response)
				responseChan <- string(jsonResp)
			}
		}
	}()

	return responseChan, nil
}
