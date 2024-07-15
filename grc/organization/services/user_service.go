package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	userDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	settingDTOs "jethur/grc/settings/dtos"
	settingService "jethur/grc/settings/services"
	"jethur/utility"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceInterface interface {
	GetAll(params *userDTOs.UserListQueryParamsDTO) (*userDTOs.UserListResponseDTO, error)
	GetByID(id uint) (*models.User, error)
	GetFullDetailsByID(id uint) (*userDTOs.UserDetailsDTO, error)
	GetByEmail(email string) (*models.User, error)
	Create(currentUserID uint, userFormDTO *userDTOs.UserFormDTO) (uint, error)
	Update(id uint, currentUserID uint, user *userDTOs.UserFormDTO) error
	UpdatePartial(id uint, dto *userDTOs.UserUpdateDTO, tx utility.DatabaseTransaction) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	GetUserLoginHistory(userID uint, params *userDTOs.UserLoginHistoryListQueryParamsDTO) (*userDTOs.UserLoginHistoryListResponseDTO, error)
	CreateUserLoginHistory(user *models.User, loginStatus string, ipAddress string, location string, device string, browser string) error

	GenerateRandomPassword(length int) (string, error)
	SetUserSettingService(svc UserSettingServiceInterface)
	GetPrimaryOrganizationUnitByUserID(userID uint) (*models.OrganizationUnit, error)

}

type UserService struct {
	repo                organizationRepositories.UserRepositoryInterface
	designationSvc      DesignationServiceInterface
	userStatusSvc       UserStatusServiceInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	roleSvc             settingService.RoleServiceInterface
	userSettingSvc      UserSettingServiceInterface
	systemPreferenceSvc settingService.SystemPreferenceServiceInterface
}

func NewUserService(
	repo organizationRepositories.UserRepositoryInterface,
	designationSvc DesignationServiceInterface,
	userStatusSvc UserStatusServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	roleSvc settingService.RoleServiceInterface,
	systemPreferenceSvc settingService.SystemPreferenceServiceInterface) UserServiceInterface {
	return &UserService{  
		repo:                repo,
		designationSvc:      designationSvc,
		userStatusSvc:       userStatusSvc,
		activityLogSvc:      activityLogSvc,
		roleSvc:             roleSvc,
		systemPreferenceSvc: systemPreferenceSvc,
	}
}

func (us *UserService) SetUserSettingService(svc UserSettingServiceInterface) {
	us.userSettingSvc = svc
}

func (us *UserService) GetAll(params *userDTOs.UserListQueryParamsDTO) (*userDTOs.UserListResponseDTO, error) {
	return us.repo.GetAll(params)
}

func (us *UserService) GetByID(id uint) (*models.User, error) {
	user, err := us.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return user, nil
}

func (us *UserService) GetFullDetailsByID(id uint) (*userDTOs.UserDetailsDTO, error) {
	var userDetailsDTO userDTOs.UserDetailsDTO

	user, err := us.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(user, &userDetailsDTO)
	if err != nil {
		return nil, err
	}

	roles := make([]settingDTOs.RoleDTO, 0, len(*user.UserRoles))
	for _, userRole := range *user.UserRoles {
		roles = append(roles, settingDTOs.RoleDTO{
			ID:      userRole.Role.ID,
			TitleEn: userRole.Role.TitleEn,
			TitleAr: userRole.Role.TitleAr,
			Status:  userRole.Role.Status,
		})
	}
	userDetailsDTO.Roles = &roles

	//todo :- also add extra fields which required in details screen

	return &userDetailsDTO, nil
}

func (us *UserService) GetByEmail(email string) (*models.User, error) {
	user, err := us.repo.GetByEmail(email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, email)
		}
		return nil, err
	}

	return user, nil
}

