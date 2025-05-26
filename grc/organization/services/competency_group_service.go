package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/organization/repositories"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type CompetencyGroupServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CompetencyGroup, error)
	GetByID(id uint) (*models.CompetencyGroup, error)
	Create(currentUserID uint, competencyGroup *models.CompetencyGroup) (*models.CompetencyGroup, error)
	Update(id uint, currentUserID uint, competencyGroup *models.CompetencyGroup) (*models.CompetencyGroup, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type CompetencyGroupService struct {
	repo           repositories.CompetencyGroupRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewCompetencyGroupService(repo repositories.CompetencyGroupRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) CompetencyGroupServiceInterface {
	return &CompetencyGroupService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (pgs *CompetencyGroupService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CompetencyGroup, error) {
	return pgs.repo.GetAll(sort, sortBy, keyword, status)
}

func (pgs *CompetencyGroupService) GetByID(id uint) (*models.CompetencyGroup, error) {
	competencyGroup, err := pgs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return competencyGroup, nil
}

func (pgs *CompetencyGroupService) Create(currentUserID uint, competencyGroup *models.CompetencyGroup) (*models.CompetencyGroup, error) {
	if exists, err := pgs.repo.ExistsByTitle(competencyGroup.TitleEn, competencyGroup.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdCompetencyGroup, err := pgs.repo.Create(competencyGroup)
	if err != nil {
		return nil, err
	}

	go func() {
		err := pgs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.CompetencyGroup{}.TableName(),
			createdCompetencyGroup.ID,
			createdCompetencyGroup,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdCompetencyGroup, nil
}

func (pgs *CompetencyGroupService) Update(id uint, currentUserID uint, updatedCompetencyGroup *models.CompetencyGroup) (*models.CompetencyGroup, error) {
	existingCompetencyGroup, err := pgs.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedCompetencyGroup.TitleEn != existingCompetencyGroup.TitleEn || updatedCompetencyGroup.TitleAr != existingCompetencyGroup.TitleAr {
		if exists, err := pgs.repo.ExistsByTitleExcludingID(id, updatedCompetencyGroup.TitleEn, updatedCompetencyGroup.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := pgs.repo.Update(id, updatedCompetencyGroup)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingCompetencyGroup, updatedCompetencyGroup, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := pgs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.CompetencyGroup{}.TableName(),
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

func (pgs *CompetencyGroupService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := pgs.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = pgs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := pgs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.CompetencyGroup{}.TableName(),
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

func (pgs *CompetencyGroupService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := pgs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := pgs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.CompetencyGroup{}.TableName(),
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

func (pgs *CompetencyGroupService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := pgs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := pgs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.CompetencyGroup{}.TableName(),
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
