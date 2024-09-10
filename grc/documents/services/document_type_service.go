package services

import (
	"errors"
	"fmt"
	"jethur/grc/documents/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type DocumentTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, docmentType string, status string) ([]*models.DocumentType, error)
	GetByID(id uint) (*models.DocumentType, error)
	Create(currentUserID uint, documentType *models.DocumentType) (*models.DocumentType, error)
	Update(id uint, currentUserID uint, documentType *models.DocumentType) (*models.DocumentType, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type DocumentTypeService struct {
	repo           repositories.DocumentTypeRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewDocumentTypeService(repo repositories.DocumentTypeRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) DocumentTypeServiceInterface {
	return &DocumentTypeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (dts *DocumentTypeService) GetAll(sort string, sortBy string, keyword string, docmentType string, status string) ([]*models.DocumentType, error) {
	return dts.repo.GetAll(sort, sortBy, keyword, docmentType, status)
}

func (dts *DocumentTypeService) GetByID(id uint) (*models.DocumentType, error) {
	documentType, err := dts.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return documentType, nil
}

func (dts *DocumentTypeService) Create(currentUserID uint, documentType *models.DocumentType) (*models.DocumentType, error) {
	tx := dts.repo.BeginTransaction()

	if exists, err := dts.repo.ExistsByTitle(documentType.TitleEn, documentType.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdDocumentType, err := dts.repo.Create(documentType)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	workflow := &models.Workflow{
		DocumentTypeID: &createdDocumentType.ID,
		Type:           "document",
		CreatedByID:    currentUserID,
		Status:         "active",
	}

	err = dts.repo.CreateWorkflow(workflow, tx) // Create the workflow passing the transaction
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	go func() {
		err := dts.activityLogSvc.LogActivity(
			models.CreateAction,
			models.DocumentType{}.TableName(),
			createdDocumentType.ID,
			createdDocumentType,
			currentUserID,
			nil,
		)

		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdDocumentType, nil
}

func (dts *DocumentTypeService) Update(id uint, currentUserID uint, updatedDocumentType *models.DocumentType) (*models.DocumentType, error) {
	existingDocumentType, err := dts.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updatedDocumentType.TitleEn != existingDocumentType.TitleEn || updatedDocumentType.TitleAr != existingDocumentType.TitleAr {
		if exists, err := dts.repo.ExistsByTitleExcludingID(id, updatedDocumentType.TitleEn, updatedDocumentType.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	result, err := dts.repo.Update(id, updatedDocumentType)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingDocumentType, updatedDocumentType, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := dts.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.DocumentType{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}()
	return result, err
}

func (dts *DocumentTypeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {
		_, err := dts.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = dts.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := dts.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.DocumentType{}.TableName(),
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

	return failedIds, nil
}

func (dts *DocumentTypeService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := dts.repo.Activate(ids)
	if len(successfulIds) > 0 {
		go func() {
			err := dts.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.DocumentType{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk activate activity: %v\n", err)
			}
		}()
	}
	return err
}

func (dts *DocumentTypeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := dts.repo.Deactivate(ids)
	if len(successfulIds) > 0 {
		go func() {
			err := dts.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.DocumentType{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk deactivate activity: %v\n", err)
			}
		}()
	}
	return err
}
