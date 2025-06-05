package services

import (
	"encoding/json"
	"fmt"
	"jethur/grc/documents/services"
	aiDTOs "jethur/grc/jethur-ai/dtos"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"os"
)

type DocumentAIServiceInterface interface {
	GetAISummary(documentId uint) (*aiDTOs.DocumentAISummaryDTO, error)
	GetDocumentRecommendations(documentId uint, currentUserID uint) (map[string]interface{}, error)
}

type DocumentAIService struct {
	documentService services.DocumentServiceInterface
	userService     organizationServices.UserServiceInterface
}

func NewDocumentAIService(documentService services.DocumentServiceInterface,
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
