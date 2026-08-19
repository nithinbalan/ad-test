package services

import (
	"fmt"
	"log"

	"jethur/grc/incident/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
)

type IncidentEvidenceTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentEvidenceType, error)
	GetByID(id uint) (*models.IncidentEvidenceType, error)
	Create(currentUserID uint, evidenceType *models.IncidentEvidenceType) (*models.IncidentEvidenceType, error)
	Update(id uint, currentUserID uint, evidenceType *models.IncidentEvidenceType) (*models.IncidentEvidenceType, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type IncidentEvidenceTypeService struct {
	repo           repositories.IncidentEvidenceTypeRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewIncidentEvidenceTypeService(repo repositories.IncidentEvidenceTypeRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) IncidentEvidenceTypeServiceInterface {
	return &IncidentEvidenceTypeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (iets *IncidentEvidenceTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentEvidenceType, error) {
	return iets.repo.GetAll(sort, sortBy, keyword, status)
}

func (iets IncidentEvidenceTypeService) GetByID(id uint) (*models.IncidentEvidenceType, error) {
	return iets.repo.GetByID(id)
}

func (iets *IncidentEvidenceTypeService) Create(currentUserID uint, evidenceType *models.IncidentEvidenceType) (*models.IncidentEvidenceType, error) {
	if exists, err := iets.repo.ExistsByTitle(evidenceType.TitleEn, evidenceType.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdEvidenceType, err := iets.repo.Create(evidenceType)
	if err != nil {
		return nil, err
	}

	go func() {
		err := iets.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentEvidenceType{}.TableName(),
			createdEvidenceType.ID,
			createdEvidenceType,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdEvidenceType, nil
}

func (iets *IncidentEvidenceTypeService) Update(id uint, currentUserID uint, updatedEvidenceType *models.IncidentEvidenceType) (*models.IncidentEvidenceType, error) {
	existingEvidenceType, err := iets.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedEvidenceType.TitleEn != existingEvidenceType.TitleEn || updatedEvidenceType.TitleAr != existingEvidenceType.TitleAr {
		if exists, err := iets.repo.ExistsByTitleExcludingID(id, updatedEvidenceType.TitleEn, updatedEvidenceType.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := iets.repo.Update(id, updatedEvidenceType)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingEvidenceType, updatedEvidenceType, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := iets.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IncidentEvidenceType{}.TableName(),
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

func (iets *IncidentEvidenceTypeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedIncidentEvidenceType []map[string]interface{}
	for _, id := range ids {

		incidentEvidenceType, err := iets.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = iets.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedIncidentEvidenceType = append(deletedIncidentEvidenceType, map[string]interface{}{
			"id":       id,
			"title_en": incidentEvidenceType.TitleEn,
			"title_ar": incidentEvidenceType.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := iets.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.IncidentEvidenceType{}.TableName(),
				successfulIds,
				deletedIncidentEvidenceType,
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

func (iets *IncidentEvidenceTypeService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := iets.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := iets.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.IncidentEvidenceType{}.TableName(),
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

func (iets *IncidentEvidenceTypeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := iets.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := iets.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.IncidentEvidenceType{}.TableName(),
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
