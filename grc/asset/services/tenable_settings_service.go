package services

import (
	"encoding/json"
	"fmt"
	"io"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	"jethur/grc/models"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

type TenableSettingsService interface {
	GetTenableSettings() (*dtos.TenableSettingsResponseDTO, error)
	UpdateTenableSettings(input dtos.UpdateTenableSettingsDTO) error
	DeleteTenableSettings() error
	TestConnection(input dtos.TestTenableConnectionDTO) (*dtos.TestConnectionResult, error)
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

	// update fields
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

	// automatically update status based on credentials
	if st.BaseURL != "" && st.AccessKey != "" && st.SecretKey != "" {
		st.Status = "active"
	} else {
		st.Status = "inactive"
	}

	st.UpdatedAt = time.Now().UTC()
	return tss.repo.Update(st)
}

func (tss *tenableSettingsService) DeleteTenableSettings() error {
	st, err := tss.repo.GetFirst()
	if err != nil {
		return err
	}
	if st == nil {
		return fmt.Errorf("no Tenable settings found")
	}

	st.AccessKey = ""
	st.SecretKey = ""
	st.Status = "inactive"
	st.UpdatedAt = time.Now().UTC()

	return tss.repo.Update(st)
}

func (tss *tenableSettingsService) TestConnection(input dtos.TestTenableConnectionDTO) (*dtos.TestConnectionResult, error) {

	if input.BaseURL == "" {
		st, err := tss.getOrCreate()
		if err != nil {
			return &dtos.TestConnectionResult{
				Success: false,
				Message: "Failed to load settings",
			}, nil
		}
		input.BaseURL = st.BaseURL
	}

	client := &http.Client{Timeout: 20 * time.Second}
	url := fmt.Sprintf("%s/rest/system", input.BaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &dtos.TestConnectionResult{
			Success: false,
			Message: "Unable to create request",
		}, nil
	}

	req.Header.Set("x-apikey",
		fmt.Sprintf("accessKey=%s; secretKey=%s",
			input.AccessKey, input.SecretKey,
		),
	)

	resp, err := client.Do(req)
	if err != nil {
		return &dtos.TestConnectionResult{
			Success: false,
			Message: "Connection failed: " + err.Error(),
		}, nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		var terr dtos.TenableTestErrorResponse
		if err := json.Unmarshal(body, &terr); err == nil && terr.ErrorMsg != "" {
			return &dtos.TestConnectionResult{
				Success: false,
				Message: strings.TrimSpace(terr.ErrorMsg),
			}, nil
		}

		// fallback: return plain body if unmarshal fails
		return &dtos.TestConnectionResult{
			Success: false,
			Message: string(body),
		}, nil
	}

	return &dtos.TestConnectionResult{
		Success: true,
		Message: "Connection successful",
	}, nil
}
