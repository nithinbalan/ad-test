package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/strategy/dtos"
	"jethur/grc/strategy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategyProfileSettingServiceInterface interface {
	GetAll(paraspss *dtos.StrategyProfileSettingListQueryParamsDTO) (*[]dtos.StrategyProfileSettingListDTO, error)
	GetFullDetailsByID(id uint) (*dtos.StrategyProfileSettingDetailDTO, error)
	EnableStrategyProfileSetting(id uint) (*models.StrategyProfileSetting, error)
	DisableStrategyProfileSetting(id uint) (*models.StrategyProfileSetting, error)
}

type StrategyProfileSettingService struct {
	repo repositories.StrategyProfileSettingRepositoryInterface
}

func NewStrategyProfileSettingService(repo repositories.StrategyProfileSettingRepositoryInterface) StrategyProfileSettingServiceInterface {
	return &StrategyProfileSettingService{
		repo: repo,
	}
}

func (spss *StrategyProfileSettingService) GetAll(params *dtos.StrategyProfileSettingListQueryParamsDTO) (*[]dtos.StrategyProfileSettingListDTO, error) {
	return spss.repo.GetAll(params)
}

func (spss *StrategyProfileSettingService) GetByID(id uint) (*models.StrategyProfileSetting, error) {
	strategyProfileSetting, err := spss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return strategyProfileSetting, nil
}

func (spss StrategyProfileSettingService) GetFullDetailsByID(strategyProfileSettingId uint) (*dtos.StrategyProfileSettingDetailDTO, error) {
	// Initialize your DTO
	var strategyProfileSettingDTO dtos.StrategyProfileSettingDetailDTO

	strategyProfileSetting, err := spss.repo.GetFullDetailsByID(strategyProfileSettingId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, strategyProfileSettingId)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(strategyProfileSetting, &strategyProfileSettingDTO)
	if err != nil {
		return nil, err
	}

	return &strategyProfileSettingDTO, nil
}

func (spss *StrategyProfileSettingService) EnableStrategyProfileSetting(id uint) (*models.StrategyProfileSetting, error) {

	// Enable the strategyProfileSetting
	strategyProfileSetting, err := spss.GetByID(id)
	if err != nil {
		return nil, err
	}
	strategyProfileSetting.IsEnable = true
	_, err = spss.repo.Update(strategyProfileSetting.ID, strategyProfileSetting)
	if err != nil {
		return nil, err
	}

	return strategyProfileSetting, nil
}

func (spss *StrategyProfileSettingService) DisableStrategyProfileSetting(id uint) (*models.StrategyProfileSetting, error) {

	strategyProfileSetting, err := spss.GetByID(id)
	if err != nil {
		return nil, err
	}

	strategyProfileSetting.IsEnable = false
	_, err = spss.repo.Update(strategyProfileSetting.ID, strategyProfileSetting)
	if err != nil {
		return nil, err
	}

	return strategyProfileSetting, nil
}
