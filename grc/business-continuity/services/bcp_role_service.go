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

type BcpRolesServiceInterface interface {
	GetAll(params *businessContinuityDTOs.BcpRolesListQueryParamsDTO, bcpId uint) ([]*businessContinuityDTOs.BcpRolesListDTO, error)
	GetByID(id uint) (*models.BcpRole, error)
	Create(bcpId uint, currentUserID uint, bcpRoleForm *businessContinuityDTOs.BcpRolesFormDTO) (uint, error)
	Update(id uint, currentUserID uint, bcpRoleForm *businessContinuityDTOs.BcpRolesFormDTO) error
	Delete(roleIDs []uint, currentUserID uint) ([]uint, error)
}

type BcpRolesService struct {
	repo    businessFunctionRepositories.BcpRolesRepositoryInterface
	userSvc organizationServices.UserServiceInterface
}

func NewBcpRolesService(repo businessFunctionRepositories.BcpRolesRepositoryInterface, userSvc organizationServices.UserServiceInterface) BcpRolesServiceInterface {
	return &BcpRolesService{
		repo:    repo,
		userSvc: userSvc,
	}
}

func (brs *BcpRolesService) GetAll(params *businessContinuityDTOs.BcpRolesListQueryParamsDTO, bcpId uint) ([]*businessContinuityDTOs.BcpRolesListDTO, error) {
	return brs.repo.GetAll(params, bcpId)
}

func (brs *BcpRolesService) GetByID(id uint) (*models.BcpRole, error) {
	role, err := brs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return role, nil
}

func (brs *BcpRolesService) Create(bcpId uint, currentUserID uint, bcpRoleForm *businessContinuityDTOs.BcpRolesFormDTO) (uint, error) {
	tx := brs.repo.BeginTransaction()

	organizationUnitID, err := brs.userSvc.GetUserPrimaryUnitID(uint(bcpRoleForm.ResponsibleUserID))
	if err != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "responsible_user_unit")
	}

	bcpRole := &models.BcpRole{
		BusinessContinuityPlanID: bcpId,
		ResponsibleUserID:        bcpRoleForm.ResponsibleUserID,
		ResponsibleUnitID:        organizationUnitID,
		Responsibility:           bcpRoleForm.Responsibility,
		BackupUserID:             bcpRoleForm.BackupUserID,
		CreatedByID:              currentUserID,
	}

	id, err := brs.repo.Create(bcpRole, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return id, nil
}

func (brs *BcpRolesService) Update(id uint, currentUserID uint, bcpRoleForm *businessContinuityDTOs.BcpRolesFormDTO) error {
	tx := brs.repo.BeginTransaction()

	existingRole, err := brs.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if existingRole.ResponsibleUserID != bcpRoleForm.ResponsibleUserID {
		organizationUnitID, err := brs.userSvc.GetUserPrimaryUnitID(bcpRoleForm.ResponsibleUserID)
		if err != nil {
			tx.Rollback()
			return utility.NewValidationError(
				fmt.Errorf("organization unit not set for responsible user"),
				"not_exist",
				"responsible_user_unit",
			)
		}
		existingRole.ResponsibleUnitID = organizationUnitID
	}

	existingRole.Responsibility = bcpRoleForm.Responsibility
	existingRole.BackupUserID = bcpRoleForm.BackupUserID

	if err := brs.repo.Update(existingRole, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (brs *BcpRolesService) Delete(roleIDs []uint, currentUserID uint) ([]uint, error) {
	tx := brs.repo.BeginTransaction()

	var failedIDs []uint

	for _, roleID := range roleIDs {
		if _, err := brs.repo.GetByID(roleID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, roleID)
				continue
			}
			tx.Rollback()
			return nil, err
		}

		if err := brs.repo.Delete(roleID, tx); err != nil {
			failedIDs = append(failedIDs, roleID)
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
