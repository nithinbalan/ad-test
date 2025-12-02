package services

import (
	"encoding/json"
	"fmt"
	"io"
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"
	"net/http"
	"strings"
	"time"
)

type TenableSettingsService interface {
	GetTenableSettings() (*models.TenableSettings, error)
	SaveTenableSettings(currentUserID uint, dto dtos.UpdateTenableSettingsDTO) error
	DeleteTenableSettings(currentUserID uint) error
	TestConnection(input dtos.TestTenableConnectionDTO) error
}

type tenableSettingsService struct {
	repo           repositories.TenableSettingsRepository
	integrationsvc IntegrationServiceInterface
}

func NewTenableSettingsService(repo repositories.TenableSettingsRepository, integrationsvc IntegrationServiceInterface) TenableSettingsService {
	return &tenableSettingsService{repo: repo, integrationsvc: integrationsvc}
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

		return tss.repo.Update(existingSetting)
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

	return tss.repo.Create(&newSetting)
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

	return tss.repo.Delete(setting)
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
