package services

import (
	"errors"
	"fmt"
	"log"

	"jethur/grc/incident/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentImpactServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentImpact, error)
	GetByID(id uint) (*models.IncidentImpact, error)
	Create(currentUserID uint, impact *models.IncidentImpact) (*models.IncidentImpact, error)
	Update(id uint, currentUserID uint, impact *models.IncidentImpact) (*models.IncidentImpact, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type IncidentImpactService struct {
	repo           repositories.IncidentImpactRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewIncidentImpactService(repo repositories.IncidentImpactRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) IncidentImpactServiceInterface {
	return &IncidentImpactService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (iis *IncidentImpactService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentImpact, error) {
	return iis.repo.GetAll(sort, sortBy, keyword, status)
}

func (iis IncidentImpactService) GetByID(id uint) (*models.IncidentImpact, error) {
	incidentImpact, err := iis.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return incidentImpact, nil
}

func (iis *IncidentImpactService) Create(currentUserID uint, impact *models.IncidentImpact) (*models.IncidentImpact, error) {
	if exists, err := iis.repo.ExistsByTitle(impact.TitleEn, impact.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdImpact, err := iis.repo.Create(impact)
	if err != nil {
		return nil, err
	}

	go func() {
		err := iis.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentImpact{}.TableName(),
			createdImpact.ID,
			createdImpact,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdImpact, nil
}

func (iis *IncidentImpactService) Update(id uint, currentUserID uint, updatedImpact *models.IncidentImpact) (*models.IncidentImpact, error) {
	existingImpact, err := iis.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedImpact.TitleEn != existingImpact.TitleEn || updatedImpact.TitleAr != existingImpact.TitleAr {
		if exists, err := iis.repo.ExistsByTitleExcludingID(id, updatedImpact.TitleEn, updatedImpact.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := iis.repo.Update(id, updatedImpact)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Type", "Description", "Status"}
	changes, err := utility.CompareFields(existingImpact, updatedImpact, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := iis.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IncidentImpact{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()
	return result, err
}

func (iis *IncidentImpactService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedIncidentImpact []map[string]interface{}
	for _, id := range ids {

		incidentImpact, err := iis.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = iis.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedIncidentImpact = append(deletedIncidentImpact, map[string]interface{}{
			"id":       id,
			"title_en": incidentImpact.TitleEn,
			"title_ar": incidentImpact.TitleAr,
		})

	}

	if len(successfulIds) > 0 {
		go func() {
			err := iis.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.IncidentImpact{}.TableName(),
				successfulIds,
				deletedIncidentImpact,
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

func (iis *IncidentImpactService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := iis.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := iis.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.IncidentImpact{}.TableName(),
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

func (iis *IncidentImpactService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := iis.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := iis.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.IncidentImpact{}.TableName(),
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
