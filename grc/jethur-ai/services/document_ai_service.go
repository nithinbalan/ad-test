package services

import (
	"encoding/json"
	"fmt"
	"jethur/grc/documents/services"
	aiDTOs "jethur/grc/jethur-ai/dtos"
)

type DocumentAIServiceInterface interface {
	GetAISummary(documentId uint) (*aiDTOs.DocumentAISummaryDTO, error)
}

type DocumentAIService struct {
	documentService services.DocumentServiceInterface
}

func NewDocumentAIService(documentService services.DocumentServiceInterface) DocumentAIServiceInterface {
	return &DocumentAIService{
		documentService: documentService,
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

	// Get the latest version of the document
	documentVersion, err := as.documentService.GetDocumentVersionByID(document.LatestDocumentVersion.ID)
	if err != nil {
		return nil, err
	}

	// Initialize DTO
	aiSummaryDTO := &aiDTOs.DocumentAISummaryDTO{
		Summary: documentVersion.Summary,
	}

	// Unmarshal sample_questions (JSON string -> []string)
	if documentVersion.SampleQuestions != nil {
		if err := json.Unmarshal([]byte(*documentVersion.SampleQuestions), &aiSummaryDTO.SampleQuestions); err != nil {
			return nil, fmt.Errorf("error unmarshalling sample_questions: %w", err)
		}
	}

	// Unmarshal key_points (JSON string -> []string)
	if documentVersion.KeyPoints != nil {
		if err := json.Unmarshal([]byte(*documentVersion.KeyPoints), &aiSummaryDTO.KeyPoints); err != nil {
			return nil, fmt.Errorf("error unmarshalling key_points: %w", err)
		}
	}

	// Unmarshal faq (JSON string -> []FAQEntry)
	if documentVersion.FAQ != nil {
		if err := json.Unmarshal([]byte(*documentVersion.FAQ), &aiSummaryDTO.FAQ); err != nil {
			return nil, fmt.Errorf("error unmarshalling faq: %w", err)
		}
	}

	return aiSummaryDTO, nil
}
