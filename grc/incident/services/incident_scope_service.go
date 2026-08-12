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

type IncidentScopeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentScope, error)
	GetByID(id uint) (*models.IncidentScope, error)
	Create(currentUserID uint, scope *models.IncidentScope) (*models.IncidentScope, error)
	Update(id uint, currentUserID uint, scope *models.IncidentScope) (*models.IncidentScope, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type IncidentScopeService struct {
	repo           repositories.IncidentScopeRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewIncidentScopeService(repo repositories.IncidentScopeRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) IncidentScopeServiceInterface {
	return &IncidentScopeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (iss *IncidentScopeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentScope, error) {
	return iss.repo.GetAll(sort, sortBy, keyword, status)
}

func (iss IncidentScopeService) GetByID(id uint) (*models.IncidentScope, error) {
	incidentScope, err := iss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return incidentScope, nil
}

func (iss *IncidentScopeService) Create(currentUserID uint, scope *models.IncidentScope) (*models.IncidentScope, error) {
	if exists, err := iss.repo.ExistsByTitle(scope.TitleEn, scope.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdScope, err := iss.repo.Create(scope)
	if err != nil {
		return nil, err
	}

	go func() {
		err := iss.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentScope{}.TableName(),
			createdScope.ID,
			createdScope,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdScope, nil
}

func (iss *IncidentScopeService) Update(id uint, currentUserID uint, updatedScope *models.IncidentScope) (*models.IncidentScope, error) {
	existingScope, err := iss.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedScope.TitleEn != existingScope.TitleEn || updatedScope.TitleAr != existingScope.TitleAr {
		if exists, err := iss.repo.ExistsByTitleExcludingID(id, updatedScope.TitleEn, updatedScope.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := iss.repo.Update(id, updatedScope)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Type", "Status"}
	changes, err := utility.CompareFields(existingScope, updatedScope, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := iss.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IncidentScope{}.TableName(),
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

func (iss *IncidentScopeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedIncidentScope []map[string]interface{}
	for _, id := range ids {

		incidentScope, err := iss.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = iss.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedIncidentScope = append(deletedIncidentScope, map[string]interface{}{
			"id":       id,
			"title_en": incidentScope.TitleEn,
			"title_ar": incidentScope.TitleAr,
		})

	}

	if len(successfulIds) > 0 {
		go func() {
			err := iss.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.IncidentScope{}.TableName(),
				successfulIds,
				deletedIncidentScope,
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

func (iss *IncidentScopeService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := iss.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := iss.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.IncidentScope{}.TableName(),
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

func (iss *IncidentScopeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := iss.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := iss.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.IncidentScope{}.TableName(),
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
