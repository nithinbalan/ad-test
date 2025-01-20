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

type BcpTeamServiceInterface interface {
	GetAll(params *businessContinuityDTOs.BcpTeamsListQueryParamsDTO, bcpId uint) ([]*businessContinuityDTOs.BcpTeamsListDTO, error)
	GetByID(id uint) (*models.BcpTeam, error)
	Create(bcpId uint, currentUserID uint, bcpTeamForm *businessContinuityDTOs.BcpTeamsFormDTO) (uint, error)
	Update(id uint, currentUserID uint, bcpTeamForm *businessContinuityDTOs.BcpTeamsFormDTO) error
	Delete(TeamIDs []uint, currentUserID uint) ([]uint, error)
}

type BcpTeamService struct {
	repo    businessFunctionRepositories.BcpTeamRepositoryInterface
	userSvc organizationServices.UserServiceInterface
}

func NewBcpTeamService(repo businessFunctionRepositories.BcpTeamRepositoryInterface, userSvc organizationServices.UserServiceInterface) BcpTeamServiceInterface {
	return &BcpTeamService{
		repo:    repo,
		userSvc: userSvc,
	}
}

func (brs *BcpTeamService) GetAll(params *businessContinuityDTOs.BcpTeamsListQueryParamsDTO, bcpId uint) ([]*businessContinuityDTOs.BcpTeamsListDTO, error) {
	return brs.repo.GetAll(params, bcpId)
}

func (brs *BcpTeamService) GetByID(id uint) (*models.BcpTeam, error) {
	Team, err := brs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return Team, nil
}

func (brs *BcpTeamService) Create(bcpId uint, currentUserID uint, bcpTeamForm *businessContinuityDTOs.BcpTeamsFormDTO) (uint, error) {
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
		CreatedBYID:              currentUserID,
	}

	id, err := brs.repo.Create(bcpTeam, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return id, nil
}

func (brs *BcpTeamService) Update(id uint, currentUserID uint, bcpTeamForm *businessContinuityDTOs.BcpTeamsFormDTO) error {
	tx := brs.repo.BeginTransaction()

	// Fetch the existing record
	existingTeam, err := brs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	// Update ResponsibleUserID and recalculate ResponsibleUnitID if necessary
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
		existingTeam.ResponsibleUserID = bcpTeamForm.ResponsibleUserID
		existingTeam.ResponsibleUnitID = organizationUnitID
	}

	// Update responsibility if provided
	if bcpTeamForm.Responsibility != "" {
		existingTeam.Responsibility = bcpTeamForm.Responsibility
	}

	// Update BackupUserID (can be nil)
	existingTeam.BackupUserID = bcpTeamForm.BackupUserID

	// Perform the update
	if err := brs.repo.Update(existingTeam, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (brs *BcpTeamService) Delete(TeamIDs []uint, currentUserID uint) ([]uint, error) {
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
