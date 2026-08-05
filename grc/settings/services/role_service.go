package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"

	"jethur/utility"
	"log"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RoleServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, licenseTypeID string) ([]*dtos.RoleDTO, error)
	GetFullDetailsByID(id uint) (*dtos.RoleDetailDTO, error)
	GetByID(id uint) (*models.Role, error)
	Create(currentUserID uint, role *models.Role) (*models.Role, error)
	Update(id uint, currentUserID uint, role *models.Role) (*models.Role, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
	GetByType(roleType string) (*models.Role, error)

	UpdateRolePermissions(roleID uint, currentUserID uint, rolePermissions []*dtos.RolePermissionFormDTO) error
	GetAllWithPermissions(params *dtos.ModuleListQueryParamsDTO, roleID uint) (*[]dtos.ModuleListDTO, error)

	GetSuperAdminRole() (*models.Role, error)
}

type RoleService struct {
	repo                  repositories.RoleRepositoryInterface
	activityLogSvc        services.ActivityLogServiceInterface
	moduleSvc             ModuleServiceInterface
	rolePermissionTypeSvc RolePermissionTypeServiceInterface
	licenseTypeSvc        masterServices.LicenseTypeServiceInterface
}

func NewRoleService(repo repositories.RoleRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	moduleSvc ModuleServiceInterface,
	rolePermissionTypeSvc RolePermissionTypeServiceInterface, licenseTypeSvc masterServices.LicenseTypeServiceInterface) RoleServiceInterface {
	return &RoleService{
		repo:                  repo,
		activityLogSvc:        activityLogSvc,
		moduleSvc:             moduleSvc,
		rolePermissionTypeSvc: rolePermissionTypeSvc,
		licenseTypeSvc:        licenseTypeSvc,
	}
}

func (rs *RoleService) GetAll(sort string, sortBy string, keyword string, status string, licenseTypeID string) ([]*dtos.RoleDTO, error) {
	var roleDTOs []*dtos.RoleDTO
	roles, err := rs.repo.GetAll(sort, sortBy, keyword, status, licenseTypeID)
	if err != nil {
		return nil, err
	}

	err = mapstructure.Decode(roles, &roleDTOs)
	if err != nil {
		return nil, err
	}
	return roleDTOs, nil
}

func (rs RoleService) GetFullDetailsByID(roleId uint) (*dtos.RoleDetailDTO, error) {
	var roleDTO dtos.RoleDetailDTO

	role, err := rs.repo.GetFullDetailsByID(roleId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, roleId)
		}
		return nil, err
	}
	err = mapstructure.Decode(role, &roleDTO)
	if err != nil {
		return nil, err
	}

	return &roleDTO, nil
}

