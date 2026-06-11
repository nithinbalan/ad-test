package services

import (
	"fmt"
	"log"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
)

type GovernanceSettingsServiceInterface interface {
	Get() (*dtos.GovernanceSettingDTO, error)
	Update(currentUserID uint, setting *models.GovernanceSetting) (*models.GovernanceSetting, error)
}

type GovernanceSettingsService struct {
	repo           repositories.GovernanceSettingsRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewGovernanceSettingsService(
	repo repositories.GovernanceSettingsRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) GovernanceSettingsServiceInterface {
	return &GovernanceSettingsService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (gss *GovernanceSettingsService) Get() (*dtos.GovernanceSettingDTO, error) {
	var governanceSettingDTO dtos.GovernanceSettingDTO

	governanceSetting, err := gss.repo.Get()
	if err != nil {
		return nil, err
	}

	err = mapstructure.Decode(governanceSetting, &governanceSettingDTO)
	if err != nil {
		return nil, err
	}

	return &governanceSettingDTO, nil
}

func (gss *GovernanceSettingsService) Update(currentUserID uint, updatedSetting *models.GovernanceSetting) (*models.GovernanceSetting, error) {

	existingSetting, err := gss.repo.Get()
	if err != nil {
		return nil, err
	}

	result, err := gss.repo.Update(updatedSetting)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{
		"ApproveDepartmentID",
	}

	changes, err := utility.CompareFields(existingSetting, updatedSetting, fieldsToCompare)
	if err != nil {
		log.Println("error:", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error:", err.Error())
	}

	go func() {
		err := gss.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.GovernanceSetting{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fmt.Printf("Failed to log update activity: %v\n", err)
		}
	}()

	return result, nil
}
