package services

import (
	"errors"
	"fmt"
	activityLogServices "jethur/grc/logs/services"
	integrationCategoryServices "jethur/grc/masters/services"
	"jethur/grc/models"

	settingsDTOs "jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"
	"log"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IntegrationServiceInterface interface {
	GetAll(params *settingsDTOs.IntegrationListQueryParamsDTO) ([]*settingsDTOs.IntegrationDTO, error)
	GetCredentailByID(id uint) (*models.IntegrationCredential, error)
	GetCredentialFullDetailsByID(integrationID uint, tx ...utility.DatabaseTransaction) (*settingsDTOs.IntegrationCredentialDetailDTO, error)
	EnableIntegration(currentUserID uint, integrtaionCredentialFormDTO *settingsDTOs.IntegrationCredentialFormDTO) (uint, error)
	UpdateCredential(id uint, currentUserID uint, integrationCredential *settingsDTOs.IntegrationCredentialFormDTO) (uint, error)
	DisableIntegration(id uint, currentUserID uint) (uint, error)
}

// Service and implementations
type IntegrationService struct {
	repo                   repositories.IntegrationRepositoryInterface
	activityLogSvc         activityLogServices.ActivityLogServiceInterface
	integrationCategorySvc integrationCategoryServices.IntegrationCategoryServiceInterface
}

func NewIntegrationService(repo repositories.IntegrationRepositoryInterface, activityLogSvc activityLogServices.ActivityLogServiceInterface, integrationCategorySvc integrationCategoryServices.IntegrationCategoryServiceInterface) IntegrationServiceInterface {
	return &IntegrationService{
		repo:                   repo,
		activityLogSvc:         activityLogSvc,
		integrationCategorySvc: integrationCategorySvc,
	}
}

func (is *IntegrationService) GetAll(params *settingsDTOs.IntegrationListQueryParamsDTO) ([]*settingsDTOs.IntegrationDTO, error) {
	return is.repo.GetAll(params)
}

func (rs IntegrationService) GetCredentailByID(id uint) (*models.IntegrationCredential, error) {
	integrationCredential, err := rs.repo.GetCredentialByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return integrationCredential, nil
}
func (is IntegrationService) GetCredentialFullDetailsByID(integrationId uint, tx ...utility.DatabaseTransaction) (*settingsDTOs.IntegrationCredentialDetailDTO, error) {
	// Initialize your DTO
	var integrationCredentialDTO settingsDTOs.IntegrationCredentialDetailDTO

	integration, err := is.repo.GetCredentialFullDetailsByID(integrationId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, integrationId)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(integration, &integrationCredentialDTO)
	if err != nil {
		return nil, err
	}

	// integrationCredentialDTO.CreatedAt = integration.CreatedAt
	// integrationCredentialDTO.UpdatedAt = integration.UpdatedAt

	return &integrationCredentialDTO, nil
}

func (is *IntegrationService) EnableIntegration(currentUserID uint, credentialFormDTO *settingsDTOs.IntegrationCredentialFormDTO) (uint, error) {

	isActive := true
	integrationCredential := &models.IntegrationCredential{
		IntegrationID:  credentialFormDTO.IntegrationID,
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
	id, err := is.repo.Create(integrationCredential)
	if err != nil {
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
	if _, err := is.repo.GetByID(id); err != nil {
		return 0, err
	}

	// Update existing Integration Credential
	existingIntegrationCredential, err := is.repo.GetCredentialByID(id)
	if err != nil {
		return 0, err
	}
	if existingIntegrationCredential == nil {
		return 0, errors.New("existing integration credential not found")
	}

	// existingIntegrationCredential, err := is.repo.GetCredentailByID(id)

	// if err != nil {
	// 	return 0, err
	// }

	existingIntegrationCredential.IntegrationID = credentialFormDTO.IntegrationID
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

	err = is.repo.Update(existingIntegrationCredential, tx)
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

func (is *IntegrationService) DisableIntegration(id uint, currentUserID uint) (uint, error) {
	var successfulIds uint

	_, err := is.repo.GetByID(id)
	if err != nil {
		return 0, err
	}
	// Delete the integrationCredential
	integration, err := is.repo.GetFullDetailsByID(id)
	if err != nil {
		return 0, err
	}

	if len(integration.Credentials) > 0 {
		var deleteCredentialIds []*uint
		for _, update := range integration.Credentials {
			deleteCredentialIds = append(deleteCredentialIds, &update.ID)
		}
		if err := is.repo.Delete(id, deleteCredentialIds); err != nil {
			return 0, err
		}
	}
	successfulIds = id

	go func() {
		err := is.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.IntegrationCredential{}.TableName(),
			successfulIds,
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
