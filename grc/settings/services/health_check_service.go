package services

import (
	"errors"

	"jethur/grc/settings/repositories"
)

type HealthCheckService interface {
	CheckSystemPreference() error
}

type healthCheckService struct {
	repo              repositories.HealthCheckRepository
	systemPrefService SystemPreferenceServiceInterface
}

func NewHealthCheckService(repo repositories.HealthCheckRepository, systemPrefSvc SystemPreferenceServiceInterface) HealthCheckService {
	return &healthCheckService{
		repo:              repo,
		systemPrefService: systemPrefSvc,
	}
}

func (s *healthCheckService) CheckSystemPreference() error {
	// ensure system preference exists
	pref, err := s.repo.GetSystemPreference()
	if err != nil {
		return err
	}

	if pref == nil {
		return errors.New("system preference not found")
	}

	return nil
}
