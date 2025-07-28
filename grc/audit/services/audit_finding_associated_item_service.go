package services

import (
	"fmt"
	"jethur/grc/audit/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	"jethur/grc/models"
	"jethur/utility"
	"time"
)

type AuditFindingAssociatedItemServiceInterface interface {
	GetAllAuditFindingDocuments(findingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	CreateAuditFindingDocument(findingID uint, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error
	DeleteAuditFindingDocument(findingID uint, documentID uint) error
}

// Service and implementations
type AuditFindingAssociatedItemService struct {
	repo            repositories.AuditFindingAssociatedItemRepositoryInterface
	AuditFindingSvc AuditFindingServiceInterface
	validationSvc   commonServices.AssociatedItemValidationServiceInterface
}

func NewAuditFindingAssociatedItemService(repo repositories.AuditFindingAssociatedItemRepositoryInterface,
	AuditFindingSvc AuditFindingServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,
) AuditFindingAssociatedItemServiceInterface {
	return &AuditFindingAssociatedItemService{
		repo:            repo,
		AuditFindingSvc: AuditFindingSvc,
		validationSvc:   validationSvc,
	}
}

func (afais *AuditFindingAssociatedItemService) GetAllAuditFindingDocuments(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	if _, err := afais.AuditFindingSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return afais.repo.GetAllAuditFindingDocuments(assetID, params)
}

func (afais *AuditFindingAssociatedItemService) CreateAuditFindingDocument(findingID, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error {

	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("AuditFinding does not exist for ID: %d", findingID), "not_exist", "audit_finding_id")
	}

	if err := afais.validationSvc.ValidateDocumentMapping(findingID, dto.DocumentIDs, "audit_finding_associated_documents", "audit_finding_id"); err != nil {
		return err
	}

	tx := afais.repo.BeginTransaction()
	var mappings []*models.AuditFindingAssociatedDocument
	currentTime := time.Now()

	for _, documentID := range dto.DocumentIDs {
		mappings = append(mappings, &models.AuditFindingAssociatedDocument{
			DocumentID:     documentID,
			AuditFindingID: findingID,
			CreatedByID:    userID,
			CreatedAt:      currentTime,
		})
	}

	if err := afais.repo.CreateAuditFindingDocument(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (afais *AuditFindingAssociatedItemService) DeleteAuditFindingDocument(findingID uint, documentID uint) error {
	// Validate control existence
	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("AuditFinding does not exist for ID: %d", findingID), "not_exist", "audit_finding_id")
	}

	tx := afais.repo.BeginTransaction()

	if err := afais.repo.DeleteAuditFindingDocument(findingID, documentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
