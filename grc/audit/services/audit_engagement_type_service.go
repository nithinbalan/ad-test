package services

import (
	"fmt"
	"jethur/grc/audit/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"log"
)

type AuditEngagementTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditEngagementType, error)
	GetByID(id uint) (*models.AuditEngagementType, error)
	Create(currentUserID uint, auditEngagementType *models.AuditEngagementType) (*models.AuditEngagementType, error)
	Update(id uint, currentUserID uint, auditEngagementType *models.AuditEngagementType) (*models.AuditEngagementType, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type AuditEngagementTypeService struct {
	repo           repositories.AuditEngagementTypeRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewAuditEngagementTypeService(repo repositories.AuditEngagementTypeRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) AuditEngagementTypeServiceInterface {
	return &AuditEngagementTypeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (aets *AuditEngagementTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditEngagementType, error) {
	return aets.repo.GetAll(sort, sortBy, keyword, status)
}

func (aets AuditEngagementTypeService) GetByID(id uint) (*models.AuditEngagementType, error) {
	return aets.repo.GetByID(id)
}

func (aets *AuditEngagementTypeService) Create(currentUserID uint, auditEngagementType *models.AuditEngagementType) (*models.AuditEngagementType, error) {

	if exists, err := aets.repo.ExistsByTitle(auditEngagementType.TitleEn, auditEngagementType.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdEngagementType, err := aets.repo.Create(auditEngagementType)
	if err != nil {
		return nil, err
	}

	go func() {
		err := aets.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditEngagementType{}.TableName(),
			createdEngagementType.ID,
			createdEngagementType,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdEngagementType, nil
}

func (aets *AuditEngagementTypeService) Update(id uint, currentUserID uint, updatedEngagementType *models.AuditEngagementType) (*models.AuditEngagementType, error) {
	existingEngagementType, err := aets.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedEngagementType.TitleEn != existingEngagementType.TitleEn || updatedEngagementType.TitleAr != existingEngagementType.TitleAr {
		if exists, err := aets.repo.ExistsByTitleExcludingID(id, updatedEngagementType.TitleEn, updatedEngagementType.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := aets.repo.Update(id, updatedEngagementType)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "DescriptionEn", "DescriptionAr", "Status"}
	changes, err := utility.CompareFields(existingEngagementType, updatedEngagementType, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := aets.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditEngagementType{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()
	return result, err
}

func (aets *AuditEngagementTypeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedAuditEngagementType []map[string]interface{}
	for _, id := range ids {

		engagementType, err := aets.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = aets.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
		deletedAuditEngagementType = append(deletedAuditEngagementType, map[string]interface{}{
			"id":       id,
			"title_en": engagementType.TitleEn,
			"title_ar": engagementType.TitleAr,
		})
	}

	if len(successfulIds) > 0 {

		go func() {
			err := aets.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.AuditEngagementType{}.TableName(),
				successfulIds,
				deletedAuditEngagementType,
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

func (aets *AuditEngagementTypeService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := aets.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := aets.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.AuditEngagementType{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}

func (aets *AuditEngagementTypeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := aets.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := aets.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.AuditEngagementType{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				// You might want to log this error or handle it appropriately
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}
