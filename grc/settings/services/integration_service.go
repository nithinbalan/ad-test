package services

import (
	"errors"
	"fmt"
	activityLogServices "jethur/grc/logs/services"
	integrationCategoryServices "jethur/grc/masters/services"
	"jethur/grc/models"

	// integrationDTOs "jethur/grc/settings/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type IntegrationServiceInterface interface {
	GetAll(params *settingsDTOs.IntegrationListQueryParamsDTO) ([]*settingsDTOs.IntegrationDTO, error)
	GetCredentailByID(id uint) (*models.IntegrationCredential, error)
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
	integrationCredential, err := rs.repo.GetCredentailByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return integrationCredential, nil
}

func (is *IntegrationService) EnableIntegration(currentUserID uint, credential *settingsDTOs.IntegrationCredentialFormDTO) (uint, error) {

	integrationCredential := &models.IntegrationCredential{
		IntegrationID:  credential.IntegrationID,
		ApiKey:         credential.ApiKey,
		ApiSecret:      credential.ApiSecret,
		RedirectUrl:    credential.RedirectUrl,
		TokenType:      credential.TokenType,
		Scope:          credential.Scope,
		AccessToken:    credential.AccessToken,
		ExpireAt:       credential.ExpireAt,
		RefreshToken:   credential.RefreshToken,
		AdditionalInfo: credential.AdditionalInfo,
		IsActive:       credential.IsActive,
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

func (is *IntegrationService) UpdateCredential(id uint, currentUserID uint, updatedIntegrationCredential *settingsDTOs.IntegrationCredentialFormDTO) (uint, error) {
	existingIntegrationCredential, err := is.repo.GetCredentailByID(id)

	if err != nil {
		return 0, err
	}

	existingIntegrationCredential.IntegrationID = updatedIntegrationCredential.IntegrationID
	existingIntegrationCredential.ApiKey = updatedIntegrationCredential.ApiKey
	existingIntegrationCredential.ApiSecret = updatedIntegrationCredential.ApiSecret
	existingIntegrationCredential.RedirectUrl = updatedIntegrationCredential.RedirectUrl
	existingIntegrationCredential.TokenType = updatedIntegrationCredential.TokenType
	existingIntegrationCredential.Scope = updatedIntegrationCredential.Scope
	existingIntegrationCredential.AccessToken = updatedIntegrationCredential.AccessToken
	existingIntegrationCredential.ExpireAt = updatedIntegrationCredential.ExpireAt
	existingIntegrationCredential.RefreshToken = updatedIntegrationCredential.RefreshToken
	existingIntegrationCredential.AdditionalInfo = updatedIntegrationCredential.AdditionalInfo
	existingIntegrationCredential.IsActive = updatedIntegrationCredential.IsActive
	existingIntegrationCredential.UpdatedByID = &currentUserID

	err = is.repo.Update(id, existingIntegrationCredential)
	if err != nil {
		return 0, err
	}

	go func() {
		err := is.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IntegrationCredential{}.TableName(),
			id,
			updatedIntegrationCredential,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}()
	return existingIntegrationCredential.IntegrationID, nil
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
