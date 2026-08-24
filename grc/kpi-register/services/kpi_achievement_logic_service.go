package services

import (
	"jethur/grc/kpi-register/repositories"
	"jethur/grc/models"
)

type KPIAchievementLogicServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.KPIAchievementLogic, error)
	GetByID(id uint) (*models.KPIAchievementLogic, error)
	GetByType(logicType string) (*models.KPIAchievementLogic, error)
}

// Service and implementations
type KPIAchievementLogicService struct {
	repo repositories.KPIAchievementLogicRepositoryInterface
}

func NewKPIAchievementLogicService(repo repositories.KPIAchievementLogicRepositoryInterface) KPIAchievementLogicServiceInterface {
	return &KPIAchievementLogicService{
		repo: repo,
	}
}

func (kals *KPIAchievementLogicService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.KPIAchievementLogic, error) {
	return kals.repo.GetAll(sort, sortBy, keyword, status)
}

func (kals *KPIAchievementLogicService) GetByID(id uint) (*models.KPIAchievementLogic, error) {
	return kals.repo.GetByID(id)
}

func (kals *KPIAchievementLogicService) GetByType(logicType string) (*models.KPIAchievementLogic, error) {
	logicResult, err := kals.repo.GetByType(logicType)
	if err != nil {
		return nil, err
	}

	return logicResult, nil
}
