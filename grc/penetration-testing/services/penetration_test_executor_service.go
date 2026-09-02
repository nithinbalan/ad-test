package services

import (
	"errors"
	"fmt"
	"log"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/penetration-testing/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type PenetrationTestExecutorServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestExecutor, error)
	GetByID(id uint) (*models.PenetrationTestExecutor, error)
	Create(currentUserID uint, penetrationTestExecutor *models.PenetrationTestExecutor) (*models.PenetrationTestExecutor, error)
	Update(id uint, currentUserID uint, penetrationTestExecutor *models.PenetrationTestExecutor) (*models.PenetrationTestExecutor, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type PenetrationTestExecutorService struct {
	repo           repositories.PenetrationTestExecutorRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewPenetrationTestExecutorService(repo repositories.PenetrationTestExecutorRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface) PenetrationTestExecutorServiceInterface {
	return &PenetrationTestExecutorService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ptes *PenetrationTestExecutorService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestExecutor, error) {
	return ptes.repo.GetAll(sort, sortBy, keyword, status)
}

func (ptes *PenetrationTestExecutorService) GetByID(id uint) (*models.PenetrationTestExecutor, error) {
	penetrationTestExecutor, err := ptes.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return penetrationTestExecutor, nil
}

func (ptes *PenetrationTestExecutorService) Create(currentUserID uint, penetrationTestExecutor *models.PenetrationTestExecutor) (*models.PenetrationTestExecutor, error) {
	if exists, err := ptes.repo.ExistsByTitle(penetrationTestExecutor.TitleEn, penetrationTestExecutor.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	created, err := ptes.repo.Create(penetrationTestExecutor)
	if err != nil {
		return nil, err
	}

	go func() {
		err := ptes.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PenetrationTestExecutor{}.TableName(),
			created.ID,
			created,
			currentUserID,
			"penetration_testing", "penetration_testing_setting", "create penetration test executor", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return created, nil
}

func (ptes *PenetrationTestExecutorService) Update(id uint, currentUserID uint, updated *models.PenetrationTestExecutor) (*models.PenetrationTestExecutor, error) {
	existing, err := ptes.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if updated.TitleEn != existing.TitleEn || updated.TitleAr != existing.TitleAr {
		if exists, err := ptes.repo.ExistsByTitleExcludingID(id, updated.TitleEn, updated.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := ptes.repo.Update(id, updated)
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
		err := ptes.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.PenetrationTestExecutor{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"penetration_testing", "penetration_testing_setting", "update penetration test executor", nil,
		)
		if err != nil {
			fmt.Printf("Failed to log update activity: %v\n", err)
		}
	}()
	return result, err
}

func (ptes *PenetrationTestExecutorService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deleted []map[string]any
	for _, id := range ids {
		penetrationTestExecutor, err := ptes.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		if err := ptes.repo.Delete(id); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
			deleted = append(deleted, map[string]any{
				"id":       id,
				"title_en": penetrationTestExecutor.TitleEn,
				"title_ar": penetrationTestExecutor.TitleAr,
			})
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := ptes.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.PenetrationTestExecutor{}.TableName(),
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

func (ptes *PenetrationTestExecutorService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := ptes.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := ptes.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.PenetrationTestExecutor{}.TableName(),
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

func (ptes *PenetrationTestExecutorService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := ptes.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := ptes.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.PenetrationTestExecutor{}.TableName(),
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
