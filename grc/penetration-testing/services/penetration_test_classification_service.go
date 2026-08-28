package services

import (
	"jethur/grc/models"
	"jethur/grc/penetration-testing/repositories"
)

type PenetrationTestClassificationServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestClassification, error)
}

type PenetrationTestClassificationService struct {
	repo repositories.PenetrationTestClassificationRepositoryInterface
}

func NewPenetrationTestClassificationService(repo repositories.PenetrationTestClassificationRepositoryInterface) PenetrationTestClassificationServiceInterface {
	return &PenetrationTestClassificationService{repo: repo}
}

func (ptcs *PenetrationTestClassificationService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestClassification, error) {
	return ptcs.repo.GetAll(sort, sortBy, keyword, status)
}
