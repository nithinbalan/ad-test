package services

import (
	"errors"
	"time"

	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	businessContinuityRepositories "jethur/grc/business-continuity/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BiaProcessServiceInterface interface {
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.BiaProcess, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BiaProcessDetailDTO, error)
	GetFullDetailsByProcessID(processID uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BiaProcessDetailDTO, error)
	GetByProcessID(processID uint, tx ...utility.DatabaseTransaction) (*models.BiaProcess, error)
	CreateIfNotExists(processID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error
	Update(biaProcess *models.BiaProcess, tx ...utility.DatabaseTransaction) error
	DeleteByProcessID(processID uint, tx ...utility.DatabaseTransaction) error
}

type BiaProcessService struct {
	repo businessContinuityRepositories.BiaProcessRepositoryInterface
}

func NewBiaProcessService(repo businessContinuityRepositories.BiaProcessRepositoryInterface) BiaProcessServiceInterface {
	return &BiaProcessService{
		repo: repo,
	}
}

func (svc *BiaProcessService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.BiaProcess, error) {
	return svc.repo.GetByID(id, tx...)
}

func (svc *BiaProcessService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BiaProcessDetailDTO, error) {
	var dto businessContinuityDTOs.BiaProcessDetailDTO

	biaProcess, err := svc.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(biaProcess, &dto)
	if err != nil {
		return nil, err
	}

	dto.SubmittedAt = biaProcess.SubmittedAt
	dto.ApprovedAt = biaProcess.ApprovedAt
	dto.CreatedAt = biaProcess.CreatedAt
	dto.UpdatedAt = biaProcess.UpdatedAt

	return &dto, nil
}

func (svc *BiaProcessService) GetFullDetailsByProcessID(processID uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BiaProcessDetailDTO, error) {
	var dto businessContinuityDTOs.BiaProcessDetailDTO

	biaProcess, err := svc.repo.GetFullDetailsByProcessID(processID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, processID)
		}
		return nil, err
	}

	err = mapstructure.Decode(biaProcess, &dto)
	if err != nil {
		return nil, err
	}

	dto.SubmittedAt = biaProcess.SubmittedAt
	dto.ApprovedAt = biaProcess.ApprovedAt
	dto.CreatedAt = biaProcess.CreatedAt
	dto.UpdatedAt = biaProcess.UpdatedAt

	return &dto, nil
}

func (svc *BiaProcessService) GetByProcessID(processID uint, tx ...utility.DatabaseTransaction) (*models.BiaProcess, error) {
	return svc.repo.GetByProcessID(processID, tx...)
}

func (svc *BiaProcessService) CreateIfNotExists(processID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	_, err := svc.repo.GetByProcessID(processID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			now := time.Now()
			newBiaProcess := &models.BiaProcess{
				ProcessID:   processID,
				CreatedBYID: &currentUserID,
				CreatedAt:   now,
			}
			return svc.repo.Create(newBiaProcess, tx...)
		}
		return err
	}
	// BiaProcess already exists, do nothing
	return nil
}

func (svc *BiaProcessService) Update(biaProcess *models.BiaProcess, tx ...utility.DatabaseTransaction) error {
	return svc.repo.Update(biaProcess, tx...)
}

func (svc *BiaProcessService) DeleteByProcessID(processID uint, tx ...utility.DatabaseTransaction) error {
	return svc.repo.DeleteByProcessID(processID, tx...)
}
