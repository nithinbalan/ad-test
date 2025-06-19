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
	Delete() error
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
	if !integration.IsEnable {
		return utility.NewForbiddenError("SIEM integration is disabled")
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

func (sss *SIEMSettingService) Delete() error {
	setting, err := sss.repo.GetFirst()
	if err != nil {
		return err
	}
	if setting == nil {
		return err
	}
	return sss.repo.Delete(setting.ID)
}
