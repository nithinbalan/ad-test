package services

import (
	"jethur/grc/models"
	"jethur/grc/tprm/dtos"
	"jethur/grc/tprm/repositories"
)

type TPRMSettingServiceInterface interface {
	GetTPRMSetting() (*models.TPRMSetting, error)
	UpdateTPRMSetting(currentUserID uint, formDTO dtos.TprmSettingFormDTO) error
}

type TPRMSettingService struct {
	repo repositories.TPRMSettingRepositoryInterface
}

func NewTPRMSettingService(repo repositories.TPRMSettingRepositoryInterface) TPRMSettingServiceInterface {
	return &TPRMSettingService{
		repo: repo,
	}
}

func (tss *TPRMSettingService) GetTPRMSetting() (*models.TPRMSetting, error) {
	return tss.repo.GetTPRMSetting()
}

func (tss *TPRMSettingService) UpdateTPRMSetting(currentUserID uint, formDTO dtos.TprmSettingFormDTO) error {
	setting, err := tss.repo.GetTPRMSetting()
	if err != nil {
		return err
	}
	setting.IsExport = formDTO.IsExport
	setting.IsExternalVendor = formDTO.IsExternalVendor
	setting.UpdatedBYID = &currentUserID

	err = tss.repo.Update(*setting)
	if err != nil {
		return err
	}

	return err
}
