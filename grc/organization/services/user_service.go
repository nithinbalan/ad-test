package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonService "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	userDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	settingDTOs "jethur/grc/settings/dtos"
	settingService "jethur/grc/settings/services"
	"jethur/utility"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceInterface interface {
	GetAll(params *userDTOs.UserListQueryParamsDTO) (*userDTOs.UserListResponseDTO, error)
	GetAllUsers() ([]*userDTOs.UserPrimaryDetailsDTO, error)
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
	GetUserPrimaryUnitID(userID uint) (uint, error)
	GetUserIdsByRole(roleID uint, tx ...utility.DatabaseTransaction) ([]uint, error)

	GetCurrentUserDetails(id uint) (*userDTOs.CurrentUserDetailsDTO, error)
	GetUserNotifications(userID uint, params *userDTOs.UserNotificationListQueryParamsDTO) (*userDTOs.UserNotificationListResponseDTO, error)
	GetUserNewNotificationCount(userID uint) (int64, error)

	MarkNotificationsAsRead(userID uint, formDTO *userDTOs.MarkNotificationAsReadFormDTO) (int64, error)
	MarkNotificationsAsUnread(userID uint, formDTO *userDTOs.MarkNotificationAsReadFormDTO) (int64, error)

	ExportUsers(data []*userDTOs.UserListDTO, currentUserID uint) (*os.File, error)

	GetUserPermissions(userID uint, params *settingDTOs.ModuleListQueryParamsDTO) (*[]settingDTOs.ModuleListDTO, error)
}

type UserService struct {
	repo                     organizationRepositories.UserRepositoryInterface
	designationSvc           DesignationServiceInterface
	userStatusSvc            UserStatusServiceInterface
	activityLogSvc           logServices.ActivityLogServiceInterface
	roleSvc                  settingService.RoleServiceInterface
	userSettingSvc           UserSettingServiceInterface
	systemPreferenceSvc      settingService.SystemPreferenceServiceInterface
	organizationLevelSvc     settingService.OrganizationLevelServiceInterface
	moduleSvc                settingService.ModuleServiceInterface
	emailTemplateSettingsSvc settingService.EmailTemplateSettingsServiceInterface
	notificationSvc          commonService.NotificationServiceInterface
}

func NewUserService(
	repo organizationRepositories.UserRepositoryInterface,
	designationSvc DesignationServiceInterface,
	userStatusSvc UserStatusServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	roleSvc settingService.RoleServiceInterface,
	systemPreferenceSvc settingService.SystemPreferenceServiceInterface,
	organizationLevelSvc settingService.OrganizationLevelServiceInterface,
	moduleSvc settingService.ModuleServiceInterface,
	emailTemplateSettingsSvc settingService.EmailTemplateSettingsServiceInterface,
	notificationSvc commonService.NotificationServiceInterface,
) UserServiceInterface {
	return &UserService{
		repo:                     repo,
		designationSvc:           designationSvc,
		userStatusSvc:            userStatusSvc,
		activityLogSvc:           activityLogSvc,
		roleSvc:                  roleSvc,
		systemPreferenceSvc:      systemPreferenceSvc,
		organizationLevelSvc:     organizationLevelSvc,
		moduleSvc:                moduleSvc,
		emailTemplateSettingsSvc: emailTemplateSettingsSvc,
		notificationSvc:          notificationSvc,
	}
}

func (us *UserService) SetUserSettingService(svc UserSettingServiceInterface) {
	us.userSettingSvc = svc
}

