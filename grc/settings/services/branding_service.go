package services

import (
	"errors"
	"fmt"
	"jethur/config"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type BrandingServiceInterface interface {
	GetAll(params *dtos.BrandingListQueryParamsDTO) ([]*dtos.BrandingListDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Branding, error)
	GetFirst() (*models.Branding, error)
	Update(id uint, currentUserID uint, brandingFormDTO *dtos.BrandingFormDTO) error
	GetBrandingPreviewImage(previewType string) (string, error)
}

type BrandingService struct {
	repo           repositories.BrandingRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewBrandingService(repo repositories.BrandingRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) BrandingServiceInterface {
	return &BrandingService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (bs *BrandingService) GetAll(params *dtos.BrandingListQueryParamsDTO) ([]*dtos.BrandingListDTO, error) {
	branding, err := bs.repo.GetAll(params)
	if err != nil {
		return nil, err
	}

	return branding, nil
}

func (bs BrandingService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Branding, error) {
	branding, err := bs.repo.GetByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, nil)
		}

		return nil, err
	}

	return branding, nil
}

func (bs BrandingService) GetFirst() (*models.Branding, error) {
	branding, err := bs.repo.GetFirst()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, nil)
		}
		return nil, err
	}
	return branding, nil
}

func (bs *BrandingService) Update(id uint, currentUserID uint, brandingFormDTO *dtos.BrandingFormDTO) error {
	tx := bs.repo.BeginTransaction()

	existing, err := bs.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	handleBrandingUpdate := func(doc *dtos.BrandingDocumentFormDTO, storageFolder string, urlField **string) error {
		if doc != nil && doc.IsNewFile != nil && *doc.IsNewFile {

			var oldRelativePath *string
			if *urlField != nil && **urlField != "" {
				tmp := **urlField
				oldRelativePath = &tmp
			}

			//define storage path
			storagePath := filepath.Join("branding", "branding-settings", storageFolder)
			if err := utility.StoreFile(doc.FileName, storagePath); err != nil {
				return err
			}

			// Generate URL for the new document,
			url := utility.GenerateDocumentURL("branding", "branding-settings", storageFolder, doc.FileName)

			*urlField = &url

			// Remove old file if exists
			if oldRelativePath != nil {
				oldDir := filepath.Dir(*oldRelativePath)
				oldFile := filepath.Base(*oldRelativePath)
				_ = utility.RemoveFile(oldFile, oldDir)
			}
		}

		return nil
	}

	if err := handleBrandingUpdate(brandingFormDTO.LoginPageBgImage, "login-page-bg", &existing.LoginPageBgImagePath); err != nil {
		tx.Rollback()
		return err
	}

	if err := handleBrandingUpdate(brandingFormDTO.OrgLogoImage, "org-logo", &existing.OrgLogoImagePath); err != nil {
		tx.Rollback()
		return err
	}

	if err := handleBrandingUpdate(brandingFormDTO.FaviconImage, "favicon", &existing.FaviconImagePath); err != nil {
		tx.Rollback()
		return err
	}

	if err := handleBrandingUpdate(brandingFormDTO.LoaderGif, "loader", &existing.LoaderGifPath); err != nil {
		tx.Rollback()
		return err
	}

	now := time.Now()

	existing.LoginPageTitleEn = brandingFormDTO.LoginPageTitleEn
	existing.LoginPageTitleAr = brandingFormDTO.LoginPageTitleAr
	existing.LoginPageSubtitleEn = brandingFormDTO.LoginPageSubtitleEn
	existing.LoginPageSubtitleAr = brandingFormDTO.LoginPageSubtitleAr
	existing.UpdatedBYID = &currentUserID
	existing.UpdatedAt = &now

	if err := bs.repo.Update(existing, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log activity in the background
	go func() {
		err := bs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.Branding{}.TableName(),
			existing.ID,
			brandingFormDTO,
			currentUserID,
			"setting",
			"general_setting",
			"",
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.Branding{}.TableName(),
				"recordID":  existing.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (bs *BrandingService) GetBrandingPreviewImage(previewType string) (string, error) {
	// Retrieve the branding details
	branding, err := bs.GetFirst()
	if err != nil {
		return "", err
	}

	var storagePath *string

	// Map previewType to the corresponding storage path
	switch previewType {
	case "login-page-bg":
		storagePath = branding.LoginPageBgImagePath
	case "loader":
		storagePath = branding.LoaderGifPath
	case "org-logo":
		storagePath = branding.OrgLogoImagePath
	case "favicon":
		storagePath = branding.FaviconImagePath
	default:
		return "", fmt.Errorf("invalid preview type: %s", previewType)
	}

	// Check if the storagePath pointer is nil
	if storagePath == nil || *storagePath == "" {
		return "", fmt.Errorf("Path not set for preview type: %s", previewType)
	}

	// Build the full file path
	basePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(basePath, *storagePath)

	// Check if the file exists
	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file not found")
		}
		return "", err
	}

	return fullPath, nil
}
