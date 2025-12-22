package services

import (
	"encoding/json"
	"fmt"
	"io"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"
	"log"
	"net/http"
	"strings"
	"time"
)

type TenableSettingsServiceInterface interface {
	GetTenableSettings() (*models.TenableSettings, error)
	SaveTenableSettings(currentUserID uint, dto dtos.UpdateTenableSettingsDTO) error
	DeleteTenableSettings(currentUserID uint) error
	TestConnection(input dtos.TestTenableConnectionDTO) error
	CheckSavedConnection() error
}

type tenableSettingsService struct {
	repo           repositories.TenableSettingsRepository
	integrationsvc IntegrationServiceInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewTenableSettingsService(repo repositories.TenableSettingsRepository, integrationsvc IntegrationServiceInterface, activityLogSvc logServices.ActivityLogServiceInterface) TenableSettingsServiceInterface {
	return &tenableSettingsService{repo: repo, integrationsvc: integrationsvc, activityLogSvc: activityLogSvc}
}

func (tss *tenableSettingsService) GetTenableSettings() (*models.TenableSettings, error) {
	return tss.repo.GetFirst()
}

func (tss *tenableSettingsService) SaveTenableSettings(currentUserID uint, dto dtos.UpdateTenableSettingsDTO) error {
	integration, err := tss.integrationsvc.GetByType("tenable")
	if err != nil {
		return err
	}

	_, err = tss.integrationsvc.EnableIntegrationSettings(integration, currentUserID)
	if err != nil {
		return fmt.Errorf("failed to enable Tenable integration: %w", err)
	}

	existingSetting, err := tss.repo.GetFirst()
	if err != nil {
		return err
	}

	if existingSetting != nil {
		existingSetting.BaseURL = dto.BaseURL
		existingSetting.AccessKey = dto.AccessKey
		existingSetting.SecretKey = dto.SecretKey
		existingSetting.SyncFrequencyHours = dto.SyncFrequencyHours
		existingSetting.LookbackDays = dto.LookbackDays
		existingSetting.Status = "active"
		existingSetting.UpdatedAt = time.Now().UTC()

		err = tss.repo.Update(existingSetting)
		if err != nil {
			return err
		}

		go func() {
			err := tss.activityLogSvc.LogActivity(
				models.UpdateAction,
				"tenable_settings",
				existingSetting.ID,
				existingSetting,
				currentUserID,
				"", "", "", nil,
			)
			if err != nil {
				log.Printf("Failed to log Tenable settings update: %v", err)
			}
		}()
		return nil

	}

	newSetting := models.TenableSettings{
		BaseURL:            dto.BaseURL,
		AccessKey:          dto.AccessKey,
		SecretKey:          dto.SecretKey,
		SyncFrequencyHours: dto.SyncFrequencyHours,
		LookbackDays:       dto.LookbackDays,
		Status:             "active",
		CreatedAt:          time.Now().UTC(),
	}

	err = tss.repo.Create(&newSetting)

	if err != nil {
		return err
	}

	go func() {
		err := tss.activityLogSvc.LogActivity(
			models.CreateAction,
			"tenable_settings",
			newSetting.ID,
			newSetting,
			currentUserID,
			"setting", "general_setting", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log Tenable settings create: %v", err)
		}
	}()

	return nil
}

func (tss *tenableSettingsService) DeleteTenableSettings(currentUserID uint) error {
	setting, err := tss.repo.GetFirst()
	if err != nil {
		return err
	}
	if setting == nil {
		return fmt.Errorf("no Tenable settings found")
	}

	integration, err := tss.integrationsvc.GetByType("tenable")
	if err != nil {
		return err
	}

	_, err = tss.integrationsvc.DisableIntegrationSettings(integration, currentUserID)
	if err != nil {
		return err
	}

	err = tss.repo.Delete(setting)
	if err != nil {
		return err
	}

	go func() {
		err := tss.activityLogSvc.LogActivity(
			models.DeleteAction,
			"tenable_settings",
			setting.ID,
			setting,
			currentUserID,
			"setting", "general_setting", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log Tenable settings delete: %v", err)
		}
	}()

	return nil

}

func (tss *tenableSettingsService) TestConnection(input dtos.TestTenableConnectionDTO) error {

	client := &http.Client{Timeout: 20 * time.Second}
	url := fmt.Sprintf("%s/rest/system", input.BaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return utility.NewTenableConnectionError(
			fmt.Errorf("unable to create request: %w", err),
		)
	}

	req.Header.Set("x-apikey",
		fmt.Sprintf("accessKey=%s; secretKey=%s",
			input.AccessKey, input.SecretKey,
		),
	)

	resp, err := client.Do(req)
	if err != nil {
		return utility.NewTenableConnectionError(
			fmt.Errorf("connection failed: %w", err),
		)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		var terr dtos.TenableTestErrorResponse
		if err := json.Unmarshal(body, &terr); err == nil && terr.ErrorMsg != "" {
			return utility.NewTenableConnectionError(
				fmt.Errorf("%s", strings.TrimSpace(terr.ErrorMsg)),
			)
		}

		return utility.NewTenableConnectionError(
			fmt.Errorf("tenable returned %d: %s", resp.StatusCode, string(body)),
		)
	}

	return nil
}

func (tss *tenableSettingsService) CheckSavedConnection() error {

	integration, err := tss.integrationsvc.GetByType("tenable")
	if err != nil {
		return utility.NewTenableConnectionError(
			fmt.Errorf("failed to find integration: %w", err),
		)
	}

	if !integration.IsEnable {
		return utility.NewTenableConnectionError(
			fmt.Errorf("integration 'tenable' is disabled"),
		)
	}

	setting, err := tss.repo.GetFirst()
	if err != nil {
		return utility.NewTenableConnectionError(
			fmt.Errorf("failed to fetch saved Tenable settings: %w", err),
		)
	}

	if setting == nil {
		return utility.NewTenableConnectionError(
			fmt.Errorf("no Tenable settings found"),
		)
	}

	dto := dtos.TestTenableConnectionDTO{
		BaseURL:   setting.BaseURL,
		AccessKey: setting.AccessKey,
		SecretKey: setting.SecretKey,
	}

	return tss.TestConnection(dto)
}
