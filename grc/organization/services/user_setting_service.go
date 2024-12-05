package services

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonService "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	userDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	settingService "jethur/grc/settings/services"
	"jethur/utility"
	"sort"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserSettingServiceInterface interface {
	SetUserService(svc UserServiceInterface)

	GetUserSetting(userID uint) (*userDTOs.UserSettingDTO, error)
	GetUserSettingsByUserID(userID uint) (*models.UserSetting, error)
	CreateUserSetting(userID uint, currentUserID uint, userSettingFormDTO *userDTOs.UserSettingFormDTO, tx ...utility.DatabaseTransaction) error
	UpdateUserSetting(userID uint, currentUserID uint, userSettingFormDTO *userDTOs.UserSettingFormDTO) error
	EnableMainNotification(userID uint) error
	DisableMainNotification(userID uint) error

	GetUserSecurityByUserID(userID uint) (*models.UserSecurity, error)
	CreateUserSecurity(userID uint, userSecurityFormDTO *userDTOs.UserSecurityFormDTO, tx ...utility.DatabaseTransaction) error
	UpdateUserSecurity(userSecurity *models.UserSecurity) error
	EnableMFA(userID uint, currentUserID uint) error
	DisableMFA(userID uint, currentUserID uint) error
	ResetPassword(userID uint) error
	ChangePassword(userID uint, oldPassword, newPassword string) error

	GetNotificationSettings(userID uint) ([]userDTOs.NotificationSettingListResponseDTO, error)
	UpdateNotificationSetting(userID uint, notificationMasterID uint, enable bool) error
}

type UserSettingService struct {
	repo                     organizationRepositories.UserSettingRepositoryInterface
	userSvc                  UserServiceInterface
	activityLogSvc           logServices.ActivityLogServiceInterface
	notificationSvc          commonService.NotificationServiceInterface
	emailTemplateSettingsSvc settingService.EmailTemplateSettingsServiceInterface
}

func NewUserSettingService(
	repo organizationRepositories.UserSettingRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	notificationSvc commonService.NotificationServiceInterface,
	emailTemplateSettingsSvc settingService.EmailTemplateSettingsServiceInterface) UserSettingServiceInterface {
	return &UserSettingService{
		repo:                     repo,
		activityLogSvc:           activityLogSvc,
		notificationSvc:          notificationSvc,
		emailTemplateSettingsSvc: emailTemplateSettingsSvc,
	}
}

func (uss *UserSettingService) SetUserService(svc UserServiceInterface) {
	uss.userSvc = svc
}

func (uss *UserSettingService) GetUserSetting(userID uint) (*userDTOs.UserSettingDTO, error) {
	userSettingDTO, err := uss.repo.GetUserSettingsWithMfaByUserID(userID)

	if err != nil {
		return nil, utility.NewResourceNotFoundError(err, userID)
	}

	return userSettingDTO, nil
}

func (uss *UserSettingService) GetUserSettingsByUserID(userID uint) (*models.UserSetting, error) {
	userSetting, err := uss.repo.GetUserSettingsByUserID(userID)

	if err != nil {
		return nil, utility.NewResourceNotFoundError(err, userID)
	}

	return userSetting, nil
}

