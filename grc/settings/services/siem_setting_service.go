package services

import (
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"
)

type SIEMSettingServiceInterface interface {
	GetFirst() (*models.SIEMSetting, error)
	Update(currentUserID uint, dto dtos.SIEMSettingDTO) error
	Delete(currentUserID uint) error
}

type SIEMSettingService struct {
	repo           repositories.SIEMSettingRepositoryInterface
	integrationSvc IntegrationServiceInterface
}

func NewSIEMSettingService(repo repositories.SIEMSettingRepositoryInterface, integrationSvc IntegrationServiceInterface) SIEMSettingServiceInterface {
	return &SIEMSettingService{
		repo:           repo,
		integrationSvc: integrationSvc,
	}
}

func (sss *SIEMSettingService) GetFirst() (*models.SIEMSetting, error) {
	return sss.repo.GetFirst()
}

func (sss *SIEMSettingService) Update(currentUserID uint, dto dtos.SIEMSettingDTO) error {
	integration, err := sss.integrationSvc.GetByType("siem")
	if err != nil {
		return err
	}

	_, err = sss.integrationSvc.EnableIntegrationSettings(integration, currentUserID)
	if err != nil {
		return utility.NewForbiddenError("failed to enable SIEM integration")
	}

	existingSiemSettings, err := sss.repo.GetFirst()
	if err != nil {
		return err
	}

	if existingSiemSettings != nil {
		// Update
		existingSiemSettings.IntegrationID = integration.ID
		existingSiemSettings.SyslogServer = dto.SyslogServer
		existingSiemSettings.SyslogPort = dto.SyslogPort
		existingSiemSettings.SyslogProtocol = dto.SyslogProtocol
		existingSiemSettings.SyslogTimezone = dto.SyslogTimezone
		existingSiemSettings.SyslogDatetimeFormat = dto.SyslogDatetimeFormat
		existingSiemSettings.UpdatedByID = &currentUserID

		return sss.repo.Update(existingSiemSettings)
	}

	// Create new
	newSetting := models.SIEMSetting{
		IntegrationID:        integration.ID,
		SyslogServer:         dto.SyslogServer,
		SyslogPort:           dto.SyslogPort,
		SyslogProtocol:       dto.SyslogProtocol,
		SyslogTimezone:       dto.SyslogTimezone,
		SyslogDatetimeFormat: dto.SyslogDatetimeFormat,
		UpdatedByID:          &currentUserID,
	}

	return sss.repo.Create(&newSetting)
}

func (sss *SIEMSettingService) Delete(currentUserID uint) error {
	setting, err := sss.repo.GetFirst()
	if err != nil {
		return err
	}
	if setting == nil {
		return err
	}

	integration, err := sss.integrationSvc.GetByType("siem")
	if err != nil {
		return err
	}

	// DISABLE integration settings
	_, err = sss.integrationSvc.DisableIntegrationSettings(integration, currentUserID)
	if err != nil {
		return err
	}

	return sss.repo.Delete(setting.ID)
}
