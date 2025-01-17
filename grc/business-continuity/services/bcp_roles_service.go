package services

import (
	"errors"
	"strconv"

	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	businessFunctionRepositories "jethur/grc/business-continuity/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BcpRolesServiceInterface interface {
	GetAll(bcpId uint, params *businessContinuityDTOs.BcpRolesListQueryParamsDTO) (*businessContinuityDTOs.BcpRolesListResponseDTO, error)
	GetByID(id uint) (*models.BcpRole, error)
	Create(bcpId uint, currentUserID uint, bcpRoleForm *businessContinuityDTOs.BcpRolesFormDTO) (uint, error)
	Update(id uint, currentUserID uint, bcpRoleForm *businessContinuityDTOs.BcpRolesFormDTO) error
	Delete(roleID uint) error
}

type BcpRolesService struct {
	repo businessFunctionRepositories.BcpRolesRepositoryInterface
}

func NewBcpRolesService(repo businessFunctionRepositories.BcpRolesRepositoryInterface) BcpRolesServiceInterface {
	return &BcpRolesService{
		repo: repo,
	}
}

func (brs *BcpRolesService) GetAll(bcpId uint, params *businessContinuityDTOs.BcpRolesListQueryParamsDTO) (*businessContinuityDTOs.BcpRolesListResponseDTO, error) {
	if params.Filters == nil {
		params.Filters = make(map[string][]string)
	}
	params.Filters["bcp_roles.business_continuity_plan_id"] = []string{strconv.FormatUint(uint64(bcpId), 10)}

	return brs.repo.GetAll(params)
}

// GetByID fetches a BCP role by its ID.
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

// Create adds a new BCP role associated with a specific business continuity plan.
func (brs *BcpRolesService) Create(bcpId uint, currentUserID uint, bcpRoleForm *businessContinuityDTOs.BcpRolesFormDTO) (uint, error) {
	// Begin a new transaction
	tx := brs.repo.BeginTransaction()

	// Map the form DTO to the model
	bcpRole := &models.BcpRole{
		BusinessContinuityPlanID: bcpId,
		ResponsibleUserID:        bcpRoleForm.ResponsibleUserID,
		ResponsibleUnitID:        bcpRoleForm.ResponsibleUnitID,
		Responsibility:           bcpRoleForm.Responsibility,
		BackupUserID:             bcpRoleForm.BackupUserID,
		CreatedByID:              currentUserID,
	}

	// Create the BCP role
	id, err := brs.repo.Create(bcpRole, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return id, nil
}

// Update modifies an existing BCP role.
func (brs *BcpRolesService) Update(id uint, currentUserID uint, bcpRoleForm *businessContinuityDTOs.BcpRolesFormDTO) error {
	tx := brs.repo.BeginTransaction()

	// Fetch the existing role
	existingRole, err := brs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	// Update the fields
	existingRole.ResponsibleUserID = bcpRoleForm.ResponsibleUserID
	existingRole.ResponsibleUnitID = bcpRoleForm.ResponsibleUnitID
	existingRole.Responsibility = bcpRoleForm.Responsibility
	existingRole.BackupUserID = bcpRoleForm.BackupUserID
	// existingRole.UpdatedAt = utility.GetCurrentTimePointer()

	// Save the updated role
	if err := brs.repo.Update(existingRole, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (brs *BcpRolesService) Delete(roleID uint) error {
	tx := brs.repo.BeginTransaction()

	if _, err := brs.repo.GetByID(roleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, roleID)
		}
		return err
	}

	if err := brs.repo.Delete(roleID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
