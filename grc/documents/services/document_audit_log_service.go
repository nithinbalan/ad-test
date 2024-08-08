package services

import (
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/repositories"
	"jethur/grc/models"
	"jethur/utility"
)

type DocumentAuditLogServiceInterface interface {
	GetAll(documentID uint, params *documentDTOs.DocumentAuditLogListQueryParamsDTO) (*documentDTOs.DocumentAuditLogListResponseDTO, error)
	GetByType(documentAuditLogType string) (*models.DocumentAuditLogType, error)
	Create(documentAuditLogFormDTOs *documentDTOs.DocumentAuditLogFormDTO, tx ...utility.DatabaseTransaction) error
	CheckExists(documentID uint, documentVersionID uint, documentTypeID uint, userID uint, tx ...utility.DatabaseTransaction) (bool, error)
}

// Service and implementations
type DocumentAuditLogService struct {
	repo         repositories.DocumentAuditLogRepositoryInterface
	documentRepo repositories.DocumentRepositoryInterface
}

func NewDocumentAuditLogService(repo repositories.DocumentAuditLogRepositoryInterface,
	documentRepo repositories.DocumentRepositoryInterface,
) DocumentAuditLogServiceInterface {
	return &DocumentAuditLogService{
		repo:         repo,
		documentRepo: documentRepo,
	}
}

func (dals *DocumentAuditLogService) GetAll(documentID uint, params *documentDTOs.DocumentAuditLogListQueryParamsDTO) (*documentDTOs.DocumentAuditLogListResponseDTO, error) {
	tx := dals.repo.BeginTransaction()

	document, err := dals.documentRepo.GetFullDetailsByID(documentID, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	documentVersion, err := dals.documentRepo.GetDocumentVersionByID(document.LatestDocumentVersion.ID, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return dals.repo.GetAll(documentID, documentVersion.ID, params)
}

func (dals *DocumentAuditLogService) GetByType(documentAuditLogType string) (*models.DocumentAuditLogType, error) {
	typeResult, err := dals.repo.GetByType(documentAuditLogType)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return typeResult, nil
}

func (dals *DocumentAuditLogService) CheckExists(documentID uint, documentVersionID uint, documentAuditLogTypeID uint, userID uint, tx ...utility.DatabaseTransaction) (bool, error) {
	exists, err := dals.repo.Exists(documentID, documentVersionID, documentAuditLogTypeID, userID, tx...)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (dals *DocumentAuditLogService) Create(documentAuditLogFormDTO *documentDTOs.DocumentAuditLogFormDTO, tx ...utility.DatabaseTransaction) error {
	documentAuditLog := &models.DocumentAuditLog{
		DocumentID:             documentAuditLogFormDTO.DocumentID,
		DocumentVersionID:      documentAuditLogFormDTO.DocumentVersionID,
		UserID:                 documentAuditLogFormDTO.UserID,
		DocumentAuditLogTypeID: documentAuditLogFormDTO.DocumentAuditLogTypeID,
		CreatedAt:              documentAuditLogFormDTO.CreatedAt,
		CreatedBYID:            documentAuditLogFormDTO.CreatedByID,
	}

	if err := dals.repo.Create(documentAuditLog, tx...); err != nil {
		return err
	}
	return nil

}
