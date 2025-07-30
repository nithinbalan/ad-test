package services

import (
	"context"
	"encoding/json"
	"fmt"
	"jethur/grc/documents/services"
	aiDTOs "jethur/grc/jethur-ai/dtos"
	organizationServices "jethur/grc/organization/services"
	settingsServices "jethur/grc/settings/services"
	"jethur/utility"
)

type DocumentAIServiceInterface interface {
	GetAISummary(documentId uint) (*aiDTOs.DocumentAISummaryDTO, error)
	GetDocumentRecommendations(documentId uint, currentUserID uint) (map[string]interface{}, error)
	ProcessChat(ctx context.Context, content, role, documentID, contextID string) (<-chan string, error)
	GetDocumentRequirements(documentID uint) ([]map[string]interface{}, error)
}

type DocumentAIService struct {
	documentService     services.DocumentServiceInterface
	userService         organizationServices.UserServiceInterface
	systemPreferenceSvc settingsServices.SystemPreferenceServiceInterface
}

func NewDocumentAIService(
	documentService services.DocumentServiceInterface,
	userService organizationServices.UserServiceInterface,
	systemPreferenceSvc settingsServices.SystemPreferenceServiceInterface,
) DocumentAIServiceInterface {
	return &DocumentAIService{
		documentService:     documentService,
		userService:         userService,
		systemPreferenceSvc: systemPreferenceSvc,
	}
}

func (as *DocumentAIService) GetAISummary(documentId uint) (*aiDTOs.DocumentAISummaryDTO, error) {
	systemPreference, err := as.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}
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

func (dai *DocumentAIService) GetDocumentRecommendations(documentId uint, currentUserID uint) (map[string]interface{}, error) {
	systemPreference, err := dai.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}
	userDetails, err := dai.userService.GetCurrentUserDetails(currentUserID)
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

	return utility.CallAIAPI("/api/v1/policy/policy-recomendation", payload)
}

func (dai *DocumentAIService) ProcessChat(ctx context.Context, content, role, documentID, contextID string) (<-chan string, error) {
	systemPreference, err := dai.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	payload := map[string]interface{}{
		"content":     content,
		"role":        role,
		"document_id": documentID,
		"context_id":  contextID,
	}

	return utility.StreamAIAPI(ctx, "/api/v1/document-chat", payload)
}

func (dai *DocumentAIService) GetDocumentRequirements(documentID uint) ([]map[string]interface{}, error) {
	systemPreference, err := dai.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/policy/%d/requirements-extracting", documentID)
	return utility.CallAIAPIArrayResponse(path)
}
