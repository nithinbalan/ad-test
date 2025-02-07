package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
	"jethur/utility"
	"log"
)

type RequiredDocumentTitleServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RequiredDocumentTitle, error)
	GetByID(id uint) (*models.RequiredDocumentTitle, error)
	Create(currentUserID uint, document *models.RequiredDocumentTitle) (*models.RequiredDocumentTitle, error)
	Update(id uint, currentUserID uint, document *models.RequiredDocumentTitle) (*models.RequiredDocumentTitle, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type RequiredDocumentTitleService struct {
	repo           repositories.RequiredDocumentTitleRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewRequiredDocumentTitleService(repo repositories.RequiredDocumentTitleRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) RequiredDocumentTitleServiceInterface {
	return &RequiredDocumentTitleService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (rdts *RequiredDocumentTitleService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RequiredDocumentTitle, error) {
	return rdts.repo.GetAll(sort, sortBy, keyword, status)
}

func (rdts RequiredDocumentTitleService) GetByID(id uint) (*models.RequiredDocumentTitle, error) {
	return rdts.repo.GetByID(id)
}

func (rdts *RequiredDocumentTitleService) Create(currentUserID uint, document *models.RequiredDocumentTitle) (*models.RequiredDocumentTitle, error) {
	if exists, err := rdts.repo.ExistsByTitle(document.TitleEn, document.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdDocument, err := rdts.repo.Create(document)
	if err != nil {
		return nil, err
	}

	go func() {
		err := rdts.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RequiredDocumentTitle{}.TableName(),
			document.ID,
			createdDocument,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdDocument, nil
}

func (rdts *RequiredDocumentTitleService) Update(id uint, currentUserID uint, updatedDocument *models.RequiredDocumentTitle) (*models.RequiredDocumentTitle, error) {
	existingDocument, err := rdts.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedDocument.TitleEn != existingDocument.TitleEn || updatedDocument.TitleAr != existingDocument.TitleAr {
		if exists, err := rdts.repo.ExistsByTitleExcludingID(id, updatedDocument.TitleEn, updatedDocument.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := rdts.repo.Update(id, updatedDocument)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingDocument, updatedDocument, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := rdts.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RequiredDocumentTitle{}.TableName(),
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

func (rdts *RequiredDocumentTitleService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedDocuments []map[string]interface{}

	for _, id := range ids {
		document, err := rdts.repo.GetByID(id)
		if err != nil {
			continue
		}

		err = rdts.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedDocuments = append(deletedDocuments, map[string]interface{}{
			"id":       id,
			"title_en": document.TitleEn,
			"title_ar": document.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := rdts.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.RequiredDocumentTitle{}.TableName(),
				successfulIds,
				deletedDocuments,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (rdts *RequiredDocumentTitleService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := rdts.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := rdts.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.RequiredDocumentTitle{}.TableName(),
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

func (rdts *RequiredDocumentTitleService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := rdts.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := rdts.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.RequiredDocumentTitle{}.TableName(),
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
