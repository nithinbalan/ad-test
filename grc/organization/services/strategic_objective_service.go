package services

import (
	"fmt"
	logModels "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	"jethur/grc/organization/models"
	"jethur/grc/organization/repositories"
	"jethur/utility"
	"log"
)

type StrategicObjectiveServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.StrategicObjective, error)
	GetByID(id uint) (*models.StrategicObjective, error)
	Create(currentUserID uint, strategicObjective *models.StrategicObjective) (*models.StrategicObjective, error)
	Update(id uint, currentUserID uint, strategicObjective *models.StrategicObjective) (*models.StrategicObjective, error)
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

func (ccs *StrategicObjectiveService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.StrategicObjective, error) {
	return ccs.repo.GetAll(sort, sortBy, keyword, status)
}

func (ccs StrategicObjectiveService) GetByID(id uint) (*models.StrategicObjective, error) {
	return ccs.repo.GetByID(id)
}

func (ccs *StrategicObjectiveService) Create(currentUserID uint, strategicObjective *models.StrategicObjective) (*models.StrategicObjective, error) {

	if exists, err := ccs.repo.ExistsByTitle(strategicObjective.TitleEn, strategicObjective.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	strategicObjective.CreatedBy = currentUserID

	createdStrategicObjective, err := ccs.repo.Create(strategicObjective)
	if err != nil {
		return nil, err
	}

	// go func() {
	// 	err := ccs.activityLogSvc.LogActivity(
	// 		logModels.CreateAction,
	// 		models.StrategicObjective{}.TableName(),
	// 		createdStrategicObjective.ID,
	// 		createdStrategicObjective,
	// 		currentUserID,
	// 		nil,
	// 	)

	// 	if err != nil {
	// 		// @todo - save to log file
	// 		log.Printf("Failed to log activity: %v", err)
	// 	}
	// }()

	return createdStrategicObjective, nil
}

func (ccs *StrategicObjectiveService) Update(id uint, currentUserID uint, updatedStrategicObjective *models.StrategicObjective) (*models.StrategicObjective, error) {
	existingStrategicObjective, err := ccs.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedStrategicObjective.TitleEn != existingStrategicObjective.TitleEn || updatedStrategicObjective.TitleAr != existingStrategicObjective.TitleAr {
		if exists, err := ccs.repo.ExistsByTitleExcludingID(id, updatedStrategicObjective.TitleEn, updatedStrategicObjective.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := ccs.repo.Update(id, updatedStrategicObjective)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingStrategicObjective, updatedStrategicObjective, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := ccs.activityLogSvc.LogActivity(
			logModels.UpdateAction,
			models.StrategicObjective{}.TableName(),
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

func (ccs *StrategicObjectiveService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := ccs.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = ccs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := ccs.activityLogSvc.LogBulkActivities(
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

func (ccs *StrategicObjectiveService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := ccs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ccs.activityLogSvc.LogBulkActivities(
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

func (ccs *StrategicObjectiveService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := ccs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ccs.activityLogSvc.LogBulkActivities(
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
