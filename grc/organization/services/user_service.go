package services

import (
	"crypto/rand"
	"fmt"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	userDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	settingService "jethur/grc/settings/services"
	"jethur/utility"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	GetAll(params *userDTOs.UserListQueryParamsDTO) (*userDTOs.UserListResponseDTO, error)
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(currentUserID uint, userFormDTO *userDTOs.UserFormDTO) (uint, error)
	Update(id uint, currentUserID uint, user *userDTOs.UserFormDTO) error
	UpdatePartial(id uint, dto *userDTOs.UserUpdateDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	//User Security
	GetUserSecurityByUserID(userID uint) (*models.UserSecurity, error)
	UpdateUserSecurity(userSecurity *models.UserSecurity) error
	EnableMFA(userID uint) error
	DisableMFA(userID uint) error
	ResetPassword(userID uint) error

	//user Settings
	GetUserSetting(userID uint) (*userDTOs.UserSettingDTO, error)
	UpdateUserSetting(userID uint, currentUserID uint, userSettingFormDTO *userDTOs.UserSettingFormDTO) error
	EnableNotification(userID uint) error
	DisableNotification(userID uint) error
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
	// @todo - need to validate organization unit in validateUserFormDTO

	password, err := us.generateRandomPassword(12)
	if err != nil {
		return 0, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	userStatus, err := us.userStatusSvc.GetByType("active")
	if err != nil {
		return 0, err
	}

	user := &models.User{
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

	userID, err := us.repo.Create(user, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	userSecurity := models.UserSecurity{
		UserID:       userID,
		PasswordHash: string(passwordHash),
	}

	if err = us.repo.CreateUserSecurity(userSecurity, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	//@todo - timezone to be corrected
	userSetting := models.UserSetting{
		UserID:     userID,
		TimezoneID: 5,
		Language:   "en",
		Timeout:    5,
	}

	if err = us.repo.CreateUserSetting(userSetting, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	userRoles := make([]models.UserRole, len(userFormDTO.RoleIDs))
	for i, roleID := range userFormDTO.RoleIDs {
		userRoles[i] = models.UserRole{
			UserID: userID,
			RoleID: roleID,
		}
	}

	if err = us.repo.CreateUserRole(userRoles, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	// @todo profile image upload

	tx.Commit()

	go func() {
		err := us.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.User{}.TableName(),
			userID, // use userID instead of id
			userFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": models.User{}.TableName(),
				"recordID":  userID, // use userID instead of id
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return userID, nil
}

func (us *UserService) Update(id uint, currentUserID uint, userFormDTO *userDTOs.UserFormDTO) error {

	tx := us.repo.BeginTransaction()

	//@Review
	// Fetch roles as part of user
	// existingUser, err := us.repo.GetByID(id, "Roles")
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
	existingUser.UpdatedBy = &currentUserID

	err = us.repo.Update(existingUser, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	//@review
	//Get user final role ids
	//delete existing user roles
	//create new set of roles.
	if err = us.repo.UpdateUserRoles(id, userFormDTO.RoleIDs, tx); err != nil {
		tx.Rollback()
		return err
	}

	// @todo profile image update

	tx.Commit()

	return nil
}

func (us *UserService) UpdatePartial(id uint, userUpdateDTO *userDTOs.UserUpdateDTO) error {

	updates := make(map[string]interface{})

	v := reflect.ValueOf(userUpdateDTO).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.IsValid() && !field.IsNil() {
			jsonTag := t.Field(i).Tag.Get("json")
			tagKey := strings.Split(jsonTag, ",")[0]
			updates[tagKey] = field.Elem().Interface()
		}
	}
	err := us.repo.UpdatePartial(id, updates)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) GetByID(id uint) (*models.User, error) {
	return us.repo.GetByID(id)
}

func (crs *UserService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := crs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = crs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := crs.activityLogSvc.LogBulkActivities(
				logModels.DeleteAction,
				models.User{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
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
			_, err := us.roleSvc.GetByID(roleID)
			if err != nil {
				return utility.NewValidationError(fmt.Errorf("role does not exist"), "not_exist", "role")
			}

		}
	}

	if id == nil {
		if exists, err := us.repo.ExistsByEmail(userFormDTO.Email); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("email already exists: %s", userFormDTO.Email), "unique", "email")
		}
	} else {
		if exists, err := us.repo.ExistsByEmailExcludingID(*id, userFormDTO.Email); err != nil {
			return err
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

func (us *UserService) GetByEmail(email string) (*models.User, error) {
	return us.repo.GetByEmail(email)
}

func (us *UserService) GetUserSecurityByUserID(userId uint) (*models.UserSecurity, error) {
	return us.repo.GetUserSecurityByUserID(userId)
}

func (us *UserService) UpdateUserSecurity(userSecurity *models.UserSecurity) error {
	return us.repo.UpdateUserSecurity(userSecurity)
}

func (us *UserService) GetUserSetting(userID uint) (*userDTOs.UserSettingDTO, error) {
	var userSettingDTO *userDTOs.UserSettingDTO
	userSetting, err := us.repo.GetUserSettingsByUserID(userID)

	if err != nil {
		return nil, utility.NewResourceNotFoundError(err, userID)
	}

	// Map the result to your DTO
	err = mapstructure.Decode(userSetting, &userSettingDTO)
	if err != nil {
		return nil, err
	}

	return userSettingDTO, nil
}

func (us *UserService) UpdateUserSetting(userID uint, currentUserID uint, userSettingFormDTO *userDTOs.UserSettingFormDTO) error {
	existingUserSetting, err := us.repo.GetUserSettingsByUserID(userID)

	if err != nil {
		return utility.NewResourceNotFoundError(err, nil)
	}
	// @review
	//do we need to validate timezone
	existingUserSetting.Language = userSettingFormDTO.Language
	existingUserSetting.TimezoneID = userSettingFormDTO.TimezoneID
	existingUserSetting.Timeout = userSettingFormDTO.Timeout
	existingUserSetting.UpdatedBy = &currentUserID

	if err = us.repo.UpdateUserSetting(existingUserSetting); err != nil {
		return err
	}

	return nil

}

func (us *UserService) EnableMFA(userID uint) error {
	existingUserSecurity, err := us.repo.GetUserSecurityByUserID(userID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, nil)
	}

	existingUserSecurity.IsMFA = true

	if err = us.repo.UpdateUserSecurity(existingUserSecurity); err != nil {
		return err
	}
	//@todo - send email to user
	//@todo - need capture activity log

	return nil

}

func (us *UserService) DisableMFA(userID uint) error {
	existingUserSecurity, err := us.repo.GetUserSecurityByUserID(userID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, nil)
	}

	existingUserSecurity.IsMFA = false

	if err = us.repo.UpdateUserSecurity(existingUserSecurity); err != nil {
		return err
	}
	//@todo - send email to user
	//@todo - need capture activity log

	return nil
}

func (us *UserService) EnableNotification(userID uint) error {
	existingUserSetting, err := us.repo.GetUserSettingsByUserID(userID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, nil)
	}

	existingUserSetting.IsNotification = true

	if err = us.repo.UpdateUserSetting(existingUserSetting); err != nil {
		return err
	}
	//@todo - send email to user
	//@todo - need capture activity log

	return nil

}

func (us *UserService) DisableNotification(userID uint) error {
	existingUserSetting, err := us.repo.GetUserSettingsByUserID(userID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, nil)
	}

	existingUserSetting.IsNotification = false

	if err = us.repo.UpdateUserSetting(existingUserSetting); err != nil {
		return err
	}
	//@todo - send email to user
	//@todo - need capture activity log

	return nil
}

func (us *UserService) ResetPassword(userID uint) error {

	existingUserSecurity, err := us.repo.GetUserSecurityByUserID(userID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, nil)
	}

	newPassword, err := us.generateRandomPassword(12) // Assume you have this function implemented
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	existingUserSecurity.PasswordHash = string(hashedPassword)

	if err = us.repo.UpdateUserSecurity(existingUserSecurity); err != nil {
		return err
	}

	user, err := us.repo.GetByID(userID)
	if err != nil {
		return err
	}

	emailData := userDTOs.SendNewPasswordEmailDTO{
		FirstNameEn: user.FirstNameEn,
		LastNameEn:  user.LastNameEn,
		Email:       user.Email,
		Password:    newPassword,
	}

	to := []string{user.Email}
	subject := "Important: Your Account Password Has Been Updated"
	templateRelativePath := "organization/templates/email/send_new_password.html"

	err = utility.SendTemplatedEmail(to, subject, templateRelativePath, emailData, nil)

	if err != nil {
		return err
	}
	return nil
}
