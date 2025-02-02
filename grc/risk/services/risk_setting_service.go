package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
	"jethur/utility"
	"log"

	"github.com/mitchellh/mapstructure"
)

type RiskSettingsServiceInterface interface {
	Get() (*dtos.RiskSettingDTO, error)
	Update(currentUserID uint, setting *models.RiskSetting) (*models.RiskSetting, error)
}

type RiskSettingsService struct {
	repo           repositories.RiskSettingsRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewRiskSettingsService(repo repositories.RiskSettingsRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) RiskSettingsServiceInterface {
	return &RiskSettingsService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (rss *RiskSettingsService) Get() (*dtos.RiskSettingDTO, error) {
	var riskSettingDTO dtos.RiskSettingDTO

	riskSetting, err := rss.repo.Get()
	if err != nil {
		return nil, err
	}

	err = mapstructure.Decode(riskSetting, &riskSettingDTO)
	if err != nil {
		return nil, err
	}

	return &riskSettingDTO, nil
}

func (rss *RiskSettingsService) Update(currentUserID uint, updatedSetting *models.RiskSetting) (*models.RiskSetting, error) {
	existingSetting, err := rss.repo.Get()
	if err != nil {
		return nil, err
	}

	result, err := rss.repo.Update(updatedSetting)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"ReviewRoleID"}
	changes, err := utility.CompareFields(existingSetting, updatedSetting, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	// Log the update activity
	go func() {
		err := rss.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskSetting{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			fmt.Printf("Failed to log update activity: %v\n", err)
		}
	}()

	return result, nil
}
