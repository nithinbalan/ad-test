package services

import (
	"crypto/rand"
	"fmt"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	userDTOs "jethur/grc/organization/dtos"
	organizationModels "jethur/grc/organization/models"
	organizationRepositories "jethur/grc/organization/repositories"
	settingService "jethur/grc/settings/services"
	"jethur/utility"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	GetAll(params *userDTOs.UserListQueryParamsDTO) (*userDTOs.UserListResponseDTO, error)
	// GetByID(id uint) (*models.ComplianceRegister, error)
	Create(currentUserID uint, userFormDTO *userDTOs.UserFormDTO) (uint, error)
	Update(id uint, currentUserID uint, user *userDTOs.UserFormDTO) error
	// Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type UserService struct {
	repo           organizationRepositories.UserRepositoryInterface
	designationSvc DesignationServiceInterface
	userStatusSvc  UserStatusServiceInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	roleSvc        settingService.RoleServiceInterface
}

func NewUserService(
	repo organizationRepositories.UserRepositoryInterface,
	designationSvc DesignationServiceInterface,
	userStatusSvc UserStatusServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	roleSvc settingService.RoleServiceInterface) UserServiceInterface {
	return &UserService{
		repo:           repo,
		designationSvc: designationSvc,
		userStatusSvc:  userStatusSvc,
		activityLogSvc: activityLogSvc,
		roleSvc:        roleSvc,
	}
}

func (us *UserService) GetAll(params *userDTOs.UserListQueryParamsDTO) (*userDTOs.UserListResponseDTO, error) {
	return us.repo.GetAll(params)
}

func (us *UserService) Create(currentUserID uint, userFormDTO *userDTOs.UserFormDTO) (uint, error) {

	tx := us.repo.BeginTransaction()

	if err := us.validateUserFormDTO(nil, userFormDTO); err != nil {
		return 0, err
	}

	password, err := us.generateRandomPassword(12)
	if err != nil {
		return 0, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	userStatus, err := us.userStatusSvc.GetByStatus("active")
	if err != nil {
		return 0, err
	}

	user := &organizationModels.User{
		OrganizationUnit1ID: userFormDTO.OrganizationUnit1ID,
		OrganizationUnit2ID: userFormDTO.OrganizationUnit2ID,
		OrganizationUnit3ID: userFormDTO.OrganizationUnit3ID,
		OrganizationUnit4ID: userFormDTO.OrganizationUnit4ID,
		OrganizationUnit5ID: userFormDTO.OrganizationUnit5ID,
		DesignationID:       userFormDTO.DesignationID,
		UserStatusID:        userStatus.ID,
		FirstNameEn:         userFormDTO.FirstNameEn,
		LastNameEn:          userFormDTO.LastNameEn,
		FirstNameAr:         userFormDTO.FirstNameAr,
		LastNameAr:          userFormDTO.LastNameAr,
		Email:               userFormDTO.Email,
		MobileNumber:        userFormDTO.MobileNumber,
		ProfileImagePath:    userFormDTO.ProfileImagePath,
		CreatedBy:           currentUserID,
	}

	userID, err := us.repo.Create(tx, user)
	if err != nil {
		tx.Rollback()
		return 0, err // Only return the error
	}

	userSecurity := organizationModels.UserSecurity{
		UserID:       userID,
		PasswordHash: string(passwordHash),
	}

	if err = us.repo.CreateUserSecurity(tx, userSecurity); err != nil {
		tx.Rollback()
		return 0, err
	}

	userSetting := organizationModels.UserSetting{
		UserID:     userID,
		TimezoneID: 5,
		Language:   "en",
		Timeout:    5,
	}

	if err = us.repo.CreateUserSetting(tx, userSetting); err != nil {
		tx.Rollback()
		return 0, err
	}

	// Prepare a slice of UserRole for bulk insert
	userRoles := make([]organizationModels.UserRole, len(userFormDTO.RoleIDs))
	for i, roleID := range userFormDTO.RoleIDs {
		userRoles[i] = organizationModels.UserRole{
			UserID: userID,
			RoleID: roleID,
		}
	}

	// Now pass the entire array to add roles
	if err = us.repo.CreateUserRole(tx, userRoles); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := us.activityLogSvc.LogActivity(
			logModels.CreateAction,
			organizationModels.User{}.TableName(),
			userID, // use userID instead of id
			userFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": organizationModels.User{}.TableName(),
				"recordID":  userID, // use userID instead of id
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return userID, nil // Use userID, which was modified within the transaction block
}

func (us *UserService) Update(id uint, currentUserID uint, userFormDTO *userDTOs.UserFormDTO) error {

	tx := us.repo.BeginTransaction()
	existingUser, err := us.repo.GetByID(id)

	if err != nil {
		return err
	}

	if err := us.validateUserFormDTO(&id, userFormDTO); err != nil {
		return err
	}

	existingUser.OrganizationUnit1ID = userFormDTO.OrganizationUnit1ID
	existingUser.OrganizationUnit2ID = userFormDTO.OrganizationUnit2ID
	existingUser.OrganizationUnit3ID = userFormDTO.OrganizationUnit3ID
	existingUser.OrganizationUnit4ID = userFormDTO.OrganizationUnit4ID
	existingUser.OrganizationUnit5ID = userFormDTO.OrganizationUnit5ID
	existingUser.DesignationID = userFormDTO.DesignationID
	existingUser.FirstNameEn = userFormDTO.FirstNameEn
	existingUser.LastNameEn = userFormDTO.LastNameEn
	existingUser.FirstNameAr = userFormDTO.FirstNameAr
	existingUser.LastNameAr = userFormDTO.LastNameAr
	existingUser.Email = userFormDTO.Email
	existingUser.MobileNumber = userFormDTO.MobileNumber
	existingUser.ProfileImagePath = userFormDTO.ProfileImagePath
	existingUser.UpdatedBy = currentUserID

	err = us.repo.Update(tx, existingUser)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = us.repo.UpdateUserRoles(tx, id, userFormDTO.RoleIDs); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (us *UserService) validateUserFormDTO(id *uint, userFormDTO *userDTOs.UserFormDTO) error {
	_, err := us.designationSvc.GetByID(userFormDTO.DesignationID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("designation does not exist"), "not_exist", "designation")
	}

	if len(userFormDTO.RoleIDs) == 0 {
		return utility.NewValidationError(fmt.Errorf("no roles provided, at least one role is required"), "required", "role")
	} else {
		for _, roleID := range userFormDTO.RoleIDs {
			// Attempt to fetch each role by ID
			_, err := us.roleSvc.GetByID(roleID)
			if err != nil {
				// If any role ID is invalid, return an error
				return utility.NewValidationError(fmt.Errorf("role does not exist"), "not_exist", "role")
			}

		}
	}

	if id == nil {
		if exists, err := us.repo.ExistsByEmail(userFormDTO.Email); err != nil {
			return err // Handle error from checking email existence
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("email already exists: %s", userFormDTO.Email), "unique", "email")
		}
	} else {
		if exists, err := us.repo.ExistsByEmailExcludingID(*id, userFormDTO.Email); err != nil {
			return err // Handle error from checking email existence
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("email already exists: %s", userFormDTO.Email), "unique", "email")
		}
	}

	return nil
}

func (us *UserService) generateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random password: %w", err)
	}

	for i, v := range b {
		b[i] = charset[v%byte(len(charset))]
	}

	return string(b), nil
}
