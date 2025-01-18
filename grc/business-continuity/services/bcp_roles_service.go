package services

import (
	"errors"

	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	businessFunctionRepositories "jethur/grc/business-continuity/repositories"
	"jethur/grc/models"
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
	repo businessFunctionRepositories.BcpRolesRepositoryInterface
}

func NewBcpRolesService(repo businessFunctionRepositories.BcpRolesRepositoryInterface) BcpRolesServiceInterface {
	return &BcpRolesService{
		repo: repo,
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

	bcpRole := &models.BcpRole{
		BusinessContinuityPlanID: bcpId,
		ResponsibleUserID:        bcpRoleForm.ResponsibleUserID,
		ResponsibleUnitID:        bcpRoleForm.ResponsibleUnitID,
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

	existingRole, err := brs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	existingRole.ResponsibleUserID = bcpRoleForm.ResponsibleUserID
	existingRole.ResponsibleUnitID = bcpRoleForm.ResponsibleUnitID
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
