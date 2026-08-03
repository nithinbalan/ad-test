package services

import (
	"errors"
	"fmt"
	"jethur/grc/identity-access-review/dtos"
	"jethur/grc/identity-access-review/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	settingsServices "jethur/grc/settings/services"
	"jethur/utility"

	"gorm.io/gorm"
)

type IdentityReviewUserServiceInterface interface {
	GetAll(params *dtos.IdentityReviewUserListQueryParamsDTO) (*dtos.IdentityReviewUserListResponseDTO, error)
	GetByID(id uint) (*models.IdentityReviewUser, error)
	GetFullDetailsByID(id uint) (*models.IdentityReviewUser, error)
	Create(currentUserID uint, formDTO *dtos.IdentityReviewUserFormDTO) (*models.IdentityReviewUser, error)
	Update(id uint, currentUserID uint, formDTO *dtos.IdentityReviewUserFormDTO) (*models.IdentityReviewUser, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type IdentityReviewUserService struct {
	repo                 repositories.IdentityReviewUserRepositoryInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
	organizationUnitSvc  settingsServices.OrganizationUnitServiceInterface
	privilegedAccountSvc PrivilegedAccountServiceInterface
	accountTypeSvc       AccountTypeServiceInterface
}

func NewIdentityReviewUserService(
	repo repositories.IdentityReviewUserRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	organizationUnitSvc settingsServices.OrganizationUnitServiceInterface,
	privilegedAccountSvc PrivilegedAccountServiceInterface,
	accountTypeSvc AccountTypeServiceInterface,
) IdentityReviewUserServiceInterface {
	return &IdentityReviewUserService{
		repo:                 repo,
		activityLogSvc:       activityLogSvc,
		organizationUnitSvc:  organizationUnitSvc,
		privilegedAccountSvc: privilegedAccountSvc,
		accountTypeSvc:       accountTypeSvc,
	}
}

func (irus *IdentityReviewUserService) GetAll(params *dtos.IdentityReviewUserListQueryParamsDTO) (*dtos.IdentityReviewUserListResponseDTO, error) {
	return irus.repo.GetAll(params)
}

func (irus *IdentityReviewUserService) GetByID(id uint) (*models.IdentityReviewUser, error) {
	identityReviewUser, err := irus.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return identityReviewUser, nil
}

func (irus *IdentityReviewUserService) GetFullDetailsByID(id uint) (*models.IdentityReviewUser, error) {
	identityReviewUser, err := irus.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return identityReviewUser, nil
}

func (irus *IdentityReviewUserService) Create(currentUserID uint, formDTO *dtos.IdentityReviewUserFormDTO) (*models.IdentityReviewUser, error) {
	if err := irus.validateForm(formDTO); err != nil {
		return nil, err
	}

	identityReviewUser := &models.IdentityReviewUser{
		Username:      formDTO.Username,
		DepartmentID:  formDTO.DepartmentID,
		PrivilegeID:   formDTO.PrivilegeID,
		AccountTypeID: formDTO.AccountTypeID,
		CreatedByID:   currentUserID,
	}

	id, err := irus.repo.Create(identityReviewUser)
	if err != nil {
		return nil, err
	}

	createdUser, err := irus.GetFullDetailsByID(id)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := irus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IdentityReviewUser{}.TableName(),
			id,
			createdUser,
			currentUserID,
			"identity_access_review",
			"identity_review_user",
			"user",
			nil,
		); err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

	return createdUser, nil
}

func (irus *IdentityReviewUserService) Update(id uint, currentUserID uint, formDTO *dtos.IdentityReviewUserFormDTO) (*models.IdentityReviewUser, error) {
	existingUser, err := irus.GetByID(id)
	if err != nil {
		return nil, err
	}

	oldUser, err := irus.GetFullDetailsByID(id)
	if err != nil {
		return nil, err
	}

	if err := irus.validateForm(formDTO); err != nil {
		return nil, err
	}

	existingUser.Username = formDTO.Username
	existingUser.DepartmentID = formDTO.DepartmentID
	existingUser.PrivilegeID = formDTO.PrivilegeID
	existingUser.AccountTypeID = formDTO.AccountTypeID
	existingUser.UpdatedByID = &currentUserID

	if err := irus.repo.Update(existingUser); err != nil {
		return nil, err
	}

	updatedUser, err := irus.GetFullDetailsByID(id)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := irus.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IdentityReviewUser{}.TableName(),
			id,
			map[string]interface{}{"old": oldUser, "new": updatedUser},
			currentUserID,
			"identity_access_review",
			"identity_review_user",
			"user",
			nil,
		); err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

	return updatedUser, nil
}

func (irus *IdentityReviewUserService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint
	var deletedUsers []interface{}

	for _, id := range ids {
		identityReviewUser, err := irus.GetFullDetailsByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if err := irus.repo.Delete(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
		deletedUsers = append(deletedUsers, identityReviewUser)
	}

	if len(successfulIDs) > 0 {
		go func() {
			if err := irus.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.IdentityReviewUser{}.TableName(),
				successfulIDs,
				deletedUsers,
				currentUserID,
				nil,
			); err != nil {
				utility.LogError("Failed to log bulk delete activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}

	return failedIDs, nil
}

func (irus *IdentityReviewUserService) validateForm(formDTO *dtos.IdentityReviewUserFormDTO) error {
	if _, err := irus.organizationUnitSvc.GetByID(formDTO.DepartmentID); err != nil {
		return utility.NewValidationError(
			fmt.Errorf("department does not exist"),
			"not_exist",
			"department_id",
		)
	}

	if formDTO.PrivilegeID != nil {
		if _, err := irus.privilegedAccountSvc.GetByID(*formDTO.PrivilegeID); err != nil {
			return utility.NewValidationError(
				fmt.Errorf("privilege does not exist"),
				"not_exist",
				"privilege_id",
			)
		}
	}

	if formDTO.AccountTypeID != nil {
		if _, err := irus.accountTypeSvc.GetByID(*formDTO.AccountTypeID); err != nil {
			return utility.NewValidationError(
				fmt.Errorf("account type does not exist"),
				"not_exist",
				"account_type_id",
			)
		}
	}

	return nil
}
