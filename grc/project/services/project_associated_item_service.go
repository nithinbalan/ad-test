package services

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	"jethur/grc/models"
	"jethur/grc/project/repositories"
	"jethur/utility"
	"time"
)

type ProjectAssociatedItemServiceInterface interface {
	GetAllProjectDocuments(projectID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	CreateProjectDocument(projectID uint, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error
	DeleteProjectDocument(projectID uint, documentID uint) error
}

// Service and implementations
type ProjectAssociatedItemService struct {
	repo          repositories.ProjectAssociatedItemRepositoryInterface
	projectSvc    ProjectServiceInterface
	validationSvc commonServices.AssociatedItemValidationServiceInterface
}

func NewProjectAssociatedItemService(repo repositories.ProjectAssociatedItemRepositoryInterface,
	projectSvc ProjectServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,
) ProjectAssociatedItemServiceInterface {
	return &ProjectAssociatedItemService{
		repo:          repo,
		projectSvc:    projectSvc,
		validationSvc: validationSvc,
	}
}

func (pais *ProjectAssociatedItemService) GetAllProjectDocuments(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	if _, err := pais.projectSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllProjectDocuments(assetID, params)
}

func (pais *ProjectAssociatedItemService) CreateProjectDocument(projectID, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error {

	if _, err := pais.projectSvc.GetByID(projectID); err != nil {
		return utility.NewValidationError(fmt.Errorf("project does not exist for ID: %d", projectID), "not_exist", "asset_id")
	}

	if err := pais.validationSvc.ValidateDocumentMapping(projectID, dto.DoumentIDs, "project_associated_documents", "project_id"); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()
	var mappings []*models.ProjectAssociatedDocument
	currentTime := time.Now()

	for _, documentID := range dto.DoumentIDs {
		mappings = append(mappings, &models.ProjectAssociatedDocument{
			DocumentID:  documentID,
			ProjectID:   projectID,
			CreatedByID: userID,
			CreatedAt:   currentTime,
		})
	}

	if err := pais.repo.CreateProjectDocument(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *ProjectAssociatedItemService) DeleteProjectDocument(projectID uint, documentID uint) error {
	// Validate control existence
	if _, err := pais.projectSvc.GetByID(projectID); err != nil {
		return utility.NewValidationError(fmt.Errorf("project does not exist for ID: %d", projectID), "not_exist", "asset_register_id")
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteProjectDocument(projectID, documentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
