package services

import (
	"errors"
	"fmt"

	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	businessFunctionRepositories "jethur/grc/business-continuity/repositories"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"gorm.io/gorm"
)

type BcpTeamsServiceInterface interface {
	GetAll(params *businessContinuityDTOs.BcpTeamsListQueryParamsDTO, bcpId uint) ([]*businessContinuityDTOs.BcpTeamsListDTO, error)
	GetByID(id uint) (*models.BcpTeam, error)
	Create(bcpId uint, currentUserID uint, bcpTeamForm *businessContinuityDTOs.BcpTeamsFormDTO) (uint, error)
	Update(id uint, currentUserID uint, bcpTeamForm *businessContinuityDTOs.BcpTeamsFormDTO) error
	Delete(TeamIDs []uint, currentUserID uint) ([]uint, error)
}

type BcpTeamsService struct {
	repo    businessFunctionRepositories.BcpTeamsRepositoryInterface
	userSvc organizationServices.UserServiceInterface
}

func NewBcpTeamsService(repo businessFunctionRepositories.BcpTeamsRepositoryInterface, userSvc organizationServices.UserServiceInterface) BcpTeamsServiceInterface {
	return &BcpTeamsService{
		repo:    repo,
		userSvc: userSvc,
	}
}

func (brs *BcpTeamsService) GetAll(params *businessContinuityDTOs.BcpTeamsListQueryParamsDTO, bcpId uint) ([]*businessContinuityDTOs.BcpTeamsListDTO, error) {
	return brs.repo.GetAll(params, bcpId)
}

func (brs *BcpTeamsService) GetByID(id uint) (*models.BcpTeam, error) {
	Team, err := brs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return Team, nil
}

func (brs *BcpTeamsService) Create(bcpId uint, currentUserID uint, bcpTeamForm *businessContinuityDTOs.BcpTeamsFormDTO) (uint, error) {
	tx := brs.repo.BeginTransaction()

	organizationUnitID, err := brs.userSvc.GetUserPrimaryUnitID(uint(bcpTeamForm.ResponsibleUserID))
	if err != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "responsible_user_unit")
	}

	bcpTeam := &models.BcpTeam{
		BusinessContinuityPlanID: bcpId,
		ResponsibleUserID:        bcpTeamForm.ResponsibleUserID,
		ResponsibleUnitID:        organizationUnitID,
		Responsibility:           bcpTeamForm.Responsibility,
		BackupUserID:             bcpTeamForm.BackupUserID,
		CreatedByID:              currentUserID,
	}

	id, err := brs.repo.Create(bcpTeam, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return id, nil
}

func (brs *BcpTeamsService) Update(id uint, currentUserID uint, bcpTeamForm *businessContinuityDTOs.BcpTeamsFormDTO) error {
	tx := brs.repo.BeginTransaction()

	existingTeam, err := brs.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if existingTeam.ResponsibleUserID != bcpTeamForm.ResponsibleUserID {
		organizationUnitID, err := brs.userSvc.GetUserPrimaryUnitID(bcpTeamForm.ResponsibleUserID)
		if err != nil {
			tx.Rollback()
			return utility.NewValidationError(
				fmt.Errorf("organization unit not set for responsible user"),
				"not_exist",
				"responsible_user_unit",
			)
		}
		existingTeam.ResponsibleUnitID = organizationUnitID
	}

	existingTeam.Responsibility = bcpTeamForm.Responsibility
	existingTeam.BackupUserID = bcpTeamForm.BackupUserID

	if err := brs.repo.Update(existingTeam, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (brs *BcpTeamsService) Delete(TeamIDs []uint, currentUserID uint) ([]uint, error) {
	tx := brs.repo.BeginTransaction()

	var failedIDs []uint

	for _, TeamID := range TeamIDs {
		if _, err := brs.repo.GetByID(TeamID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, TeamID)
				continue
			}
			tx.Rollback()
			return nil, err
		}

		if err := brs.repo.Delete(TeamID, tx); err != nil {
			failedIDs = append(failedIDs, TeamID)
			continue
		}
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, nil
	}

	tx.Commit()
	return failedIDs, nil
}