func (us *UserService) GetAll(params *userDTOs.UserListQueryParamsDTO) (*userDTOs.UserListResponseDTO, error) {
	// Check if `organization_unit_ids` exists in the filters
	orgUnitIDs, ok := params.Filters["organization_unit_id"]
	if ok && len(orgUnitIDs) > 0 {
		// Get the primary organization level
		primaryLevel, err := us.organizationLevelSvc.GetPrimaryLevel()
		if err != nil {
			return nil, err
		}

		// Determine the field name based on the primary level
		var primaryOrgUnitField string
		switch primaryLevel.Level {
		case 1:
			primaryOrgUnitField = "organization_unit_1_id"
		case 2:
			primaryOrgUnitField = "organization_unit_2_id"
		case 3:
			primaryOrgUnitField = "organization_unit_3_id"
		case 4:
			primaryOrgUnitField = "organization_unit_4_id"
		case 5:
			primaryOrgUnitField = "organization_unit_5_id"
		default:
			return nil, fmt.Errorf("unsupported primary level")
		}

		// Add the resolved column name and value to the filters
		params.Filters[primaryOrgUnitField] = orgUnitIDs
		delete(params.Filters, "organization_unit_id") // Clean up the generic filter key
	}

	return us.repo.GetAll(params)
}

func (us *UserService) GetAllUsers() ([]*userDTOs.UserPrimaryDetailsDTO, error) {
	primaryLevel, err := us.organizationLevelSvc.GetPrimaryLevel()
	if err != nil {
		return nil, err
	}

	var primaryOrgUnitField string
	switch primaryLevel.Level {
	case 1:
		primaryOrgUnitField = "organization_unit_1_id"
	case 2:
		primaryOrgUnitField = "organization_unit_2_id"
	case 3:
		primaryOrgUnitField = "organization_unit_3_id"
	case 4:
		primaryOrgUnitField = "organization_unit_4_id"
	case 5:
		primaryOrgUnitField = "organization_unit_5_id"
	default:
		return nil, fmt.Errorf("unsupported primary level")
	}

	return us.repo.GetAllUsers(primaryOrgUnitField)
}

