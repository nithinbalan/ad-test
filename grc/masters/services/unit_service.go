package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
	"jethur/utility"
	"log"
)

type UnitServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Unit, error)
	GetByID(id uint) (*models.Unit, error)
	Create(currentUserID uint, unit *models.Unit) (*models.Unit, error)
	Update(id uint, currentUserID uint, unit *models.Unit) (*models.Unit, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type UnitService struct {
	repo           repositories.UnitRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewUnitService(repo repositories.UnitRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) UnitServiceInterface {
	return &UnitService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (us *UnitService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Unit, error) {
	return us.repo.GetAll(sort, sortBy, keyword, status)
}

func (us UnitService) GetByID(id uint) (*models.Unit, error) {
	return us.repo.GetByID(id)
}

func (us *UnitService) Create(currentUserID uint, unit *models.Unit) (*models.Unit, error) {

	if err := us.validateUnitForm(unit, nil); err != nil {
		return nil, err
	}

	createdUnit, err := us.repo.Create(unit)
	if err != nil {
		return nil, err
	}

	go func() {
		err := us.activityLogSvc.LogActivity(
			models.CreateAction,
			models.Unit{}.TableName(),
			unit.ID,
			createdUnit,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdUnit, nil
}

func (us *UnitService) Update(id uint, currentUserID uint, updatedUnit *models.Unit) (*models.Unit, error) {
	existingUnit, err := us.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if err := us.validateUnitForm(updatedUnit, &id); err != nil {
		return nil, err
	}

	result, err := us.repo.Update(id, updatedUnit)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "CodeEn", "CodeAr", "Type", "Status"}
	changes, err := utility.CompareFields(existingUnit, updatedUnit, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := us.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.Unit{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}()
	return result, err
}

func (us *UnitService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := us.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = us.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := us.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.Unit{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (us *UnitService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := us.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := us.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.Unit{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}

func (us *UnitService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := us.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := us.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.Unit{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				// You might want to log this error or handle it appropriately
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}

func (us *UnitService) validateUnitForm(unitForm *models.Unit, id *uint) error {

	if id == nil {
		if exists, err := us.repo.ExistsByTitle(unitForm.TitleEn, unitForm.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
		if exists, err := us.repo.ExistsByCode(unitForm.CodeEn, unitForm.CodeAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the code must be unique"), "unique", "code")
		}

	} else {
		if exists, err := us.repo.ExistsByTitleExcludingID(*id, unitForm.TitleEn, unitForm.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
		if exists, err := us.repo.ExistsByCodeExcludingID(*id, unitForm.CodeEn, unitForm.CodeAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the code must be unique"), "unique", "code")
		}

	}

	return nil
}
