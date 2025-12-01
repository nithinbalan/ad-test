package services

import (
	"fmt"
	"io"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	"net/http"
	"time"
)

type TenableSettingsService interface {
	GetSettings() (*dtos.TenableSettingsResponseDTO, error)
	UpdateSettings(input dtos.UpdateTenableSettingsDTO) error
	ActivateSettings() error
	DeactivateSettings() error
	TestConnection(input dtos.TestTenableConnectionDTO) error
}

type tenableSettingsService struct {
	repo repositories.TenableSettingsRepository
}

func NewTenableSettingsService(repo repositories.TenableSettingsRepository) TenableSettingsService {
	return &tenableSettingsService{repo: repo}
}

func (s *tenableSettingsService) GetSettings() (*dtos.TenableSettingsResponseDTO, error) {
	st, err := s.repo.Get()
	if err != nil {
		return nil, err
	}

	return &dtos.TenableSettingsResponseDTO{
		ID:                 st.ID,
		BaseURL:            st.BaseURL,
		AccessKey:          st.AccessKey,
		SecretKey:          st.SecretKey,
		SyncFrequencyHours: st.SyncFrequencyHours,
		LookbackDays:       st.LookbackDays,
		Status:             st.Status,
	}, nil
}

func (s *tenableSettingsService) UpdateSettings(input dtos.UpdateTenableSettingsDTO) error {

	st, err := s.repo.Get()
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

	return s.repo.Update(st)
}

func (s *tenableSettingsService) ActivateSettings() error {
	return s.repo.Activate()
}

func (s *tenableSettingsService) DeactivateSettings() error {
	st, err := s.repo.Get()
	if err != nil {
		return err
	}

	st.Status = "inactive"
	st.AccessKey = ""
	st.SecretKey = ""
	st.UpdatedAt = time.Now().UTC()

	return s.repo.Update(st)
}

func (s *tenableSettingsService) TestConnection(input dtos.TestTenableConnectionDTO) error {

	if input.BaseURL == "" {
		st, err := s.repo.Get()
		if err != nil {
			return fmt.Errorf("failed to load settings from DB: %w", err)
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