func (uss *UserSettingService) CreateUserSetting(userID uint, currentUserID uint, userSettingFormDTO *userDTOs.UserSettingFormDTO, tx ...utility.DatabaseTransaction) error {
	userSettings := &models.UserSetting{
		UserID:         userID,
		TimezoneID:     userSettingFormDTO.TimezoneID,
		Language:       userSettingFormDTO.Language,
		IsNotification: userSettingFormDTO.IsNotification,
	}

	err := uss.repo.CreateUserSetting(userSettings, tx...)
	if err != nil {
		return err
	}
	go func() {
		existingUser, err := uss.GetUserSetting(userID)
		if err != nil {
			utility.LogError("Failed to fetch user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		user, err := uss.userSvc.GetByID(userID)
		if err != nil {
			return
		}
		err = uss.activityLogSvc.LogActivity(
			models.CreateAction,
			models.UserSetting{}.TableName(),
			userID, // use userID instead of id
			existingUser,
			user.CreatedByID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.UserSetting{}.TableName(),
				"recordID":  userID, // use userID instead of id
				"userID":    user.CreatedByID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (uss *UserSettingService) UpdateUserSetting(userID uint, currentUserID uint, userSettingFormDTO *userDTOs.UserSettingFormDTO) error {
	existingUserSetting, err := uss.repo.GetUserSettingsByUserID(userID)

	if err != nil {
		return utility.NewResourceNotFoundError(err, nil)
	}

	oldSetting, err := uss.repo.GetUserSettingsByUserID(userID)
	if err != nil {
		return nil
	}
	// @review
	//do we need to validate timezone
	existingUserSetting.Language = userSettingFormDTO.Language
	existingUserSetting.TimezoneID = userSettingFormDTO.TimezoneID
	existingUserSetting.Timeout = userSettingFormDTO.Timeout
	existingUserSetting.UpdatedByID = &currentUserID

	if err = uss.repo.UpdateUserSetting(existingUserSetting); err != nil {
		return err
	}

	go func() {
		existingSetting, err := uss.repo.GetUserSettingsByUserID(userID)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})
		//log changes
		oldData := map[string]interface{}{
			"language": oldSetting.Language,
			"timezone": oldSetting.Timezone.Name,
			"timeout":  oldSetting.Timeout,
		}

		newData := map[string]interface{}{
			"language": existingSetting.Language,
			"timezone": existingSetting.Timezone.Name,
			"timeout":  existingSetting.Timeout,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := uss.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.UserSetting{}.TableName(),
				existingSetting.ID,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.UserSetting{}.TableName(),
					"recordID":  existingSetting.ID,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()

	return nil
}

func (uss *UserSettingService) EnableMainNotification(userID uint) error {
	existingUserSetting, err := uss.repo.GetUserSettingsByUserID(userID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, nil)
	}

	existingUserSetting.IsNotification = true

	if err = uss.repo.UpdateUserSetting(existingUserSetting); err != nil {
		return err
	}
	//@todo - send email to user
	//@todo - need capture activity log

	return nil
}

func (uss *UserSettingService) DisableMainNotification(userID uint) error {
	existingUserSetting, err := uss.repo.GetUserSettingsByUserID(userID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, nil)
	}

	existingUserSetting.IsNotification = false

	if err = uss.repo.UpdateUserSetting(existingUserSetting); err != nil {
		return err
	}
	//@todo - send email to user
	//@todo - need capture activity log

	return nil
}

func (uss *UserSettingService) GetUserSecurityByUserID(userId uint) (*models.UserSecurity, error) {
	userSecurity, err := uss.repo.GetUserSecurityByUserID(userId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, userId)
		}
		return nil, err
	}

	return userSecurity, nil
}

func (uss *UserSettingService) CreateUserSecurity(userID uint, userSecurityFormDTO *userDTOs.UserSecurityFormDTO, tx ...utility.DatabaseTransaction) error {
	userSecurity := &models.UserSecurity{
		UserID:       userID,
		PasswordHash: userSecurityFormDTO.PasswordHash,
		IsMFA:        userSecurityFormDTO.IsMFA,
	}

	err := uss.repo.CreateUserSecurity(userSecurity, tx...)
	if err != nil {
		return err
	}

	return nil
}

func (uss *UserSettingService) UpdateUserSecurity(userSecurity *models.UserSecurity) error {
	return uss.repo.UpdateUserSecurity(userSecurity)
}

func (uss *UserSettingService) EnableMFA(userID uint, currentUserID uint) error {
	oldUserSecurity, err := uss.repo.GetUserSecurityByUserID(userID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, nil)
	}

	existingUserSecurity, err := uss.repo.GetUserSecurityByUserID(userID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, nil)
	}

	existingUserSecurity.IsMFA = true

	if err = uss.repo.UpdateUserSecurity(existingUserSecurity); err != nil {
		return err
	}
	//@todo - send email to user
	//@todo - need capture activity log
	go func() {
		// Prepare log data
		changedData := make(map[string]interface{})
		//log changes
		oldData := map[string]interface{}{
			"is_mfa": oldUserSecurity.IsMFA,
		}

		newData := map[string]interface{}{
			"is_mfa": existingUserSecurity.IsMFA,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := uss.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.UserSecurity{}.TableName(),
				oldUserSecurity.ID,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.UserSecurity{}.TableName(),
					"recordID":  oldUserSecurity.ID,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()
	return nil
}

func (uss *UserSettingService) DisableMFA(userID uint, currentUserID uint) error {
	oldUserSecurity, err := uss.repo.GetUserSecurityByUserID(userID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, nil)
	}

	existingUserSecurity, err := uss.repo.GetUserSecurityByUserID(userID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, nil)
	}

	existingUserSecurity.IsMFA = false

	if err = uss.repo.UpdateUserSecurity(existingUserSecurity); err != nil {
		return err
	}
	//@todo - send email to user
	//@todo - need capture activity log
	go func() {
		// Prepare log data
		changedData := make(map[string]interface{})
		//log changes
		oldData := map[string]interface{}{
			"is_mfa": oldUserSecurity.IsMFA,
		}

		newData := map[string]interface{}{
			"is_mfa": existingUserSecurity.IsMFA,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := uss.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.UserSecurity{}.TableName(),
				oldUserSecurity.ID,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.UserSecurity{}.TableName(),
					"recordID":  oldUserSecurity.ID,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()
	return nil
}

func (uss *UserSettingService) ResetPassword(userID uint) error {
	existingUserSecurity, err := uss.repo.GetUserSecurityByUserID(userID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, nil)
	}

	newPassword, err := uss.userSvc.GenerateRandomPassword(12) // Assume you have this function implemented
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	existingUserSecurity.PasswordHash = string(hashedPassword)

	if err = uss.repo.UpdateUserSecurity(existingUserSecurity); err != nil {
		return err
	}

	user, err := uss.userSvc.GetByID(userID)
	if err != nil {
		return err
	}

	emailTemplateSettings, err := uss.emailTemplateSettingsSvc.GetEmailTemplateSettings()
	if err != nil {
		return err
	}

	socialLinks, err := uss.notificationSvc.GetSocialMediaLinks()
	if err != nil {
		return err
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

	emailData := userDTOs.SendNewPasswordEmailDTO{
		FirstNameEn: user.FirstNameEn,
		LastNameEn:  user.LastNameEn,
		Email:       user.Email,
		Password:    newPassword,
	}

	to := []string{user.Email}
	subject := "Important: Your Account Password Has Been Updated"
	templateRelativePath := "organization/templates/email/send_new_password.html"

	err = utility.SendTemplatedEmail(to, subject, templateRelativePath, emailData, emailTemplateAndSocialDTO, nil)

	if err != nil {
		return err
	}
	return nil
}

func (uss *UserSettingService) GetNotificationSettings(userID uint) ([]userDTOs.NotificationSettingListResponseDTO, error) {
	modules, err := uss.repo.GetNotificationSetting(userID)
	if err != nil {
		return nil, err
	}

	parentModulesMap := make(map[uint]*userDTOs.NotificationSettingListResponseDTO)

	for _, module := range modules {
		if module.ModuleID == nil {
			// This is a parent module
			parentModulesMap[module.ID] = &userDTOs.NotificationSettingListResponseDTO{
				ModuleID:      module.ID,
				ModuleTitleEn: module.TitleEn,
				ModuleTitleAr: module.TitleAr,
				Order:         module.Order,
				Notifications: &[]userDTOs.NotificationListDTO{},
			}
		}
	}

	// Function to aggregate notifications and ensure uniqueness
	aggregateNotifications := func(module models.Module) map[uint]bool {
		notificationsSet := make(map[uint]bool)
		if module.NotificationMasters != nil {
			for _, notification := range *module.NotificationMasters {
				notificationsSet[notification.ID] = true
			}
		}
		return notificationsSet
	}

	for _, module := range modules {
		notificationsSet := aggregateNotifications(module)
		var notifications []userDTOs.NotificationListDTO

		for notificationID := range notificationsSet {
			var notification *models.NotificationMaster
			for _, notif := range *module.NotificationMasters {
				if notif.ID == notificationID {
					notification = &notif
					break
				}
			}

			enableStatus := false
			if notification != nil && notification.UserNotificationSettings != nil {
				for _, setting := range *notification.UserNotificationSettings {
					if setting.UserID == userID {
						enableStatus = setting.IsEnable
						break
					}
				}
			}

			if notification != nil {
				notificationDTO := userDTOs.NotificationListDTO{
					ID:         notification.ID,
					ModulesID:  notification.ModuleID,
					Identifier: notification.Identifier,
					TitleEn:    notification.TitleEn,
					TitleAr:    notification.TitleAr,
					IsEnable:   enableStatus,
				}

				notifications = append(notifications, notificationDTO)
			}
		}
		// Sort notifications by ID
		sort.Slice(notifications, func(i, j int) bool {
			return notifications[i].ID < notifications[j].ID
		})

		if parentModule, exists := parentModulesMap[module.ID]; exists {
			*parentModule.Notifications = append(*parentModule.Notifications, notifications...)
		} else if module.ModuleID != nil {
			if parentModule, exists := parentModulesMap[*module.ModuleID]; exists {
				*parentModule.Notifications = append(*parentModule.Notifications, notifications...)
			}
		}

	}

	var response []userDTOs.NotificationSettingListResponseDTO
	for _, parentModule := range parentModulesMap {
		if len(*parentModule.Notifications) > 0 {
			response = append(response, *parentModule)
		}
	}

	// Sort the final response by the Order field
	sort.Slice(response, func(i, j int) bool {
		return response[i].Order < response[j].Order
	})

	return response, nil
}

func (uss *UserSettingService) UpdateNotificationSetting(userID uint, notificationMasterID uint, enable bool) error {
	notificationSetting, err := uss.repo.FindNotificationSetting(userID, notificationMasterID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if notificationSetting.ID == 0 {
		// No setting found, create a new one
		createNotificationSetting := &models.UserNotificationSetting{
			UserID:               userID,
			NotificationMasterID: notificationMasterID,
			IsEnable:             enable,
		}
		return uss.repo.CreateNotificationSetting(createNotificationSetting)
	}

	// Existing setting found, update it
	notificationSetting.IsEnable = enable
	return uss.repo.UpdateNotificationSetting(notificationSetting)
}

func (uss *UserSettingService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	currentUserSecurity, err := uss.GetUserSecurityByUserID(userID)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(currentUserSecurity.PasswordHash), []byte(oldPassword))
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("incorrect old password"), "password_mismatch", "old_password")
	}

	if !utility.IsValidPassword(newPassword) {
		return utility.NewValidationError(fmt.Errorf("password does not meet the requirements"), "password_complexity", "password")
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	currentUserSecurity.PasswordHash = string(newPasswordHash)

	err = uss.UpdateUserSecurity(currentUserSecurity)
	if err != nil {
		return err
	}

	user, err := uss.userSvc.GetByID(userID)
	if err != nil {
		return err
	}

	emailTemplateSettings, err := uss.emailTemplateSettingsSvc.GetEmailTemplateSettings()
	if err != nil {
		return err
	}

	socialLinks, err := uss.notificationSvc.GetSocialMediaLinks()
	if err != nil {
		return err
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

	emailData := userDTOs.SendNewPasswordEmailDTO{
		FirstNameEn: user.FirstNameEn,
		LastNameEn:  user.LastNameEn,
		Email:       user.Email,
	}

	to := []string{user.Email}
	subject := "Important: Your Account Password Has Been Updated"
	templateRelativePath := "organization/templates/email/change_password.html"

	err = utility.SendTemplatedEmail(to, subject, templateRelativePath, emailData, emailTemplateAndSocialDTO, nil)

	if err != nil {
		return err
	}

	return nil
}
