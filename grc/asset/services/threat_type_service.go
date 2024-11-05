package services

import (
	"errors"
	"fmt"
	"jethur/grc/asset/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type ThreatTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThreatType, error)
	GetByID(id uint) (*models.ThreatType, error)
	Create(currentUserID uint, threatType *models.ThreatType) (*models.ThreatType, error)
	Update(id uint, currentUserID uint, threatType *models.ThreatType) (*models.ThreatType, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type ThreatTypeService struct {
	repo           repositories.ThreatTypeRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewThreatTypeService(repo repositories.ThreatTypeRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) ThreatTypeServiceInterface {
	return &ThreatTypeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (tts *ThreatTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThreatType, error) {
	return tts.repo.GetAll(sort, sortBy, keyword, status)
}

func (tts ThreatTypeService) GetByID(id uint) (*models.ThreatType, error) {
	threatType, err := tts.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return threatType, nil
}

func (tts *ThreatTypeService) Create(currentUserID uint, threatType *models.ThreatType) (*models.ThreatType, error) {
	if exists, err := tts.repo.ExistsByTitle(threatType.TitleEn, threatType.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdType, err := tts.repo.Create(threatType)
	if err != nil {
		return nil, err
	}

	go func() {
		err := tts.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ThreatType{}.TableName(),
			createdType.ID,
			createdType,
			currentUserID,
			nil,
		)

		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdType, nil
}

func (tts *ThreatTypeService) Update(id uint, currentUserID uint, updatedType *models.ThreatType) (*models.ThreatType, error) {
	existingType, err := tts.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedType.TitleEn != existingType.TitleEn || updatedType.TitleAr != existingType.TitleAr {
		if exists, err := tts.repo.ExistsByTitleExcludingID(id, updatedType.TitleEn, updatedType.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := tts.repo.Update(id, updatedType)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingType, updatedType, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := tts.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ThreatType{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()
	return result, err
}

func (tts *ThreatTypeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {
		_, err := tts.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = tts.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := tts.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ThreatType{}.TableName(),
				successfulIds,
				nil,
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

func (tts *ThreatTypeService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := tts.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := tts.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ThreatType{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}

func (tts *ThreatTypeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := tts.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := tts.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ThreatType{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}
