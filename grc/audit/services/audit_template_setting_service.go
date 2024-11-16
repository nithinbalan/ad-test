package services

import (
	"errors"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type AuditTemplateSettingServiceInterface interface {
	GetAll(params *auditDTOs.AuditTemplateSettingListQueryParamsDTO) (*auditDTOs.AuditTemplateSettingListDTO, error)
	GetByID(id uint) (*models.AuditTemplateSetting, error)
	Update(id uint, currentUserID uint, updatedTemplateSettingDTO *auditDTOs.AuditTemplateSettingFormDTO) error
}

// Service and implementations
type AuditTemplateSettingService struct {
	repo           repositories.AuditTemplateSettingRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewAuditTemplateSettingService(repo repositories.AuditTemplateSettingRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) AuditTemplateSettingServiceInterface {
	return &AuditTemplateSettingService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (atss *AuditTemplateSettingService) GetAll(params *auditDTOs.AuditTemplateSettingListQueryParamsDTO) (*auditDTOs.AuditTemplateSettingListDTO, error) {
	auditTemplateSettings, err := atss.repo.GetAll(params)
	if err != nil {
		return nil, err
	}

	// Return the retrieved data
	return auditTemplateSettings, nil
}

func (atss AuditTemplateSettingService) GetByID(id uint) (*models.AuditTemplateSetting, error) {
	auditTemplateSetting, err := atss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditTemplateSetting, nil
}

func (atss *AuditTemplateSettingService) Update(id uint, currentUserID uint, updatedTemplateSettingDTO *auditDTOs.AuditTemplateSettingFormDTO) error {
	// Begin a database transaction
	tx := atss.repo.BeginTransaction()

	// Retrieve the existing AuditTemplateSetting by ID
	existingTemplateSetting, err := atss.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	// Helper function to handle img updates
	handleDocumentUpdate := func(doc *auditDTOs.AuditTemplateSettingDocumentFormDTO, templateType, storageFolder string, urlField **string) error {
		if doc.TemplateType == templateType && doc.IsNewFile != nil && *doc.IsNewFile {
			// Generate URL and update field in the existing template setting
			url := utility.GenerateDocumentURL("audit", "template-settings", storageFolder, doc.FileName)

			*urlField = &url // Update the pointer to reflect in the model field

			// Define storage path
			storagePath := filepath.Join("audit", "template-settings", storageFolder)

			// Store file
			if err := utility.StoreFile(doc.FileName, storagePath); err != nil {
				return err
			}
		}
		return nil
	}

	// Process documents
	for _, doc := range updatedTemplateSettingDTO.Documents {
		if err := handleDocumentUpdate(doc, "logo", "logo", &existingTemplateSetting.LogoImagePath); err != nil {
			tx.Rollback()
			return err
		}
		if err := handleDocumentUpdate(doc, "commencement-letter", "commencement-letter-header-bg", &existingTemplateSetting.CommencementLetterHeaderBGImagePath); err != nil {
			tx.Rollback()
			return err
		}
		if err := handleDocumentUpdate(doc, "draft-report", "draft-report-cover-page-bg", &existingTemplateSetting.DraftReportCoverPageBGImagePath); err != nil {
			tx.Rollback()
			return err
		}
		if err := handleDocumentUpdate(doc, "final-report", "final-report-cover-page-bg", &existingTemplateSetting.FinalReportCoverPageBGImagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update the common fields
	existingTemplateSetting.UpdatedBYID = &currentUserID
	now := time.Now()
	existingTemplateSetting.UpdatedAt = &now

	// Call repository to update the record
	if err := atss.repo.Update(existingTemplateSetting, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	tx.Commit()

	// Log activity in the background
	go func() {
		err := atss.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditTemplateSetting{}.TableName(),
			id,
			updatedTemplateSettingDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditTemplateSetting{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
