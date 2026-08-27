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

type PenetrationTestTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestType, error)
	GetByID(id uint) (*models.PenetrationTestType, error)
	Create(currentUserID uint, penetrationTestType *models.PenetrationTestType) (*models.PenetrationTestType, error)
	Update(id uint, currentUserID uint, penetrationTestType *models.PenetrationTestType) (*models.PenetrationTestType, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type PenetrationTestTypeService struct {
	repo           repositories.PenetrationTestTypeRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewPenetrationTestTypeService(repo repositories.PenetrationTestTypeRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface) PenetrationTestTypeServiceInterface {
	return &PenetrationTestTypeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ptts *PenetrationTestTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestType, error) {
	return ptts.repo.GetAll(sort, sortBy, keyword, status)
}

func (ptts *PenetrationTestTypeService) GetByID(id uint) (*models.PenetrationTestType, error) {
	penetrationTestType, err := ptts.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return penetrationTestType, nil
}

func (ptts *PenetrationTestTypeService) Create(currentUserID uint, penetrationTestType *models.PenetrationTestType) (*models.PenetrationTestType, error) {
	if exists, err := ptts.repo.ExistsByTitle(penetrationTestType.TitleEn, penetrationTestType.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	created, err := ptts.repo.Create(penetrationTestType)
	if err != nil {
		return nil, err
	}

	go func() {
		err := ptts.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PenetrationTestType{}.TableName(),
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

func (ptts *PenetrationTestTypeService) Update(id uint, currentUserID uint, updated *models.PenetrationTestType) (*models.PenetrationTestType, error) {
	existing, err := ptts.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if updated.TitleEn != existing.TitleEn || updated.TitleAr != existing.TitleAr {
		if exists, err := ptts.repo.ExistsByTitleExcludingID(id, updated.TitleEn, updated.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := ptts.repo.Update(id, updated)
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
		err := ptts.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.PenetrationTestType{}.TableName(),
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

func (ptts *PenetrationTestTypeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deleted []map[string]any
	for _, id := range ids {
		penetrationTestType, err := ptts.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		if err := ptts.repo.Delete(id); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deleted = append(deleted, map[string]any{
			"id":       id,
			"title_en": penetrationTestType.TitleEn,
			"title_ar": penetrationTestType.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := ptts.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.PenetrationTestType{}.TableName(),
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

func (ptts *PenetrationTestTypeService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := ptts.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := ptts.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.PenetrationTestType{}.TableName(),
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

func (ptts *PenetrationTestTypeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := ptts.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := ptts.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.PenetrationTestType{}.TableName(),
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
