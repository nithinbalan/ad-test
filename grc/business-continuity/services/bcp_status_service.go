package services

import (
	"errors"
	"jethur/grc/business-continuity/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BcpStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.BcpStatus, error)
	GetByID(id uint) (*models.BcpStatus, error)
	GetByType(status string) (*models.BcpStatus, error)
}

// Service and implementations
type BcpStatusService struct {
	repo repositories.BcpStatusRepositoryInterface
}

func NewBcpStatusService(repo repositories.BcpStatusRepositoryInterface) BcpStatusServiceInterface {
	return &BcpStatusService{
		repo: repo,
	}
}

func (bss *BcpStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.BcpStatus, error) {
	return bss.repo.GetAll(sort, sortBy, keyword, status)
}

func (bss BcpStatusService) GetByID(id uint) (*models.BcpStatus, error) {
	bcpStatus, err := bss.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return bcpStatus, nil
}

func (bss *BcpStatusService) GetByType(status string) (*models.BcpStatus, error) {
	bcpStatus, err := bss.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return bcpStatus, nil
}
