package services

import (
	"errors"
	"fmt"
	"jethur/config"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	settingsDTOs "jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type AuditTemplateSettingServiceInterface interface {
	GetAll(params *settingsDTOs.AuditTemplateSettingListQueryParamsDTO) ([]*settingsDTOs.AuditTemplateSettingListDTO, error)
	GetByID(id uint) (*models.AuditTemplateSetting, error)
	GetFirst() (*models.AuditTemplateSetting, error)
	Update(id uint, currentUserID uint, updatedTemplateSettingDTO *settingsDTOs.AuditTemplateSettingFormDTO) error
	GetImagePath(id uint, templateType string) (string, error)
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

func (atss *AuditTemplateSettingService) GetAll(params *settingsDTOs.AuditTemplateSettingListQueryParamsDTO) ([]*settingsDTOs.AuditTemplateSettingListDTO, error) {
	auditTemplateSettings, err := atss.repo.GetAll(params)
	if err != nil {
		return nil, err
	}

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

func (atss AuditTemplateSettingService) GetFirst() (*models.AuditTemplateSetting, error) {
	auditTemplateSetting, err := atss.repo.GetFirst()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, nil)
		}
		return nil, err
	}
	return auditTemplateSetting, nil
}
func (atss *AuditTemplateSettingService) Update(id uint, currentUserID uint, updatedTemplateSettingDTO *settingsDTOs.AuditTemplateSettingFormDTO) error {
	// Begin a database transaction
	tx := atss.repo.BeginTransaction()

	// Retrieve the existing AuditTemplateSetting by ID
	existingTemplateSetting, err := atss.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	// Helper function to handle document updates
	handleDocumentUpdate := func(doc *settingsDTOs.AuditTemplateSettingDocumentFormDTO, storageFolder string, urlField **string) error {
		if doc != nil && doc.IsNewFile != nil && *doc.IsNewFile {
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

	// Process documents for each template type
	if err := handleDocumentUpdate(updatedTemplateSettingDTO.Logo, "logo", &existingTemplateSetting.LogoImagePath); err != nil {
		tx.Rollback()
		return err
	}
	if err := handleDocumentUpdate(updatedTemplateSettingDTO.CommencementLetter, "commencement-letter-header-bg", &existingTemplateSetting.CommencementLetterHeaderBGImagePath); err != nil {
		tx.Rollback()
		return err
	}
	if err := handleDocumentUpdate(updatedTemplateSettingDTO.DraftReport, "draft-report-cover-page-bg", &existingTemplateSetting.DraftReportCoverPageBGImagePath); err != nil {
		tx.Rollback()
		return err
	}
	if err := handleDocumentUpdate(updatedTemplateSettingDTO.FinalReport, "final-report-cover-page-bg", &existingTemplateSetting.FinalReportCoverPageBGImagePath); err != nil {
		tx.Rollback()
		return err
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

func (atss *AuditTemplateSettingService) GetImagePath(id uint, templateType string) (string, error) {
	// Retrieve the AuditTemplateSetting by ID
	auditTemplateSetting, err := atss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", utility.NewResourceNotFoundError(err, id)
		}
		return "", err
	}

	// Map templateType to the corresponding storage path
	var storagePath *string
	switch templateType {
	case "logo":
		storagePath = auditTemplateSetting.LogoImagePath
	case "commencement_letter":
		storagePath = auditTemplateSetting.CommencementLetterHeaderBGImagePath
	case "draft_report":
		storagePath = auditTemplateSetting.DraftReportCoverPageBGImagePath
	case "final_report":
		storagePath = auditTemplateSetting.FinalReportCoverPageBGImagePath
	default:
		return "", fmt.Errorf("invalid template type: %s", templateType)
	}
	// Check if the storagePath pointer is nil
	if storagePath == nil || *storagePath == "" {
		return "", fmt.Errorf("no path set for template type: %s", templateType)
	}
	// Build the full file path
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, *storagePath)

	// Check if the file exists
	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file not found")
		}
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
