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

type DocumentRequirementServiceInterface interface {
	GetAll(documentId uint, requirementGroupId uint, sort string, sortBy string, keyword string) ([]*models.DocumentRequirement, error)
	GetByID(documentId uint, requirementGroupId uint, id uint) (*models.DocumentRequirement, error)
	Create(documentId uint, requirementGroupId uint, currentUserID uint, dto *dtos.DocumentRequirementFormDTO) (uint, error)
	Update(documentId uint, requirementGroupId uint, id uint, currentUserID uint, dto *dtos.DocumentRequirementFormDTO) (*models.DocumentRequirement, error)
	Delete(documentId uint, requirementGroupId uint, ids []uint, currentUserID uint) ([]uint, error)
}

type DocumentRequirementService struct {
	repo                        repositories.DocumentRequirementRepositoryInterface
	activityLogSvc              services.ActivityLogServiceInterface
	documentRequirementGroupSvc DocumentRequirementGroupServiceInterface
}

func NewDocumentRequirementService(repo repositories.DocumentRequirementRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface, documentRequirementGroupSvc DocumentRequirementGroupServiceInterface) DocumentRequirementServiceInterface {
	return &DocumentRequirementService{
		repo:                        repo,
		activityLogSvc:              activityLogSvc,
		documentRequirementGroupSvc: documentRequirementGroupSvc,
	}
}

func (drs *DocumentRequirementService) GetAll(documentId uint, requirementGroupId uint, sort string, sortBy string, keyword string) ([]*models.DocumentRequirement, error) {
	_, err := drs.documentRequirementGroupSvc.GetByID(documentId, requirementGroupId)
	if err != nil {
		return nil, err
	}
	return drs.repo.GetAll(requirementGroupId, sort, sortBy, keyword)
}

func (drs *DocumentRequirementService) GetByID(documentId uint, requirementGroupId uint, id uint) (*models.DocumentRequirement, error) {
	_, err := drs.documentRequirementGroupSvc.GetByID(documentId, requirementGroupId)
	if err != nil {
		return nil, err
	}
	documentRequirement, err := drs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return documentRequirement, nil
}

func (drs *DocumentRequirementService) Create(documentId uint, requirementGroupId uint, currentUserID uint, dto *dtos.DocumentRequirementFormDTO) (uint, error) {
	_, err := drs.documentRequirementGroupSvc.GetByID(documentId, requirementGroupId)
	if err != nil {
		return 0, err
	}
	// Create model object from DTO
	requirementGroup := &models.DocumentRequirement{
		DocumentRequirementGroupID: requirementGroupId,
		TitleEn:                    dto.TitleEn,
		TitleAr:                    dto.TitleAr,
		CreatedBYID:                currentUserID,
	}

	// Save the entity
	id, err := drs.repo.Create(requirementGroup)
	if err != nil {
		return 0, err
	}

	// Log activity asynchronously
	go func() {
		err := drs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.DocumentRequirement{}.TableName(),
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

func (drs *DocumentRequirementService) Update(documentId uint, requirementGroupId uint, id uint, currentUserID uint, dto *dtos.DocumentRequirementFormDTO) (*models.DocumentRequirement, error) {
	_, err := drs.documentRequirementGroupSvc.GetByID(documentId, requirementGroupId)
	if err != nil {
		return nil, err
	}
	// Fetch existing record
	existingRequirementGroup, err := drs.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	existingRequirementGroup.TitleEn = dto.TitleEn
	existingRequirementGroup.TitleAr = dto.TitleAr

	existingRequirementGroup.UpdatedBYID = &currentUserID

	// Update the entity
	result, err := drs.repo.Update(id, existingRequirementGroup)
	if err != nil {
		return nil, err
	}

	// Log changes
	fieldsToCompare := []string{"TitleEn", "TitleAr"}
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
		err := drs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.DocumentRequirement{}.TableName(),
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

func (drs *DocumentRequirementService) Delete(documentId uint, requirementGroupId uint, ids []uint, currentUserID uint) ([]uint, error) {
	_, err := drs.documentRequirementGroupSvc.GetByID(documentId, requirementGroupId)
	if err != nil {
		return nil, err
	}
	var failedIds []uint
	var successfulIds []uint
	var deletedDocumentRequirement []map[string]interface{}
	for _, id := range ids {
		documentRequirement, err := drs.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = drs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedDocumentRequirement = append(deletedDocumentRequirement, map[string]interface{}{
			"id":       id,
			"title_en": documentRequirement.TitleEn,
			"title_ar": documentRequirement.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := drs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.DocumentRequirement{}.TableName(),
				successfulIds,
				deletedDocumentRequirement,
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
