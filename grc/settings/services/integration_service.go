package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"log"

	settingsDTOs "jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IntegrationServiceInterface interface {
	GetAll(params *settingsDTOs.IntegrationListQueryParamsDTO) ([]*settingsDTOs.IntegrationDTO, error)
	GetByID(id uint) (*models.Integration, error)
	GetFullDetailsByID(integrationID uint, tx ...utility.DatabaseTransaction) (*settingsDTOs.IntegrationDetailDTO, error)
	EnableIntegration(integrationID uint, currentUserID uint, integrtaionCredentialFormDTO *settingsDTOs.IntegrationCredentialFormDTO) (uint, error)
	UpdateCredential(id uint, currentUserID uint, integrationCredential *settingsDTOs.IntegrationCredentialFormDTO) (uint, error)
	DisableIntegration(id uint, currentUserID uint) (uint, error)
	// UpdateIntegration(id uint)
}

// Service and implementations
type IntegrationService struct {
	repo                   repositories.IntegrationRepositoryInterface
	activityLogSvc         logServices.ActivityLogServiceInterface
	integrationCategorySvc masterServices.IntegrationCategoryServiceInterface
}

func NewIntegrationService(repo repositories.IntegrationRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface, integrationCategorySvc masterServices.IntegrationCategoryServiceInterface) IntegrationServiceInterface {
	return &IntegrationService{
		repo:                   repo,
		activityLogSvc:         activityLogSvc,
		integrationCategorySvc: integrationCategorySvc,
	}
}

func (is *IntegrationService) GetAll(params *settingsDTOs.IntegrationListQueryParamsDTO) ([]*settingsDTOs.IntegrationDTO, error) {
	return is.repo.GetAll(params)
}

func (rs IntegrationService) GetByID(id uint) (*models.Integration, error) {
	integration, err := rs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return integration, nil
}
func (is IntegrationService) GetFullDetailsByID(integrationId uint, tx ...utility.DatabaseTransaction) (*settingsDTOs.IntegrationDetailDTO, error) {
	// Initialize your DTO
	var integrationDTO settingsDTOs.IntegrationDetailDTO

	integration, err := is.repo.GetFullDetailsByID(integrationId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, integrationId)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(integration, &integrationDTO)
	if err != nil {
		return nil, err
	}

	return &integrationDTO, nil
}

func (is *IntegrationService) EnableIntegration(integrationID uint, currentUserID uint, credentialFormDTO *settingsDTOs.IntegrationCredentialFormDTO) (uint, error) {
	tx := is.repo.BeginTransaction()
	existingIntegration, err := is.GetByID(integrationID)

	if err != nil {
		return 0, err
	}

	existingIntegration.IsEnable = true

	id, err := is.repo.Update(integrationID, existingIntegration)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	isActive := true
	integrationCredential := &models.IntegrationCredential{
		IntegrationID:  integrationID,
		ApiKey:         credentialFormDTO.ApiKey,
		ApiSecret:      credentialFormDTO.ApiSecret,
		RedirectUrl:    credentialFormDTO.RedirectUrl,
		TokenType:      credentialFormDTO.TokenType,
		Scope:          credentialFormDTO.Scope,
		AccessToken:    credentialFormDTO.AccessToken,
		ExpireAt:       credentialFormDTO.ExpireAt,
		RefreshToken:   credentialFormDTO.RefreshToken,
		AdditionalInfo: credentialFormDTO.AdditionalInfo,
		IsActive:       &isActive,
		CreatedByID:    currentUserID,
	}
	err = is.repo.CreateCredential(integrationCredential)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	go func() {
		err := is.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IntegrationCredential{}.TableName(),
			integrationCredential.ID,
			integrationCredential,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return id, nil
}

func (is *IntegrationService) UpdateCredential(id uint, currentUserID uint, credentialFormDTO *settingsDTOs.IntegrationCredentialFormDTO) (uint, error) {
	tx := is.repo.BeginTransaction()
	if _, err := is.GetByID(id); err != nil {
		return 0, err
	}

	// Update existing Integration Credential
	existingIntegration, err := is.repo.GetFullDetailsByID(id)
	if err != nil {
		return 0, err
	}
	if existingIntegration == nil {
		return 0, errors.New("existing integration not found")
	}

	existingIntegrationCredential := existingIntegration.IntegrationCredential

	existingIntegrationCredential.ApiKey = credentialFormDTO.ApiKey
	existingIntegrationCredential.ApiSecret = credentialFormDTO.ApiSecret
	existingIntegrationCredential.RedirectUrl = credentialFormDTO.RedirectUrl
	existingIntegrationCredential.TokenType = credentialFormDTO.TokenType
	existingIntegrationCredential.Scope = credentialFormDTO.Scope
	existingIntegrationCredential.AccessToken = credentialFormDTO.AccessToken
	existingIntegrationCredential.ExpireAt = credentialFormDTO.ExpireAt
	existingIntegrationCredential.RefreshToken = credentialFormDTO.RefreshToken
	existingIntegrationCredential.AdditionalInfo = credentialFormDTO.AdditionalInfo
	existingIntegrationCredential.IsActive = credentialFormDTO.IsActive
	existingIntegrationCredential.UpdatedByID = &currentUserID

	err = is.repo.UpdateCredential(existingIntegrationCredential, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	go func() {
		err := is.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IntegrationCredential{}.TableName(),
			id,
			credentialFormDTO,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}()
	return id, nil
}

func (is *IntegrationService) DisableIntegration(integrationID uint, currentUserID uint) (uint, error) {
	tx := is.repo.BeginTransaction()

	existingIntegration, err := is.GetByID(integrationID)

	if err != nil {
		return 0, err
	}

	existingIntegration.IsEnable = false

	id, err := is.repo.Update(integrationID, existingIntegration)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	integration, err := is.repo.GetFullDetailsByID(id)
	if err != nil {
		return 0, err
	}

	// Delete the integrationCredential

	err = is.repo.DeleteCredential(integration.IntegrationCredential.ID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	go func() {
		err := is.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.IntegrationCredential{}.TableName(),
			id,
			nil,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return id, nil
}
