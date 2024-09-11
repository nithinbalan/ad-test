package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/settings/repositories"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type RoleServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Role, error)
	GetByID(id uint) (*models.Role, error)
	Create(currentUserID uint, role *models.Role) (*models.Role, error)
	Update(id uint, currentUserID uint, role *models.Role) (*models.Role, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type RoleService struct {
	repo           repositories.RoleRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewRoleService(repo repositories.RoleRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) RoleServiceInterface {
	return &RoleService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (rs *RoleService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Role, error) {
	return rs.repo.GetAll(sort, sortBy, keyword, status)
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

func (rs *RoleService) Create(currentUserID uint, role *models.Role) (*models.Role, error) {

	if err := rs.validateRoleForm(role, nil); err != nil {
		return nil, err
	}

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
			nil,
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

	result, err := rs.repo.Update(id, updatedRole)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
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
			nil,
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
