package services

import (
	"fmt"
	"jethur/grc/audit/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"log"
)

type AuditActivityTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditActivityType, error)
	GetByID(id uint) (*models.AuditActivityType, error)
	Create(currentUserID uint, auditActivityType *models.AuditActivityType) (*models.AuditActivityType, error)
	Update(id uint, currentUserID uint, auditActivityType *models.AuditActivityType) (*models.AuditActivityType, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type AuditActivityTypeService struct {
	repo           repositories.AuditActivityTypeRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewAuditActivityTypeService(repo repositories.AuditActivityTypeRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) AuditActivityTypeServiceInterface {
	return &AuditActivityTypeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (aats *AuditActivityTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditActivityType, error) {
	return aats.repo.GetAll(sort, sortBy, keyword, status)
}

func (aats AuditActivityTypeService) GetByID(id uint) (*models.AuditActivityType, error) {
	return aats.repo.GetByID(id)
}

func (aats *AuditActivityTypeService) Create(currentUserID uint, auditActivityType *models.AuditActivityType) (*models.AuditActivityType, error) {

	if exists, err := aats.repo.ExistsByTitle(auditActivityType.TitleEn, auditActivityType.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdAuditActivityType, err := aats.repo.Create(auditActivityType)
	if err != nil {
		return nil, err
	}

	go func() {
		err := aats.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditActivityType{}.TableName(),
			createdAuditActivityType.ID,
			createdAuditActivityType,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdAuditActivityType, nil
}

func (aats *AuditActivityTypeService) Update(id uint, currentUserID uint, updatedAuditActivityType *models.AuditActivityType) (*models.AuditActivityType, error) {
	existingAuditActivityType, err := aats.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedAuditActivityType.TitleEn != existingAuditActivityType.TitleEn || updatedAuditActivityType.TitleAr != existingAuditActivityType.TitleAr {
		if exists, err := aats.repo.ExistsByTitleExcludingID(id, updatedAuditActivityType.TitleEn, updatedAuditActivityType.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := aats.repo.Update(id, updatedAuditActivityType)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "DescriptionEn", "DescriptionAr", "Status"}
	changes, err := utility.CompareFields(existingAuditActivityType, updatedAuditActivityType, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := aats.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditActivityType{}.TableName(),
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

func (aats *AuditActivityTypeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := aats.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = aats.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := aats.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.AuditActivityType{}.TableName(),
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

	return failedIds, nil
}

func (aats *AuditActivityTypeService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := aats.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := aats.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.AuditActivityType{}.TableName(),
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

func (aats *AuditActivityTypeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := aats.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := aats.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.AuditActivityType{}.TableName(),
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
