package services

import (
	"encoding/xml"
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

type QualysSettingsServiceInterface interface {
	GetQualysSettings() (*models.QualysSettings, error)
	SaveQualysSettings(currentUserID uint, dto dtos.UpdateQualysSettingsDTO) error
	DeleteQualysSettings(currentUserID uint) error
	TestConnection(input dtos.TestQualysConnectionDTO) error
	CheckSavedConnection() error
}

type qualysSettingsService struct {
	repo           repositories.QualysSettingsRepository
	integrationsvc IntegrationServiceInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewQualysSettingsService(
	repo repositories.QualysSettingsRepository,
	integrationsvc IntegrationServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) QualysSettingsServiceInterface {
	return &qualysSettingsService{repo: repo, integrationsvc: integrationsvc, activityLogSvc: activityLogSvc}
}

func (qss *qualysSettingsService) GetQualysSettings() (*models.QualysSettings, error) {
	return qss.repo.GetFirst()
}

func (qss *qualysSettingsService) SaveQualysSettings(currentUserID uint, dto dtos.UpdateQualysSettingsDTO) error {

	if dto.LookbackDays <= 0 || dto.LookbackDays > 90 {
		return fmt.Errorf("lookback days must be between 1 and 90 days")
	}

	integration, err := qss.integrationsvc.GetByType("qualys")
	if err != nil {
		return err
	}

	_, err = qss.integrationsvc.EnableIntegrationSettings(integration, currentUserID)
	if err != nil {
		return fmt.Errorf("failed to enable Qualys integration: %w", err)
	}

	existingSetting, err := qss.repo.GetFirst()
	if err != nil {
		return err
	}

	if existingSetting != nil {
		now := time.Now()
		existingSetting.PlatformURL = dto.PlatformURL
		existingSetting.Username = dto.Username
		existingSetting.Password = dto.Password
		existingSetting.LookbackDays = dto.LookbackDays
		existingSetting.Status = "active"
		existingSetting.UpdatedByID = &currentUserID
		existingSetting.UpdatedAt = &now

		if err := qss.repo.Update(existingSetting); err != nil {
			return err
		}

		go func() {
			err := qss.activityLogSvc.LogActivity(
				models.UpdateAction,
				"qualys_settings",
				existingSetting.ID,
				existingSetting,
				currentUserID,
				"setting",
				"general_setting",
				"",
				nil,
			)
			if err != nil {
				log.Printf("Failed to log Qualys settings update: %v", err)
			}
		}()
		return nil
	}

	newSetting := models.QualysSettings{
		PlatformURL:  dto.PlatformURL,
		Username:     dto.Username,
		Password:     dto.Password,
		LookbackDays: dto.LookbackDays,
		Status:       "active",
		CreatedByID:  currentUserID,
		CreatedAt:    time.Now().UTC(),
	}

	if err := qss.repo.Create(&newSetting); err != nil {
		return err
	}

	go func() {
		err := qss.activityLogSvc.LogActivity(
			models.CreateAction,
			"qualys_settings",
			newSetting.ID,
			newSetting,
			currentUserID,
			"setting",
			"general_setting",
			"", nil,
		)
		if err != nil {
			log.Printf("Failed to log Qualys settings create: %v", err)
		}
	}()
	return nil
}

func (qss *qualysSettingsService) DeleteQualysSettings(currentUserID uint) error {

	setting, err := qss.repo.GetFirst()
	if err != nil {
		return err
	}
	if setting == nil {
		return fmt.Errorf("no Qualys settings found")
	}

	integration, err := qss.integrationsvc.GetByType("qualys")
	if err != nil {
		return err
	}

	_, err = qss.integrationsvc.DisableIntegrationSettings(integration, currentUserID)
	if err != nil {
		return err
	}

	if err := qss.repo.Delete(setting); err != nil {
		return err
	}

	go func() {
		err := qss.activityLogSvc.LogActivity(
			models.DeleteAction,
			"qualys_settings",
			setting.ID,
			setting,
			currentUserID,
			"setting", "general_setting", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log Qualys settings delete: %v", err)
		}
	}()

	return nil
}

func (qss *qualysSettingsService) TestConnection(input dtos.TestQualysConnectionDTO) error {

	client := &http.Client{Timeout: 20 * time.Second}

	url := fmt.Sprintf("%s/api/4.0/fo/asset/host/vm/detection/", strings.TrimRight(input.PlatformURL, "/"))

	payload := strings.NewReader("action=list&truncation_limit=1")

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return utility.NewQualysConnectionError(
			fmt.Errorf("unable to create request: %w", err),
		)
	}

	req.SetBasicAuth(input.Username, input.Password)
	req.Header.Set("X-Requested-With", "GRC-Test")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return utility.NewQualysConnectionError(
			fmt.Errorf("connection failed: %w", err),
		)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {

    var xmlResp struct {
        Response struct {
            Code string `xml:"CODE"`
            Text string `xml:"TEXT"`
        } `xml:"RESPONSE"`
    }

    if err := xml.Unmarshal(body, &xmlResp); err == nil {
        if strings.TrimSpace(xmlResp.Response.Text) != "" {
            return utility.NewQualysConnectionError(
                fmt.Errorf("%s", strings.TrimSpace(xmlResp.Response.Text)),
            )
        }
    }

    return utility.NewQualysConnectionError(
        fmt.Errorf("authentication failed (HTTP %d)", resp.StatusCode),
    )
}


	return nil
}

func (qss *qualysSettingsService) CheckSavedConnection() error {

	integration, err := qss.integrationsvc.GetByType("qualys")
	if err != nil {
		return utility.NewQualysConnectionError(
			fmt.Errorf("failed to find integration: %w", err),
		)
	}

	if !integration.IsEnable {
		return utility.NewQualysConnectionError(
			fmt.Errorf("integration 'qualys' is disabled"),
		)
	}

	setting, err := qss.repo.GetFirst()
	if err != nil {
		return utility.NewQualysConnectionError(
			fmt.Errorf("failed to fetch saved Qualys settings: %w", err),
		)
	}

	if setting == nil {
		return utility.NewQualysConnectionError(fmt.Errorf("no Qualys settings found"))
	}

	dto := dtos.TestQualysConnectionDTO{
		PlatformURL: setting.PlatformURL,
		Username:    setting.Username,
		Password:    setting.Password,
	}

	return qss.TestConnection(dto)
}