// func (us *UserService) GetAllUsers() ([]*models.User, error) {
// 	return us.repo.GetAllUsers()
// }

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
			ID:           userRole.Role.ID,
			TitleEn:      userRole.Role.TitleEn,
			TitleAr:      userRole.Role.TitleAr,
			IsSuperAdmin: userRole.Role.IsSuperAdmin,
			Status:       userRole.Role.Status,
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

	var password string
	if userFormDTO.Password != "" {
		password = userFormDTO.Password
	} else {
		var err error
		password, err = us.GenerateRandomPassword(12)
		if err != nil {
			return 0, err
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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
		UserStatusID:        userFormDTO.UserStatusID,
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

	if err = us.userSettingSvc.CreateUserSetting(userID, currentUserID, userSettingDTO, tx); err != nil {
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
		existingUser, err := us.GetFullDetailsByID(userID)
		if err != nil {
			utility.LogError("Failed to fetch user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = us.activityLogSvc.LogActivity(
			models.CreateAction,
			models.User{}.TableName(),
			userID, // use userID instead of id
			existingUser,
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

	if userFormDTO.IsSentEmail {

		emailTemplateSettings, err := us.emailTemplateSettingsSvc.GetEmailTemplateSettings()
		if err != nil {
			return 0, err
		}

		socialLinks, err := us.notificationSvc.GetSocialMediaLinks()
		if err != nil {
			return 0, err
		}

		emailTemplateAndSocialDTO := &commonDTOs.EmailTemplateAndSocialDTO{
			HeaderImageURL: emailTemplateSettings.HeaderImageURL,
			SupportText:    emailTemplateSettings.SupportText,
			ClosingText:    emailTemplateSettings.ClosingText,
			FooterLogoURL:  emailTemplateSettings.FooterLogoURL,
			FooterAddress:  emailTemplateSettings.FooterAddress,
			FooterEmail:    emailTemplateSettings.FooterEmail,
			Linkedin:       &socialLinks.Linkedin,
			X:              &socialLinks.X,
			Facebook:       &socialLinks.Facebook,
		}

		appURL := utility.GetAppURL()
		redirectURL := fmt.Sprintf("%sauth/login", appURL)

		emailData := userDTOs.SendNewPasswordEmailDTO{
			FirstNameEn: userFormDTO.FirstNameEn,
			LastNameEn:  userFormDTO.LastNameEn,
			Email:       userFormDTO.Email,
			Password:    password,
			RedirectURL: redirectURL,
		}

		to := []string{user.Email}
		subject := "Important: Welcome to Jethur!"
		templateRelativePath := "organization/templates/email/user_welcome.html"

		err = utility.SendTemplatedEmail(to, subject, templateRelativePath, emailData, emailTemplateAndSocialDTO, nil)

		if err != nil {
			return 0, err
		}
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

	oldUser, err := us.GetFullDetailsByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := us.validateUserFormDTO(&id, userFormDTO); err != nil {
		return err
	}

	initial := existingUser.Initial

	if existingUser.FirstNameEn != userFormDTO.FirstNameEn || existingUser.LastNameEn != userFormDTO.LastNameEn {
		initial = generateInitial(userFormDTO.FirstNameEn, userFormDTO.LastNameEn)
	}

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
	existingUser.UserStatusID = userFormDTO.UserStatusID

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

	go func() {
		existingUser, err := us.GetFullDetailsByID(id)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		//log changes
		oldData := map[string]interface{}{
			"designation":        oldUser.Designation.TitleEn,
			"first_name_en":      oldUser.FirstNameEn,
			"last_name_en":       oldUser.LastNameEn,
			"first_name_ar":      oldUser.FirstNameAr,
			"last_name_ar":       oldUser.LastNameAr,
			"email":              oldUser.Email,
			"mobile_number":      oldUser.MobileNumber,
			"profile_image_path": oldUser.ProfileImagePath,
			"initial":            oldUser.Initial,
			"user_status":        oldUser.UserStatus.TitleEn,
		}

		newData := map[string]interface{}{
			"designation":        existingUser.Designation.TitleEn,
			"first_name_en":      existingUser.FirstNameEn,
			"last_name_en":       existingUser.LastNameEn,
			"first_name_ar":      existingUser.FirstNameAr,
			"last_name_ar":       existingUser.LastNameAr,
			"email":              existingUser.Email,
			"mobile_number":      existingUser.MobileNumber,
			"profile_image_path": existingUser.ProfileImagePath,
			"initial":            existingUser.Initial,
			"user_status":        existingUser.UserStatus.TitleEn,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := us.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.User{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.User{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()

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
	var deletedUser []map[string]interface{}
	for _, id := range ids {

		user, err := crs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = crs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedUser = append(deletedUser, map[string]interface{}{
			"id":            user.ID,
			"first_name_en": user.FirstNameEn,
			"last_name_en":  user.LastNameEn,
		})
	}

	if len(successfulIds) > 0 {

		go func() {
			err := crs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.User{}.TableName(),
				successfulIds,
				deletedUser,
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

	if userFormDTO.Password != "" {
		if !utility.IsValidPassword(userFormDTO.Password) {
			return utility.NewValidationError(fmt.Errorf("password does not meet the requirements"), "password_complexity", "password")
		}
	}

	// **Validate Max Allowed Users from License Key**

	activeStatus, err := us.userStatusSvc.GetByType("active")
	if err != nil {
		return err
	}

	licenseKey, err := us.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return err
	}

	if *licenseKey.LicenseKey != "" {
		maxAllowedUsers, err := us.DecodeMaxUsersFromLicense(*licenseKey.LicenseKey)
		if err != nil {
			return err
		}

		currentActiveUsers, err := us.repo.CountActiveUsers(activeStatus.ID)
		if err != nil {
			return err
		}

		// If adding a new user (id == nil) or activating an existing one, check against max allowed users
		if id == nil || (id != nil && userFormDTO.UserStatusID == activeStatus.ID) {
			if currentActiveUsers >= int64(maxAllowedUsers) {
				return utility.NewValidationError(
					fmt.Errorf("maximum allowed users limit reached (%d users)", int64(maxAllowedUsers)),
					"user_exceed",
					"user",
				)
			}
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

func (us *UserService) GetUserPrimaryUnitID(userID uint) (uint, error) {
	_, err := us.GetByID(userID)
	if err != nil {
		return 0, err
	}

	primaryLevel, err := us.organizationLevelSvc.GetPrimaryLevel()
	if err != nil {
		return 0, err
	}

	var primaryOrgUnitField string

	switch primaryLevel.Level {
	case 1:
		primaryOrgUnitField = "organization_unit_1_id"
	case 2:
		primaryOrgUnitField = "organization_unit_2_id"
	case 3:
		primaryOrgUnitField = "organization_unit_3_id"
	case 4:
		primaryOrgUnitField = "organization_unit_4_id"
	case 5:
		primaryOrgUnitField = "organization_unit_5_id"
	default:
		return 0, fmt.Errorf("unsupported primary level")
	}

	return us.repo.GetUserPrimaryUnitID(primaryOrgUnitField, userID)
}

func (us *UserService) GetUserIdsByRole(roleID uint, tx ...utility.DatabaseTransaction) ([]uint, error) {
	userIds, err := us.repo.GetUserIdsByRole(roleID, tx...)
	if err != nil {
		return nil, err
	}

	return userIds, nil
}

func (us *UserService) GetCurrentUserDetails(id uint) (*userDTOs.CurrentUserDetailsDTO, error) {
	var userDetailsDTO userDTOs.CurrentUserDetailsDTO

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

	primaryLevel, err := us.organizationLevelSvc.GetPrimaryLevel()
	if err != nil {
		return nil, err
	}

	var primaryOrgUnitID *uint

	switch primaryLevel.Level {
	case 1:
		primaryOrgUnitID = &user.OrganizationUnit1ID
	case 2:
		primaryOrgUnitID = user.OrganizationUnit2ID
	case 3:
		primaryOrgUnitID = user.OrganizationUnit3ID
	case 4:
		primaryOrgUnitID = user.OrganizationUnit4ID
	case 5:
		primaryOrgUnitID = user.OrganizationUnit5ID
	default:
		return nil, fmt.Errorf("invalid primary level: %d", primaryLevel.Level)
	}

	if primaryOrgUnitID != nil {

		userDetailsDTO.PrimaryUnitID = primaryOrgUnitID
	}

	roles := make([]settingDTOs.RoleDTO, 0, len(*user.UserRoles))
	for _, userRole := range *user.UserRoles {
		roles = append(roles, settingDTOs.RoleDTO{
			ID:           userRole.Role.ID,
			TitleEn:      userRole.Role.TitleEn,
			TitleAr:      userRole.Role.TitleAr,
			Status:       userRole.Role.Status,
			IsSuperAdmin: userRole.Role.IsSuperAdmin,
		})
	}
	userDetailsDTO.Roles = &roles

	//todo :- also add extra fields which required in details screen

	return &userDetailsDTO, nil
}

func (ns *UserService) GetUserNotifications(userID uint, params *userDTOs.UserNotificationListQueryParamsDTO) (*userDTOs.UserNotificationListResponseDTO, error) {
	return ns.repo.GetUserNotifications(userID, params)
}

func (us *UserService) GetUserNewNotificationCount(userID uint) (int64, error) {
	return us.repo.GetUserNewNotificationCount(userID)
}

func (us *UserService) MarkNotificationsAsRead(userID uint, formDTO *userDTOs.MarkNotificationAsReadFormDTO) (int64, error) {
	if formDTO.Type == "all" {
		// Mark all notifications as read
		return us.repo.MarkAllNotificationsAsRead(userID)
	} else if formDTO.Type == "individual" && len(formDTO.UserNotificationIDs) > 0 {
		// Mark specific notifications as read
		return us.repo.MarkSpecificNotificationsAsRead(userID, formDTO.UserNotificationIDs)
	} else {
		return 0, fmt.Errorf("invalid parameters")
	}
}

func (us *UserService) MarkNotificationsAsUnread(userID uint, formDTO *userDTOs.MarkNotificationAsReadFormDTO) (int64, error) {
	if formDTO.Type == "all" {
		// Mark all notifications as read
		return us.repo.MarkAllNotificationsAsUnread(userID)
	} else if formDTO.Type == "individual" && len(formDTO.UserNotificationIDs) > 0 {
		// Mark specific notifications as read
		return us.repo.MarkSpecificNotificationsAsUnread(userID, formDTO.UserNotificationIDs)
	} else {
		return 0, fmt.Errorf("invalid parameters")
	}
}

func (us *UserService) ExportUsers(data []*userDTOs.UserListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := us.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	// Define headers based on user language setting
	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"User",
			"Email",
			"Mobile",
			"Designation",
			"Roles",
			"Status",
		}
	} else {
		// Default to English
		headers = []string{
			"SL No.",
			"User",
			"Email",
			"Mobile",
			"Designation",
			"Roles",
			"Status",
		}
	}

	// Convert data to [][]interface{} based on language setting
	rows := make([][]interface{}, len(data))
	for i, user := range data {
		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				user.FirstNameAr + " " + user.LastNameAr,
				user.Email,
				user.MobileNumber,
				user.DesignationTitleAr,
				user.RolesAr,
				user.UserStatusTitleAr,
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				user.FirstNameEn + " " + user.LastNameEn,
				user.Email,
				user.MobileNumber,
				user.DesignationTitleEn,
				user.RolesEn,
				user.UserStatusTitleEn,
			}
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}

func (us *UserService) GetUserPermissions(userID uint, params *settingDTOs.ModuleListQueryParamsDTO) (*[]settingDTOs.ModuleListDTO, error) {
	userRoles, err := us.repo.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	modules, err := us.moduleSvc.GetAll(params)
	if err != nil {
		return nil, err
	}

	userPermissions := map[uint]map[string]bool{}

	for _, role := range userRoles {
		for _, module := range *modules {
			for _, submodule := range module.SubModules {
				if _, exists := userPermissions[submodule.ID]; !exists {
					userPermissions[submodule.ID] = map[string]bool{}
				}

				subGrantedPermissions, err := us.repo.GetPermissionsForModule(role.RoleID, submodule.ID)
				if err != nil {
					return nil, err
				}

				for permissionType := range subGrantedPermissions {
					userPermissions[submodule.ID][permissionType] = true
				}
			}
		}
	}

	var moduleList []settingDTOs.ModuleListDTO
	for _, module := range *modules {
		moduleDTO := settingDTOs.ModuleListDTO{
			ID:         module.ID,
			Identifier: module.Identifier,
			TitleEn:    module.TitleEn,
			TitleAr:    module.TitleAr,
			Order:      module.Order,
			Icon:       module.Icon,
			Route:      module.Route,
			IsEnable:   module.IsEnable,
		}

		for _, submodule := range module.SubModules {
			if len(userPermissions[submodule.ID]) > 0 {
				submoduleDTO := settingDTOs.SubModuleListDTO{
					ID:          submodule.ID,
					Identifier:  submodule.Identifier,
					TitleEn:     submodule.TitleEn,
					TitleAr:     submodule.TitleAr,
					Order:       submodule.Order,
					Icon:        submodule.Icon,
					Route:       submodule.Route,
					IsEnable:    submodule.IsEnable,
					Permissions: userPermissions[submodule.ID],
				}
				moduleDTO.SubModules = append(moduleDTO.SubModules, submoduleDTO)
			}
		}

		if len(moduleDTO.SubModules) > 0 {
			moduleList = append(moduleList, moduleDTO)
		}
	}

	return &moduleList, nil
}

func (us *UserService) DecodeMaxUsersFromLicense(licenseKey string) (float64, error) {
	// Parse JWT token
	token, err := jwt.Parse(licenseKey, func(token *jwt.Token) (interface{}, error) {
		return []byte("f0a3e3b61e7082367820944efeb9d81c7e242ec4d22dc6fa941aaf76d333f158"), nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid license key")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}
	// Extract max_users
	maxUsers, ok := claims["max_users"].(float64) // JWT numbers are float64
	if !ok {
		return 0, fmt.Errorf("max_users field missing in license key")
	}

	return maxUsers, nil
}
