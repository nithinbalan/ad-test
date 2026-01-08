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

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BrandingServiceInterface interface {
	GetByID(id uint) (*models.Branding, error)
	GetFullDetails() (*dtos.BrandingDetailDTO, error)
	Update(id uint, currentUserID uint, brandingFormDTO *dtos.BrandingUpdateFormDTO) error
	GetPreviewImage(previewType string) (string, error)
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

func (bs *BrandingService) GetFullDetails() (*dtos.BrandingDetailDTO, error) {
	var brandingDTO dtos.BrandingDetailDTO

	branding, err := bs.repo.GetFullDetails()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, nil)
		}
		return nil, err
	}

	err = mapstructure.Decode(branding, &brandingDTO)
	if err != nil {
		return nil, err
	}

	return &brandingDTO, nil
}

func (bs BrandingService) GetByID(id uint) (*models.Branding, error) {
	branding, err := bs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, nil)
		}

		return nil, err
	}

	return branding, nil
}

func (bs *BrandingService) Update(id uint, currentUserID uint, brandingFormDTO *dtos.BrandingUpdateFormDTO) error {
	tx := bs.repo.BeginTransaction()

	existing, err := bs.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	now := time.Now()

	existing.LoginPageTitleEn = brandingFormDTO.LoginPageTitleEn
	existing.LoginPageTitleAr = brandingFormDTO.LoginPageTitleAr
	existing.LoginPageSubtitleEn = brandingFormDTO.LoginPageSubtitleEn
	existing.LoginPageSubtitleAr = brandingFormDTO.LoginPageSubtitleAr
	existing.LoginPageBgImagePath = &brandingFormDTO.LoginPageBgImagePath
	existing.OrgLogoImagePath = &brandingFormDTO.OrgLogoImagePath
	existing.FaviconImagePath = &brandingFormDTO.FaviconImagePath
	existing.LoaderGifPath = &brandingFormDTO.LoaderGifPath
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

func (bs *BrandingService) GetPreviewImage(previewType string) (string, error) {
	// Retrieve the branding details
	branding, err := bs.repo.GetFullDetails()
	if err != nil {
		return "", err
	}

	var storagePath *string
	
	// Map previewType to the corresponding storage path
	switch previewType {
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