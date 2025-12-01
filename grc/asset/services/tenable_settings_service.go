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
	UpdateConnectionSettings(input dtos.UpdateTenableConnectionDTO) error
	UpdateSyncSettings(input dtos.UpdateTenableSyncSettingsDTO) error
	ActivateSettings() error
	DeactivateSettings() error
	TestConnection() error
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

func (s *tenableSettingsService) UpdateConnectionSettings(input dtos.UpdateTenableConnectionDTO) error {
	st, err := s.repo.Get()
	if err != nil {
		return err
	}

	st.BaseURL = input.BaseURL
	st.AccessKey = input.AccessKey
	st.SecretKey = input.SecretKey
	st.UpdatedAt = time.Now().UTC()
	st.Status = "active"

	return s.repo.Update(st)
}

func (s *tenableSettingsService) UpdateSyncSettings(input dtos.UpdateTenableSyncSettingsDTO) error {
	st, err := s.repo.Get()
	if err != nil {
		return err
	}

	st.SyncFrequencyHours = input.SyncFrequencyHours
	st.LookbackDays = input.LookbackDays
	st.Status = "active"

	st.UpdatedAt = time.Now().UTC()

	return s.repo.Update(st)
}

func (s *tenableSettingsService) ActivateSettings() error {
	return s.repo.Activate()
}

func (s *tenableSettingsService) DeactivateSettings() error {
	return s.repo.Deactivate()
}

func (s *tenableSettingsService) TestConnection() error {
	st, err := s.repo.Get()
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 20 * time.Second}
	url := fmt.Sprintf("%s/rest/system", st.BaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("x-apikey", fmt.Sprintf("accessKey=%s; secretKey=%s", st.AccessKey, st.SecretKey))
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
