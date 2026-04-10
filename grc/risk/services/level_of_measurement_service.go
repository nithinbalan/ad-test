package services

import (
	"jethur/grc/models"
	"jethur/grc/risk/repositories"
)

type LevelOfMeasurementServiceInterface interface {
	GetAll(sort, sortBy, keyword, status string) ([]*models.LevelOfMeasurement, error)
	GetByID(id uint) (*models.LevelOfMeasurement, error)
}

type LevelOfMeasurementService struct {
	repo repositories.LevelOfMeasurementRepositoryInterface
}

func NewLevelOfMeasurementService(repo repositories.LevelOfMeasurementRepositoryInterface) LevelOfMeasurementServiceInterface {
	return &LevelOfMeasurementService{repo: repo}
}

func (s *LevelOfMeasurementService) GetAll(sort, sortBy, keyword, status string) ([]*models.LevelOfMeasurement, error) {
	return s.repo.GetAll(sort, sortBy, keyword, status)
}

func (s *LevelOfMeasurementService) GetByID(id uint) (*models.LevelOfMeasurement, error) {
	return s.repo.GetByID(id)
}
