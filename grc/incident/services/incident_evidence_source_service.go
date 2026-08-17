package services

import (
	"fmt"
	"log"

	"jethur/grc/incident/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
)

type IncidentEvidenceSourceServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentEvidenceSource, error)
	GetByID(id uint) (*models.IncidentEvidenceSource, error)
	Create(currentUserID uint, evidenceSource *models.IncidentEvidenceSource) (*models.IncidentEvidenceSource, error)
	Update(id uint, currentUserID uint, evidenceSource *models.IncidentEvidenceSource) (*models.IncidentEvidenceSource, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type IncidentEvidenceSourceService struct {
	repo           repositories.IncidentEvidenceSourceRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewIncidentEvidenceSourceService(repo repositories.IncidentEvidenceSourceRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) IncidentEvidenceSourceServiceInterface {
	return &IncidentEvidenceSourceService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (iess *IncidentEvidenceSourceService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentEvidenceSource, error) {
	return iess.repo.GetAll(sort, sortBy, keyword, status)
}

func (iess IncidentEvidenceSourceService) GetByID(id uint) (*models.IncidentEvidenceSource, error) {
	return iess.repo.GetByID(id)
}

func (iess *IncidentEvidenceSourceService) Create(currentUserID uint, evidenceSource *models.IncidentEvidenceSource) (*models.IncidentEvidenceSource, error) {
	if exists, err := iess.repo.ExistsByTitle(evidenceSource.TitleEn, evidenceSource.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdEvidenceSource, err := iess.repo.Create(evidenceSource)
	if err != nil {
		return nil, err
	}

	go func() {
		err := iess.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentEvidenceSource{}.TableName(),
			createdEvidenceSource.ID,
			createdEvidenceSource,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdEvidenceSource, nil
}

func (iess *IncidentEvidenceSourceService) Update(id uint, currentUserID uint, updatedEvidenceSource *models.IncidentEvidenceSource) (*models.IncidentEvidenceSource, error) {
	existingEvidenceSource, err := iess.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedEvidenceSource.TitleEn != existingEvidenceSource.TitleEn || updatedEvidenceSource.TitleAr != existingEvidenceSource.TitleAr {
		if exists, err := iess.repo.ExistsByTitleExcludingID(id, updatedEvidenceSource.TitleEn, updatedEvidenceSource.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := iess.repo.Update(id, updatedEvidenceSource)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingEvidenceSource, updatedEvidenceSource, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := iess.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IncidentEvidenceSource{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()
	return result, err
}

func (iess *IncidentEvidenceSourceService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedIncidentEvidenceSource []map[string]interface{}
	for _, id := range ids {

		incidentEvidenceSource, err := iess.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = iess.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedIncidentEvidenceSource = append(deletedIncidentEvidenceSource, map[string]interface{}{
			"id":       id,
			"title_en": incidentEvidenceSource.TitleEn,
			"title_ar": incidentEvidenceSource.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := iess.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.IncidentEvidenceSource{}.TableName(),
				successfulIds,
				deletedIncidentEvidenceSource,
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

func (iess *IncidentEvidenceSourceService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := iess.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := iess.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.IncidentEvidenceSource{}.TableName(),
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

func (iess *IncidentEvidenceSourceService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := iess.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := iess.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.IncidentEvidenceSource{}.TableName(),
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
