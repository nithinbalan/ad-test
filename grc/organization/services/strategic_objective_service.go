package services

import (
	"errors"
	"fmt"
	logModels "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type StrategicObjectiveServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*dtos.StrategicObjectiveListDTO, error)
	GetByID(id uint) (*models.StrategicObjective, error)
	Create(currentUserID uint, strategicObjectiveFormDTO *dtos.StrategicObjectiveFormDTO) (uint, error)
	Update(id uint, currentUserID uint, strategicObjectiveFormDTO *dtos.StrategicObjectiveFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type StrategicObjectiveService struct {
	repo           repositories.StrategicObjectiveRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewStrategicObjectiveService(repo repositories.StrategicObjectiveRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) StrategicObjectiveServiceInterface {
	return &StrategicObjectiveService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (sos *StrategicObjectiveService) GetAll(sort string, sortBy string, keyword string, status string) ([]*dtos.StrategicObjectiveListDTO, error) {
	return sos.repo.GetAll(sort, sortBy, keyword, status)
}

func (sos StrategicObjectiveService) GetByID(id uint) (*models.StrategicObjective, error) {
	strategicObjective, err := sos.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return strategicObjective, nil

}

func (sos *StrategicObjectiveService) Create(currentUserID uint, strategicObjectiveFormDTO *dtos.StrategicObjectiveFormDTO) (uint, error) {
	if exists, err := sos.repo.ExistsByTitle(strategicObjectiveFormDTO.TitleEn, strategicObjectiveFormDTO.TitleAr); err != nil {
		return 0, err
	} else if exists {
		return 0, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	strategicObjective := &models.StrategicObjective{
		TitleEn:   strategicObjectiveFormDTO.TitleEn,
		TitleAr:   strategicObjectiveFormDTO.TitleAr,
		Status:    strategicObjectiveFormDTO.Status,
		CreatedBy: currentUserID,
	}

	createdStrategicObjectiveID, err := sos.repo.Create(strategicObjective)
	if err != nil {
		return 0, err
	}

	go func() {
		err := sos.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.StrategicObjective{}.TableName(),
			createdStrategicObjectiveID,
			strategicObjectiveFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdStrategicObjectiveID, nil
}

func (sos *StrategicObjectiveService) Update(id uint, currentUserID uint, strategicObjectiveFormDTO *dtos.StrategicObjectiveFormDTO) error {
	existingStrategicObjective, err := sos.GetByID(id)

	if err != nil {
		return err
	}

	if strategicObjectiveFormDTO.TitleEn != existingStrategicObjective.TitleEn || strategicObjectiveFormDTO.TitleAr != existingStrategicObjective.TitleAr {
		if exists, err := sos.repo.ExistsByTitleExcludingID(id, strategicObjectiveFormDTO.TitleEn, strategicObjectiveFormDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	existingStrategicObjective.TitleEn = strategicObjectiveFormDTO.TitleEn
	existingStrategicObjective.TitleAr = strategicObjectiveFormDTO.TitleAr
	existingStrategicObjective.Status = strategicObjectiveFormDTO.Status
	existingStrategicObjective.UpdatedBy = &currentUserID

	err = sos.repo.Update(existingStrategicObjective)

	if err != nil {
		return err
	}

	//@todo - need to capture activity log
	return nil
}

func (sos *StrategicObjectiveService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := sos.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = sos.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := sos.activityLogSvc.LogBulkActivities(
				logModels.DeleteAction, // Make sure this action constant is defined appropriately
				models.StrategicObjective{}.TableName(),
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

func (sos *StrategicObjectiveService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := sos.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := sos.activityLogSvc.LogBulkActivities(
				logModels.ActivateAction,
				models.StrategicObjective{}.TableName(),
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

func (sos *StrategicObjectiveService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := sos.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := sos.activityLogSvc.LogBulkActivities(
				logModels.DeactivateAction,
				models.StrategicObjective{}.TableName(),
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
