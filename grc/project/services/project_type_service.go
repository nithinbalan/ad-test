package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/project/repositories"
	"jethur/utility"
	"log"
)

// Interface
type ProjectTypeServiceInterface interface {
	GetAll(sort, sortBy, keyword, status string) ([]*models.ProjectType, error)
	GetByID(id uint) (*models.ProjectType, error)
	Create(currentUserID uint, pt *models.ProjectType) (*models.ProjectType, error)
	Update(id uint, currentUserID uint, pt *models.ProjectType) (*models.ProjectType, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Implementation
type ProjectTypeService struct {
	repo           repositories.ProjectTypeRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewProjectTypeService(repo repositories.ProjectTypeRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) ProjectTypeServiceInterface {
	return &ProjectTypeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (pts *ProjectTypeService) GetAll(sort, sortBy, keyword, status string) ([]*models.ProjectType, error) {
	return pts.repo.GetAll(sort, sortBy, keyword, status)
}

func (pts *ProjectTypeService) GetByID(id uint) (*models.ProjectType, error) {
	return pts.repo.GetByID(id)
}

func (pts *ProjectTypeService) Create(currentUserID uint, pt *models.ProjectType) (*models.ProjectType, error) {
	if exists, err := pts.repo.ExistsByTitle(pt.TitleEn, pt.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	created, err := pts.repo.Create(pt)
	if err != nil {
		return nil, err
	}

	go func() {
		err := pts.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ProjectType{}.TableName(),
			created.ID,
			created,
			currentUserID,
			nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return created, nil
}

func (pts *ProjectTypeService) Update(id uint, currentUserID uint, updated *models.ProjectType) (*models.ProjectType, error) {
	existing, err := pts.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updated.TitleEn != existing.TitleEn || updated.TitleAr != existing.TitleAr {
		if exists, err := pts.repo.ExistsByTitleExcludingID(id, updated.TitleEn, updated.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := pts.repo.Update(id, updated)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existing, updated, fieldsToCompare)
	if err != nil {
		log.Println("compare error:", err)
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("json error:", err)
	}

	go func() {
		err := pts.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ProjectType{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			log.Printf("Failed to log update activity: %v", err)
		}
	}()

	return result, nil
}

func (pts *ProjectTypeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint
	var deletedProjectTypes []map[string]interface{}

	for _, id := range ids {
		pt, err := pts.repo.GetByID(id)
		if err != nil {
			continue
		}

		err = pts.repo.Delete(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
		} else {
			successfulIDs = append(successfulIDs, id)
			deletedProjectTypes = append(deletedProjectTypes, map[string]interface{}{
				"id":       id,
				"title_en": pt.TitleEn,
				"title_ar": pt.TitleAr,
			})
		}
	}

	if len(successfulIDs) > 0 {
		go func() {
			err := pts.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ProjectType{}.TableName(),
				successfulIDs,
				deletedProjectTypes,
				currentUserID,
				nil,
			)
			if err != nil {
				log.Printf("Failed to log delete activity: %v", err)
			}
		}()
	}

	return failedIDs, nil
}

func (pts *ProjectTypeService) Activate(ids []uint, currentUserID uint) error {
	successfulIDs, err := pts.repo.Activate(ids)
	if len(successfulIDs) > 0 {
		go func() {
			err := pts.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ProjectType{}.TableName(),
				successfulIDs,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				log.Printf("Failed to log activate activity: %v", err)
			}
		}()
	}
	return err
}

func (pts *ProjectTypeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIDs, err := pts.repo.Deactivate(ids)
	if len(successfulIDs) > 0 {
		go func() {
			err := pts.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ProjectType{}.TableName(),
				successfulIDs,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				log.Printf("Failed to log deactivate activity: %v", err)
			}
		}()
	}
	return err
}