func (us *UserService) Create(currentUserID uint, userFormDTO *userDTOs.UserFormDTO) (uint, error) {

	tx := us.repo.BeginTransaction()

	if err := us.validateUserFormDTO(nil, userFormDTO); err != nil {
		return 0, err
	}
	// @todo - need to validate organization unit in validateUserFormDTO

	password, err := us.GenerateRandomPassword(12)
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

	timezone, err := us.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return 0, err
	}

	initial := generateInitial(userFormDTO.FirstNameEn, userFormDTO.LastNameEn)

	color := utility.GenerateColorCode("user")

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
		CreatedByID:         currentUserID,
		Initial:             initial,
		Color:               color,
	}

	userID, err := us.repo.Create(user, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	userSecurityDTO := &userDTOs.UserSecurityFormDTO{
		PasswordHash: string(passwordHash),
		IsMFA:        false,
	}

	if err = us.userSettingSvc.CreateUserSecurity(userID, userSecurityDTO, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	userSettingDTO := &userDTOs.UserSettingFormDTO{
		TimezoneID:     timezone.TimezoneID,
		Language:       "en",
		Timeout:        5,
		IsNotification: false,
	}

	if err = us.userSettingSvc.CreateUserSetting(userID, userSettingDTO, tx); err != nil {
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
			models.CreateAction,
			models.User{}.TableName(),
			userID, // use userID instead of id
			userFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.User{}.TableName(),
				"recordID":  userID, // use userID instead of id
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	emailData := userDTOs.SendNewPasswordEmailDTO{
		FirstNameEn: userFormDTO.FirstNameEn,
		LastNameEn:  userFormDTO.LastNameEn,
		Email:       userFormDTO.Email,
		Password:    password,
	}

	to := []string{user.Email}
	subject := "Important: Welcome to Jethur!"
	templateRelativePath := "organization/templates/email/user_welcome.html"

	err = utility.SendTemplatedEmail(to, subject, templateRelativePath, emailData, nil)

	if err != nil {
		return 0, err
	}

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

	initial := existingUser.Initial

	if existingUser.FirstNameEn != userFormDTO.FirstNameEn || existingUser.LastNameEn != userFormDTO.LastNameEn {
		initial = generateInitial(userFormDTO.FirstNameEn, userFormDTO.LastNameEn)
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
	existingUser.Initial = initial
	existingUser.UpdatedByID = &currentUserID

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

func (us *UserService) UpdatePartial(id uint, userUpdateDTO *userDTOs.UserUpdateDTO, tx utility.DatabaseTransaction) error {

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
	err := us.repo.UpdatePartial(id, updates, tx)
	if err != nil {
		return err
	}

	return nil
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
				models.DeleteAction,
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

func (us *UserService) GetUserLoginHistory(userID uint, params *userDTOs.UserLoginHistoryListQueryParamsDTO) (*userDTOs.UserLoginHistoryListResponseDTO, error) {
	_, err := us.GetByID(userID)

	if err != nil {
		return nil, err
	}
	return us.repo.GetUserLoginHistory(userID, params)
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

func (us *UserService) GenerateRandomPassword(length int) (string, error) {
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

func (as *UserService) CreateUserLoginHistory(user *models.User, loginStatus string, ipAddress string, location string, device string, browser string) error {
	tx := as.repo.BeginTransaction()

	login := models.UserLoginHistory{
		UserID:      user.ID,
		LoginAt:     time.Now(),
		LoginStatus: loginStatus,
		IPAddress:   ipAddress,
		Location:    &location,
		DeviceInfo:  &device,
		BrowserInfo: &browser,
	}

	err := as.repo.CreateUserLoginHistory(&login, tx)
	if err != nil {
		return err
	}

	loginAt := userDTOs.UserUpdateDTO{
		LastLoginAt: &login.LoginAt,
	}

	err = as.UpdatePartial(user.ID, &loginAt, tx)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func generateInitial(firstName, lastName string) string {
	firstInitial := string(firstName[0])
	lastInitial := string(lastName[0])

	return strings.ToUpper(firstInitial + lastInitial)
}

func (us *UserService) GetPrimaryOrganizationUnitByUserID(userID uint) (*models.OrganizationUnit, error) {
	_, err := us.GetByID(userID)

	if err != nil {
		return nil, err
	}
	return us.repo.GetPrimaryOrganizationUnitByUserID(userID)
}