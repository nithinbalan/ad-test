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
	// Step 1: Check integration
	integration, err := sss.integrationSvc.GetByType("siem")
	if err != nil {
		return err
	}
	if !integration.IsEnable {
		return utility.NewForbiddenError("SIEM integration is disabled")
	}

	// Step 2: Get existing or create new
	saved, err := sss.repo.GetFirst()
	if err != nil {
		return err
	}

	if saved != nil {
		// Update
		saved.IntegrationID = integration.ID
		saved.SyslogServer = dto.SyslogServer
		saved.SyslogPort = dto.SyslogPort
		saved.SyslogProtocol = dto.SyslogProtocol
		saved.SyslogTimezone = dto.SyslogTimezone
		saved.SyslogDatetimeFormat = dto.SyslogDatetimeFormat
		saved.UpdatedByID = &currentUserID

		return sss.repo.Update(saved)
	}

	// Create new
	newSetting := models.SIEMSetting{
		IntegrationID:        integration.ID,
		SyslogServer:         dto.SyslogServer,
		SyslogPort:           dto.SyslogPort,
		SyslogProtocol:       dto.SyslogProtocol,
		SyslogTimezone:       dto.SyslogTimezone,
		SyslogDatetimeFormat: dto.SyslogDatetimeFormat,
		CreatedByID:          currentUserID,
	}

	return sss.repo.Create(&newSetting)
}
