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

type PenetrationTestPlanServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestPlan, error)
	GetByID(id uint) (*models.PenetrationTestPlan, error)
	Create(currentUserID uint, penetrationTestPlan *models.PenetrationTestPlan) (*models.PenetrationTestPlan, error)
	Update(id uint, currentUserID uint, penetrationTestPlan *models.PenetrationTestPlan) (*models.PenetrationTestPlan, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type PenetrationTestPlanService struct {
	repo           repositories.PenetrationTestPlanRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewPenetrationTestPlanService(repo repositories.PenetrationTestPlanRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface) PenetrationTestPlanServiceInterface {
	return &PenetrationTestPlanService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ptps *PenetrationTestPlanService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestPlan, error) {
	return ptps.repo.GetAll(sort, sortBy, keyword, status)
}

func (ptps *PenetrationTestPlanService) GetByID(id uint) (*models.PenetrationTestPlan, error) {
	penetrationTestPlan, err := ptps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return penetrationTestPlan, nil
}

func (ptps *PenetrationTestPlanService) Create(currentUserID uint, penetrationTestPlan *models.PenetrationTestPlan) (*models.PenetrationTestPlan, error) {
	if exists, err := ptps.repo.ExistsByTitle(penetrationTestPlan.TitleEn, penetrationTestPlan.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	created, err := ptps.repo.Create(penetrationTestPlan)
	if err != nil {
		return nil, err
	}

	go func() {
		err := ptps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PenetrationTestPlan{}.TableName(),
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

func (ptps *PenetrationTestPlanService) Update(id uint, currentUserID uint, updated *models.PenetrationTestPlan) (*models.PenetrationTestPlan, error) {
	existing, err := ptps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if updated.TitleEn != existing.TitleEn || updated.TitleAr != existing.TitleAr {
		if exists, err := ptps.repo.ExistsByTitleExcludingID(id, updated.TitleEn, updated.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := ptps.repo.Update(id, updated)
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
		err := ptps.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.PenetrationTestPlan{}.TableName(),
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

func (ptps *PenetrationTestPlanService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deleted []map[string]any
	for _, id := range ids {
		penetrationTestPlan, err := ptps.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		if err := ptps.repo.Delete(id); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deleted = append(deleted, map[string]any{
			"id":       id,
			"title_en": penetrationTestPlan.TitleEn,
			"title_ar": penetrationTestPlan.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := ptps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.PenetrationTestPlan{}.TableName(),
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

func (ptps *PenetrationTestPlanService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := ptps.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := ptps.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.PenetrationTestPlan{}.TableName(),
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

func (ptps *PenetrationTestPlanService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := ptps.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := ptps.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.PenetrationTestPlan{}.TableName(),
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
