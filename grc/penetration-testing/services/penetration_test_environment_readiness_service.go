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

type PenetrationTestEnvironmentReadinessServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestEnvironmentReadiness, error)
	GetByID(id uint) (*models.PenetrationTestEnvironmentReadiness, error)
	Create(currentUserID uint, penetrationTestEnvironmentReadiness *models.PenetrationTestEnvironmentReadiness) (*models.PenetrationTestEnvironmentReadiness, error)
	Update(id uint, currentUserID uint, penetrationTestEnvironmentReadiness *models.PenetrationTestEnvironmentReadiness) (*models.PenetrationTestEnvironmentReadiness, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type PenetrationTestEnvironmentReadinessService struct {
	repo           repositories.PenetrationTestEnvironmentReadinessRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewPenetrationTestEnvironmentReadinessService(repo repositories.PenetrationTestEnvironmentReadinessRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface) PenetrationTestEnvironmentReadinessServiceInterface {
	return &PenetrationTestEnvironmentReadinessService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (pters *PenetrationTestEnvironmentReadinessService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestEnvironmentReadiness, error) {
	return pters.repo.GetAll(sort, sortBy, keyword, status)
}

func (pters *PenetrationTestEnvironmentReadinessService) GetByID(id uint) (*models.PenetrationTestEnvironmentReadiness, error) {
	penetrationTestEnvironmentReadiness, err := pters.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return penetrationTestEnvironmentReadiness, nil
}

func (pters *PenetrationTestEnvironmentReadinessService) Create(currentUserID uint, penetrationTestEnvironmentReadiness *models.PenetrationTestEnvironmentReadiness) (*models.PenetrationTestEnvironmentReadiness, error) {
	if exists, err := pters.repo.ExistsByTitle(penetrationTestEnvironmentReadiness.TitleEn, penetrationTestEnvironmentReadiness.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	created, err := pters.repo.Create(penetrationTestEnvironmentReadiness)
	if err != nil {
		return nil, err
	}

	go func() {
		err := pters.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PenetrationTestEnvironmentReadiness{}.TableName(),
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

func (pters *PenetrationTestEnvironmentReadinessService) Update(id uint, currentUserID uint, updated *models.PenetrationTestEnvironmentReadiness) (*models.PenetrationTestEnvironmentReadiness, error) {
	existing, err := pters.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if updated.TitleEn != existing.TitleEn || updated.TitleAr != existing.TitleAr {
		if exists, err := pters.repo.ExistsByTitleExcludingID(id, updated.TitleEn, updated.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := pters.repo.Update(id, updated)
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
		err := pters.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.PenetrationTestEnvironmentReadiness{}.TableName(),
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

func (pters *PenetrationTestEnvironmentReadinessService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deleted []map[string]any
	for _, id := range ids {
		penetrationTestEnvironmentReadiness, err := pters.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		if err := pters.repo.Delete(id); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deleted = append(deleted, map[string]any{
			"id":       id,
			"title_en": penetrationTestEnvironmentReadiness.TitleEn,
			"title_ar": penetrationTestEnvironmentReadiness.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := pters.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.PenetrationTestEnvironmentReadiness{}.TableName(),
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

func (pters *PenetrationTestEnvironmentReadinessService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := pters.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := pters.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.PenetrationTestEnvironmentReadiness{}.TableName(),
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

func (pters *PenetrationTestEnvironmentReadinessService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := pters.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := pters.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.PenetrationTestEnvironmentReadiness{}.TableName(),
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
