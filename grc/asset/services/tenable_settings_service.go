package services

import (
	"fmt"
	"io"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	"jethur/grc/models"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type TenableSettingsService interface {
	GetTenableSettings() (*dtos.TenableSettingsResponseDTO, error)
	UpdateTenableSettings(input dtos.UpdateTenableSettingsDTO) error
	ActivateSettings(id uint) error
	DeactivateSettings(id uint) error
	TestConnection(input dtos.TestTenableConnectionDTO) error
}

type tenableSettingsService struct {
	repo repositories.TenableSettingsRepository
}

func NewTenableSettingsService(repo repositories.TenableSettingsRepository) TenableSettingsService {
	return &tenableSettingsService{repo: repo}
}

func (tss *tenableSettingsService) getOrCreate() (*models.TenableSettings, error) {
	st, err := tss.repo.GetFirst()

	if err == gorm.ErrRecordNotFound {
		st = &models.TenableSettings{
			BaseURL:            "",
			AccessKey:          "",
			SecretKey:          "",
			SyncFrequencyHours: 24,
			LookbackDays:       30,
			Status:             "inactive",
		}
		if err := tss.repo.Create(st); err != nil {
			return nil, err
		}
		return st, nil
	}

	return st, err
}

func (tss *tenableSettingsService) GetTenableSettings() (*dtos.TenableSettingsResponseDTO, error) {
	_, err := tss.getOrCreate()
	if err != nil {
		return nil, err
	}

	return tss.repo.GetSettingsDTO()
}

func (tss *tenableSettingsService) UpdateTenableSettings(input dtos.UpdateTenableSettingsDTO) error {
	st, err := tss.getOrCreate()
	if err != nil {
		return err
	}

	if input.BaseURL != nil {
		st.BaseURL = *input.BaseURL
	}
	if input.AccessKey != nil {
		st.AccessKey = *input.AccessKey
	}
	if input.SecretKey != nil {
		st.SecretKey = *input.SecretKey
	}
	if input.SyncFrequencyHours != nil {
		st.SyncFrequencyHours = *input.SyncFrequencyHours
	}
	if input.LookbackDays != nil {
		st.LookbackDays = *input.LookbackDays
	}

	st.Status = "active"
	st.UpdatedAt = time.Now().UTC()

	return tss.repo.Update(st)
}

func (tss *tenableSettingsService) ActivateSettings(id uint) error {
	return tss.repo.Activate(id)
}

func (tss *tenableSettingsService) DeactivateSettings(id uint) error {
	st, err := tss.repo.GetByID(id)
	if err != nil {
		return err
	}

	st.Status = "inactive"
	st.AccessKey = ""
	st.SecretKey = ""
	st.UpdatedAt = time.Now().UTC()

	return tss.repo.Update(st)
}

func (tss *tenableSettingsService) TestConnection(input dtos.TestTenableConnectionDTO) error {

	if input.BaseURL == "" {
		st, err := tss.getOrCreate()
		if err != nil {
			return fmt.Errorf("failed to load settings: %w", err)
		}
		input.BaseURL = st.BaseURL
	}

	client := &http.Client{Timeout: 20 * time.Second}
	url := fmt.Sprintf("%s/rest/system", input.BaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("x-apikey",
		fmt.Sprintf("accessKey=%s; secretKey=%s",
			input.AccessKey, input.SecretKey,
		),
	)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("tenable returned %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
