package services

import (
	"errors"
	"fmt"
	"jethur/grc/documents/dtos"
	"jethur/grc/documents/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type DocumentRequirementGroupServiceInterface interface {
	GetAll(documentId uint, sort string, sortBy string, keyword string, status string) ([]*dtos.DocumentRequirementGroupListDTO, error)
	GetByID(documentId uint, id uint) (*models.DocumentRequirementGroup, error)
	Create(documentId uint, currentUserID uint, dto *dtos.DocumentRequirementGroupFormDTO) (uint, error)
	Update(documentId uint, id uint, currentUserID uint, dto *dtos.DocumentRequirementGroupFormDTO) (*models.DocumentRequirementGroup, error)
	Delete(documentId uint, ids []uint, currentUserID uint) ([]uint, error)
	Activate(documentId uint, ids []uint, currentUserID uint) error
	Deactivate(documentId uint, ids []uint, currentUserID uint) error
}

type DocumentRequirementGroupService struct {
	repo           repositories.DocumentRequirementGroupRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
	documentSvc    DocumentServiceInterface
}

func NewDocumentRequirementGroupService(repo repositories.DocumentRequirementGroupRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface, documentSvc DocumentServiceInterface) DocumentRequirementGroupServiceInterface {
	return &DocumentRequirementGroupService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		documentSvc:    documentSvc,
	}
}

func (drgs *DocumentRequirementGroupService) GetAll(documentId uint, sort string, sortBy string, keyword string, status string) ([]*dtos.DocumentRequirementGroupListDTO, error) {
	_, err := drgs.documentSvc.GetByID(documentId, nil)
	if err != nil {
		return nil, err
	}

	documentRequirementGroups, err := drgs.repo.GetAll(documentId, sort, sortBy, keyword, status)
	if err != nil {
		return nil, err
	}

	// Convert Model to DTO
	var dtoList []*dtos.DocumentRequirementGroupListDTO
	for _, drg := range documentRequirementGroups {
		dto := &dtos.DocumentRequirementGroupListDTO{
			ID:                   drg.ID,
			DocumentID:           drg.DocumentID,
			TitleEn:              drg.TitleEn,
			TitleAr:              drg.TitleAr,
			Status:               drg.Status,
			CreatedByID:          drg.CreatedBYID,
			DocumentRequirements: []dtos.DocumentRequirementDTO{}, // Initialize as empty array
		}
		// Convert DocumentRequirements to DTO
		if len(drg.DocumentRequirements) > 0 {
			for _, requirement := range drg.DocumentRequirements {
				requirementDTO := dtos.DocumentRequirementDTO{
					ID:                         requirement.ID,
					DocumentRequirementGroupID: requirement.DocumentRequirementGroupID,
					TitleEn:                    requirement.TitleEn,
					TitleAr:                    requirement.TitleAr,
				}
				dto.DocumentRequirements = append(dto.DocumentRequirements, requirementDTO)
			}
		}

		dtoList = append(dtoList, dto)
	}

	return dtoList, nil
}

func (drgs *DocumentRequirementGroupService) GetByID(documentId uint, id uint) (*models.DocumentRequirementGroup, error) {
	_, err := drgs.documentSvc.GetByID(documentId, nil)
	if err != nil {
		return nil, err
	}
	documentRequirementGroup, err := drgs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return documentRequirementGroup, nil
}

func (drgs *DocumentRequirementGroupService) Create(documentId uint, currentUserID uint, dto *dtos.DocumentRequirementGroupFormDTO) (uint, error) {
	_, err := drgs.documentSvc.GetByID(documentId, nil)
	if err != nil {
		return 0, err
	}
	// Check if title is unique
	if exists, err := drgs.repo.ExistsByTitle(documentId, dto.TitleEn, dto.TitleAr); err != nil {
		return 0, err
	} else if exists {
		return 0, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	// Create model object from DTO
	requirementGroup := &models.DocumentRequirementGroup{
		DocumentID:  documentId,
		TitleEn:     dto.TitleEn,
		TitleAr:     dto.TitleAr,
		Status:      dto.Status,
		CreatedBYID: currentUserID,
	}

	// Save the entity
	id, err := drgs.repo.Create(requirementGroup)
	if err != nil {
		return 0, err
	}

	// Log activity asynchronously
	go func() {
		err := drgs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.DocumentRequirementGroup{}.TableName(),
			id,
			requirementGroup,
			currentUserID,
			nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return id, nil
}

func (drgs *DocumentRequirementGroupService) Update(documentId uint, id uint, currentUserID uint, dto *dtos.DocumentRequirementGroupFormDTO) (*models.DocumentRequirementGroup, error) {
	_, err := drgs.documentSvc.GetByID(documentId, nil)
	if err != nil {
		return nil, err
	}
	// Fetch existing record
	existingRequirementGroup, err := drgs.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if exists, err := drgs.repo.ExistsByTitleExcludingID(documentId, id, dto.TitleEn, dto.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}
	existingRequirementGroup.TitleEn = dto.TitleEn
	existingRequirementGroup.TitleAr = dto.TitleAr
	existingRequirementGroup.Status = dto.Status

	existingRequirementGroup.UpdatedBYID = &currentUserID

	// Update the entity
	result, err := drgs.repo.Update(id, existingRequirementGroup)
	if err != nil {
		return nil, err
	}

	// Log changes
	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingRequirementGroup, result, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}
	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	// Log activity asynchronously
	go func() {
		err := drgs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.DocumentRequirementGroup{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}()

	return result, nil
}

func (drgs *DocumentRequirementGroupService) Delete(documentId uint, ids []uint, currentUserID uint) ([]uint, error) {
	_, err := drgs.documentSvc.GetByID(documentId, nil)
	if err != nil {
		return nil, err
	}
	var failedIds []uint
	var successfulIds []uint
	var deletedDocumentRequirementGroup []map[string]interface{}
	for _, id := range ids {
		documentRequirementGroup, err := drgs.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = drgs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedDocumentRequirementGroup = append(deletedDocumentRequirementGroup, map[string]interface{}{
			"id":       id,
			"title_en": documentRequirementGroup.TitleEn,
			"title_ar": documentRequirementGroup.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := drgs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.DocumentRequirementGroup{}.TableName(),
				successfulIds,
				deletedDocumentRequirementGroup,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (drgs *DocumentRequirementGroupService) Activate(documentId uint, ids []uint, currentUserID uint) error {
	_, err := drgs.documentSvc.GetByID(documentId, nil)
	if err != nil {
		return err
	}
	successfulIds, err := drgs.repo.Activate(ids)
	if len(successfulIds) > 0 {
		go func() {
			err := drgs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.DocumentRequirementGroup{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk activate activity: %v\n", err)
			}
		}()
	}
	return err
}

func (drgs *DocumentRequirementGroupService) Deactivate(documentId uint, ids []uint, currentUserID uint) error {
	_, err := drgs.documentSvc.GetByID(documentId, nil)
	if err != nil {
		return err
	}
	successfulIds, err := drgs.repo.Deactivate(ids)
	if len(successfulIds) > 0 {
		go func() {
			err := drgs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.DocumentRequirementGroup{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk deactivate activity: %v\n", err)
			}
		}()
	}
	return err
}
