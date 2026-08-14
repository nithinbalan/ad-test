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

type CSLogDataReceptionMechanismServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CSLogDataReceptionMechanism, error)
	GetByID(id uint) (*models.CSLogDataReceptionMechanism, error)
	Create(currentUserID uint, mechanism *models.CSLogDataReceptionMechanism) (*models.CSLogDataReceptionMechanism, error)
	Update(id uint, currentUserID uint, mechanism *models.CSLogDataReceptionMechanism) (*models.CSLogDataReceptionMechanism, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type CSLogDataReceptionMechanismService struct {
	repo           repositories.CSLogDataReceptionMechanismRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewCSLogDataReceptionMechanismService(repo repositories.CSLogDataReceptionMechanismRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface) CSLogDataReceptionMechanismServiceInterface {
	return &CSLogDataReceptionMechanismService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (cldrms *CSLogDataReceptionMechanismService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CSLogDataReceptionMechanism, error) {
	return cldrms.repo.GetAll(sort, sortBy, keyword, status)
}

func (cldrms *CSLogDataReceptionMechanismService) GetByID(id uint) (*models.CSLogDataReceptionMechanism, error) {
	mechanism, err := cldrms.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return mechanism, nil
}

func (cldrms *CSLogDataReceptionMechanismService) Create(currentUserID uint, mechanism *models.CSLogDataReceptionMechanism) (*models.CSLogDataReceptionMechanism, error) {
	if exists, err := cldrms.repo.ExistsByTitle(mechanism.TitleEn, mechanism.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	created, err := cldrms.repo.Create(mechanism)
	if err != nil {
		return nil, err
	}

	go func() {
		err := cldrms.activityLogSvc.LogActivity(
			models.CreateAction,
			models.CSLogDataReceptionMechanism{}.TableName(),
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

func (cldrms *CSLogDataReceptionMechanismService) Update(id uint, currentUserID uint, updated *models.CSLogDataReceptionMechanism) (*models.CSLogDataReceptionMechanism, error) {
	existing, err := cldrms.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if updated.TitleEn != existing.TitleEn || updated.TitleAr != existing.TitleAr {
		if exists, err := cldrms.repo.ExistsByTitleExcludingID(id, updated.TitleEn, updated.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := cldrms.repo.Update(id, updated)
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
		err := cldrms.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.CSLogDataReceptionMechanism{}.TableName(),
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

func (cldrms *CSLogDataReceptionMechanismService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deleted []map[string]any
	for _, id := range ids {
		mechanism, err := cldrms.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		if err := cldrms.repo.Delete(id); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deleted = append(deleted, map[string]any{
			"id":       id,
			"title_en": mechanism.TitleEn,
			"title_ar": mechanism.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := cldrms.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.CSLogDataReceptionMechanism{}.TableName(),
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

func (cldrms *CSLogDataReceptionMechanismService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := cldrms.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := cldrms.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.CSLogDataReceptionMechanism{}.TableName(),
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

func (cldrms *CSLogDataReceptionMechanismService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := cldrms.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := cldrms.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.CSLogDataReceptionMechanism{}.TableName(),
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
