package services

import (
	"jethur/grc/models"
	"jethur/grc/strategy/repositories"
)

type DataSourceServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DataSource, error)
	GetByID(id uint) (*models.DataSource, error)
	GetByType(status string) (*models.DataSource, error) // Renamed method

}

// Service and implementations
type DataSourceService struct {
	repo repositories.DataSourceRepositoryInterface
}

func NewDataSourceService(repo repositories.DataSourceRepositoryInterface) DataSourceServiceInterface {
	return &DataSourceService{
		repo: repo,
	}
}

func (dss *DataSourceService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DataSource, error) {
	return dss.repo.GetAll(sort, sortBy, keyword, status)
}

func (dss DataSourceService) GetByID(id uint) (*models.DataSource, error) {
	return dss.repo.GetByID(id)
}

func (dss *DataSourceService) GetByType(status string) (*models.DataSource, error) {
	// Fetch the kpi status by its status
	statusResult, err := dss.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
