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

type RiskTemplateSettingServiceInterface interface {
	GetAll(params *settingsDTOs.RiskTemplateSettingListQueryParamsDTO) ([]*settingsDTOs.RiskTemplateSettingListDTO, error)
	GetByID(id uint) (*models.RiskTemplateSetting, error)
	GetFirst() (*models.RiskTemplateSetting, error)
	Update(id uint, currentUserID uint, updatedTemplateSettingDTO *settingsDTOs.RiskTemplateSettingFormDTO) error
	GetImagePath(templateType string) (string, error)
}

// Service and implementations
type RiskTemplateSettingService struct {
	repo           repositories.RiskTemplateSettingRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewRiskTemplateSettingService(repo repositories.RiskTemplateSettingRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) RiskTemplateSettingServiceInterface {
	return &RiskTemplateSettingService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (rtss *RiskTemplateSettingService) GetAll(params *settingsDTOs.RiskTemplateSettingListQueryParamsDTO) ([]*settingsDTOs.RiskTemplateSettingListDTO, error) {
	riskTemplateSettings, err := rtss.repo.GetAll(params)
	if err != nil {
		return nil, err
	}

	return riskTemplateSettings, nil
}

func (rtss RiskTemplateSettingService) GetByID(id uint) (*models.RiskTemplateSetting, error) {
	riskTemplateSetting, err := rtss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskTemplateSetting, nil
}

func (rtss RiskTemplateSettingService) GetFirst() (*models.RiskTemplateSetting, error) {
	riskTemplateSetting, err := rtss.repo.GetFirst()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, nil)
		}
		return nil, err
	}
	return riskTemplateSetting, nil
}
func (rtss *RiskTemplateSettingService) Update(id uint, currentUserID uint, updatedTemplateSettingDTO *settingsDTOs.RiskTemplateSettingFormDTO) error {
	// Begin a database transaction
	tx := rtss.repo.BeginTransaction()

	// Retrieve the existing RiskTemplateSetting by ID
	existingTemplateSetting, err := rtss.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	// Helper function to handle document updates
	handleDocumentUpdate := func(doc *settingsDTOs.RiskTemplateSettingDocumentFormDTO, storageFolder string, urlField **string) error {
		if doc != nil && doc.IsNewFile != nil && *doc.IsNewFile {
			// Generate URL and update field in the existing template setting
			url := utility.GenerateDocumentURL("risk", "template-settings", storageFolder, doc.FileName)

			*urlField = &url // Update the pointer to reflect in the model field

			// Define storage path
			storagePath := filepath.Join("risk", "template-settings", storageFolder)

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
	if err := handleDocumentUpdate(updatedTemplateSettingDTO.CoverImage, "cover-page-bg", &existingTemplateSetting.CoverPageBGImagePath); err != nil {
		tx.Rollback()
		return err
	}

	// Update the common fields
	existingTemplateSetting.UpdatedBYID = &currentUserID
	now := time.Now()
	existingTemplateSetting.UpdatedAt = &now

	// Call repository to update the record
	if err := rtss.repo.Update(existingTemplateSetting, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	tx.Commit()

	// Log activity in the background
	go func() {
		err := rtss.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskTemplateSetting{}.TableName(),
			id,
			updatedTemplateSettingDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.RiskTemplateSetting{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (rtss *RiskTemplateSettingService) GetImagePath(templateType string) (string, error) {
	// Retrieve the RiskTemplateSetting by ID
	riskTemplateSetting, err := rtss.repo.GetFirst()
	if err != nil {
		return "", err
	}

	// Map templateType to the corresponding storage path
	var storagePath *string
	switch templateType {
	case "logo":
		storagePath = riskTemplateSetting.LogoImagePath
	case "cover-page-bg":
		storagePath = riskTemplateSetting.CoverPageBGImagePath
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
