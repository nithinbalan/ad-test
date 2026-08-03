package services

import (
	"errors"
	"time"

	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	businessContinuityRepositories "jethur/grc/business-continuity/repositories"
	"jethur/grc/models"
	organizationSvc "jethur/grc/organization/services"
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
	repo    businessContinuityRepositories.BiaProcessRepositoryInterface
	userSvc organizationSvc.UserServiceInterface
}

func NewBiaProcessService(repo businessContinuityRepositories.BiaProcessRepositoryInterface, userSvc organizationSvc.UserServiceInterface) BiaProcessServiceInterface {
	return &BiaProcessService{
		repo:    repo,
		userSvc: userSvc,
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
	isCyberUser := false
	if currentUserID != 0 {
		var err error
		isCyberUser, err = svc.userSvc.IsCyberUser(currentUserID)
		if err != nil {
			return err
		}
	}

	existingBiaProcess, err := svc.repo.GetByProcessID(processID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			now := time.Now()
			newBiaProcess := &models.BiaProcess{
				ProcessID:   processID,
				CreatedBYID: &currentUserID,
				CreatedAt:   now,
				IsCyber:     isCyberUser,
			}
			return svc.repo.Create(newBiaProcess, tx...)
		}
		return err
	}

	// BiaProcess already exists - keep is_cyber in sync with the latest contributor's cyber status
	if existingBiaProcess.IsCyber != isCyberUser {
		existingBiaProcess.IsCyber = isCyberUser
		return svc.repo.Update(existingBiaProcess, tx...)
	}
	return nil
}

func (svc *BiaProcessService) Update(biaProcess *models.BiaProcess, tx ...utility.DatabaseTransaction) error {
	return svc.repo.Update(biaProcess, tx...)
}

func (svc *BiaProcessService) DeleteByProcessID(processID uint, tx ...utility.DatabaseTransaction) error {
	return svc.repo.DeleteByProcessID(processID, tx...)
}
