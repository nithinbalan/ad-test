package services

import (
	"fmt"
	"jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
	activityLogModels "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	"jethur/utility"
	"log"
)

type ComplianceServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceRegister, error)
	GetByID(id uint) (*models.ComplianceRegister, error)
	Create(currentUserID uint, compliance *models.ComplianceRegister) (*models.ComplianceRegister, error)
	Update(id uint, currentUserID uint, compliance *models.ComplianceRegister) (*models.ComplianceRegister, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type ComplianceService struct {
	repo           repositories.ComplianceRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewComplianceService(repo repositories.ComplianceRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) ComplianceServiceInterface {
	return &ComplianceService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ccs *ComplianceService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceRegister, error) {
	return ccs.repo.GetAll(sort, sortBy, keyword, status)
}

func (ccs ComplianceService) GetByID(id uint) (*models.ComplianceRegister, error) {
	return ccs.repo.GetByID(id)
}

func (ccs *ComplianceService) Create(currentUserID uint, compliance *models.ComplianceRegister) (*models.ComplianceRegister, error) {

	if exists, err := ccs.repo.ExistsByTitle(compliance.TitleEn, compliance.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdCompliance, err := ccs.repo.Create(compliance)
	if err != nil {
		return nil, err
	}

	go func() {
		err := ccs.activityLogSvc.LogActivity(
			activityLogModels.CreateAction,
			models.ComplianceRegister{}.TableName(),
			compliance.ID,
			createdCompliance,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdCompliance, nil
}

func (ccs *ComplianceService) Update(id uint, currentUserID uint, updatedCompliance *models.ComplianceRegister) (*models.ComplianceRegister, error) {
	existingCompliance, err := ccs.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedCompliance.TitleEn != existingCompliance.TitleEn || updatedCompliance.TitleAr != existingCompliance.TitleAr {
		if exists, err := ccs.repo.ExistsByTitleExcludingID(id, updatedCompliance.TitleEn, updatedCompliance.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := ccs.repo.Update(id, updatedCompliance)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingCompliance, updatedCompliance, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := ccs.activityLogSvc.LogActivity(
			activityLogModels.UpdateAction,
			models.ComplianceRegister{}.TableName(),
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

func (ccs *ComplianceService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := ccs.repo.GetByID(id)

		if err != nil {
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
				activityLogModels.DeleteAction, // Make sure this action constant is defined appropriately
				models.ComplianceRegister{}.TableName(),
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

func (ccs *ComplianceService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := ccs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ccs.activityLogSvc.LogBulkActivities(
				activityLogModels.ActivateAction,
				models.ComplianceRegister{}.TableName(),
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

func (ccs *ComplianceService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := ccs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ccs.activityLogSvc.LogBulkActivities(
				activityLogModels.DeactivateAction,
				models.ComplianceRegister{}.TableName(),
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