func (rs RoleService) GetByID(id uint) (*models.Role, error) {
	role, err := rs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return role, nil
}
func (rs RoleService) GetByType(roleType string) (*models.Role, error) {
	role, err := rs.repo.GetByType(roleType)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, roleType)
		}
		return nil, err
	}

	return role, nil
}
func (rs *RoleService) Create(currentUserID uint, role *models.Role) (*models.Role, error) {

	if err := rs.validateRoleForm(role, nil); err != nil {
		return nil, err
	}

	defaultType := "other"
	role.Type = &defaultType

	createdRole, err := rs.repo.Create(role)
	if err != nil {
		return nil, err
	}

	go func() {
		err := rs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.Role{}.TableName(),
			role.ID,
			createdRole,
			currentUserID,
			"", "", "", nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdRole, nil
}

func (rs *RoleService) Update(id uint, currentUserID uint, updatedRole *models.Role) (*models.Role, error) {
	existingRole, err := rs.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if err := rs.validateRoleForm(updatedRole, &id); err != nil {
		return nil, err
	}

	if updatedRole.Type == nil {
		updatedRole.Type = existingRole.Type
	}

	result, err := rs.repo.Update(id, updatedRole)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status", "IsCyber"}
	changes, err := utility.CompareFields(existingRole, updatedRole, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := rs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.Role{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}()
	return result, err
}

func (rs *RoleService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		role, err := rs.repo.GetByID(id)
		if err != nil {
			continue
		}

		// Restrict deletion if the role is a super admin
		if role.IsSuperAdmin {
			failedIds = append(failedIds, id)
			continue
		}

		err = rs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := rs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.Role{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (rs *RoleService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := rs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.Role{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}

func (rs *RoleService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := rs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.Role{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				// You might want to log this error or handle it appropriately
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}

func (rs *RoleService) validateRoleForm(role *models.Role, id *uint) error {

	if id == nil {

		if exists, err := rs.repo.ExistsByTitle(role.TitleEn, role.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	} else {

		if exists, err := rs.repo.ExistsByTitleExcludingID(*id, role.TitleEn, role.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	return nil
}

func (rs *RoleService) UpdateRolePermissions(roleID uint, currentUserID uint, rolePermissions []*dtos.RolePermissionFormDTO) error {
	tx := rs.repo.BeginTransaction()

	role, err := rs.GetByID(roleID)
	if err != nil {
		tx.Rollback()
		return err
	}

	passiveLicenseType, err := rs.licenseTypeSvc.GetByType("passive")
	if err != nil {
		tx.Rollback()
		return err
	}

	permissionType, err := rs.rolePermissionTypeSvc.GetByType("passive_enabled")
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, rolePermissionDTO := range rolePermissions {
		if err := rs.validateRolePermissionFormDTO(rolePermissionDTO); err != nil {
			return err
		}
	}

	err = rs.repo.DeleteRolePermissionsByRoleID(roleID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	var changedData []dtos.RolePermissionFormDTO
	for _, rolePermissionDTO := range rolePermissions {
		if *role.LicenseTypeID == passiveLicenseType.ID {
			if rolePermissionDTO.RolePermissionTypeID != permissionType.ID {
				tx.Rollback()
				return utility.NewForbiddenError("Only passive enabled permissions are allowed for passive role")
			}
		}
		rolePermission := &models.RolePermission{
			RoleID:               roleID,
			ModuleID:             rolePermissionDTO.ModuleID,
			RolePermissionTypeID: rolePermissionDTO.RolePermissionTypeID,
			IsGranted:            rolePermissionDTO.IsGranted,
			CreatedBy:            currentUserID,
		}

		_, err := rs.repo.CreateRolePermission(rolePermission, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		changedData = append(changedData, *rolePermissionDTO)
	}

	tx.Commit()

	go func() {
		err := rs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RolePermission{}.TableName(),
			roleID,
			changedData,
			currentUserID,
			"", "", "", nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.RolePermission{}.TableName(),
				"recordID":  roleID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (rs *RoleService) validateRolePermissionFormDTO(rolePermissionFormDTO *dtos.RolePermissionFormDTO) error {
	if _, err := rs.moduleSvc.GetByID(uint(rolePermissionFormDTO.ModuleID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("module does not exist"), "not_exist", "module")
	}

	if _, err := rs.rolePermissionTypeSvc.GetByID(uint(rolePermissionFormDTO.RolePermissionTypeID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("role permission type does not exist"), "not_exist", "role permission type")
	}

	return nil
}

func (rs *RoleService) GetAllWithPermissions(params *dtos.ModuleListQueryParamsDTO, roleID uint) (*[]dtos.ModuleListDTO, error) {
	role, err := rs.GetByID(roleID)
	if err != nil {
		return nil, err
	}

	passiveLicenseType, err := rs.licenseTypeSvc.GetByType("passive")
	if err != nil {
		return nil, err
	}

	isPassiveRole := role.LicenseTypeID != nil && *role.LicenseTypeID == passiveLicenseType.ID

	modules, err := rs.moduleSvc.GetAll(params)
	if err != nil {
		return nil, err
	}

	permissionTypes, err := rs.repo.GetAllPermissionTypes()
	if err != nil {
		return nil, err
	}

	moduleList := make([]dtos.ModuleListDTO, 0)

	for _, module := range *modules {
		var submoduleDTOs []dtos.SubModuleListDTO
		if module.SubModules != nil {
			for _, submodule := range module.SubModules {
				if isPassiveRole && !submodule.IsPassive {
					continue
				}
				subPermissions := initializePermissionsMap(permissionTypes)

				grantedSubPermissions, err := rs.repo.GetPermissionsForModule(roleID, submodule.ID)
				if err != nil {
					return nil, err
				}

				for permissionType := range grantedSubPermissions {
					subPermissions[permissionType] = true
				}

				submoduleDTO := dtos.SubModuleListDTO{
					ID:          submodule.ID,
					Identifier:  submodule.Identifier,
					TitleEn:     submodule.TitleEn,
					TitleAr:     submodule.TitleAr,
					Order:       submodule.Order,
					IsEnable:    submodule.IsEnable,
					Permissions: subPermissions,
					IsPassive:   submodule.IsPassive,
				}
				submoduleDTOs = append(submoduleDTOs, submoduleDTO)
			}
		}

		moduleDTO := dtos.ModuleListDTO{
			ID:         module.ID,
			Identifier: module.Identifier,
			TitleEn:    module.TitleEn,
			TitleAr:    module.TitleAr,
			Order:      module.Order,
			IsEnable:   module.IsEnable,
			SubModules: submoduleDTOs,
			IsPassive:  module.IsPassive,
		}
		moduleList = append(moduleList, moduleDTO)
	}

	return &moduleList, nil
}

func (rs *RoleService) GetSuperAdminRole() (*models.Role, error) {
	role, err := rs.repo.GetSuperAdminRole()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve super admin role: %w", err)
	}
	return role, nil
}

func initializePermissionsMap(permissionTypes []string) map[string]bool {
	permissions := make(map[string]bool)
	for _, permissionType := range permissionTypes {
		permissions[permissionType] = false
	}
	return permissions
}
