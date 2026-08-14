package services

import (
	"errors"
	"fmt"
	"log"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/security-operations/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type CSLogDataEntityTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CSLogDataEntityType, error)
	GetByID(id uint) (*models.CSLogDataEntityType, error)
	Create(currentUserID uint, entityType *models.CSLogDataEntityType) (*models.CSLogDataEntityType, error)
	Update(id uint, currentUserID uint, entityType *models.CSLogDataEntityType) (*models.CSLogDataEntityType, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type CSLogDataEntityTypeService struct {
	repo           repositories.CSLogDataEntityTypeRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewCSLogDataEntityTypeService(repo repositories.CSLogDataEntityTypeRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface) CSLogDataEntityTypeServiceInterface {
	return &CSLogDataEntityTypeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (cldets *CSLogDataEntityTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CSLogDataEntityType, error) {
	return cldets.repo.GetAll(sort, sortBy, keyword, status)
}

func (cldets *CSLogDataEntityTypeService) GetByID(id uint) (*models.CSLogDataEntityType, error) {
	entityType, err := cldets.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return entityType, nil
}

func (cldets *CSLogDataEntityTypeService) Create(currentUserID uint, entityType *models.CSLogDataEntityType) (*models.CSLogDataEntityType, error) {
	if exists, err := cldets.repo.ExistsByTitle(entityType.TitleEn, entityType.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	created, err := cldets.repo.Create(entityType)
	if err != nil {
		return nil, err
	}

	go func() {
		err := cldets.activityLogSvc.LogActivity(
			models.CreateAction,
			models.CSLogDataEntityType{}.TableName(),
			created.ID,
			created,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return created, nil
}

func (cldets *CSLogDataEntityTypeService) Update(id uint, currentUserID uint, updated *models.CSLogDataEntityType) (*models.CSLogDataEntityType, error) {
	existing, err := cldets.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if updated.TitleEn != existing.TitleEn || updated.TitleAr != existing.TitleAr {
		if exists, err := cldets.repo.ExistsByTitleExcludingID(id, updated.TitleEn, updated.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := cldets.repo.Update(id, updated)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existing, updated, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := cldets.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.CSLogDataEntityType{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fmt.Printf("Failed to log update activity: %v\n", err)
		}
	}()
	return result, err
}

func (cldets *CSLogDataEntityTypeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deleted []map[string]any
	for _, id := range ids {
		entityType, err := cldets.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		if err := cldets.repo.Delete(id); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deleted = append(deleted, map[string]any{
			"id":       id,
			"title_en": entityType.TitleEn,
			"title_ar": entityType.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := cldets.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.CSLogDataEntityType{}.TableName(),
				successfulIds,
				deleted,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (cldets *CSLogDataEntityTypeService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := cldets.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := cldets.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.CSLogDataEntityType{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk activate activity: %v\n", err)
			}
		}()
	}

	return err
}

func (cldets *CSLogDataEntityTypeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := cldets.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := cldets.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.CSLogDataEntityType{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk deactivate activity: %v\n", err)
			}
		}()
	}

	return err
}
