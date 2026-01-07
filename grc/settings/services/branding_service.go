package services

import (
	"errors"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BrandingServiceInterface interface {
	GetByID(id uint) (*models.Branding, error)
	GetFullDetails() (*dtos.BrandingDetailDTO, error)
	Update(id uint, currentUserID uint, brandingFormDTO *dtos.BrandingUpdateFormDTO) error
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, nil)
		}
		return err
	}

	now := time.Now()

	existing.LoginPageTitle = brandingFormDTO.LoginPageTitle
	existing.LoginPageSubtitle = brandingFormDTO.LoginPageSubtitle
	existing.LoginPageBgImagePath = &brandingFormDTO.LoginPageBgImagePath
	existing.OrgLogoImagePath = &brandingFormDTO.OrgLogoImagePath
	existing.FaviconImagePath = &brandingFormDTO.FaviconImagePath
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
			"general_settings",
			"branding",
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
