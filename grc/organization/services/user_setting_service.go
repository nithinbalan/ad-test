package services

import (
	"errors"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	userDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	"jethur/utility"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserSettingServiceInterface interface {
	SetUserService(svc UserServiceInterface)

	GetUserSetting(userID uint) (*userDTOs.UserSettingDTO, error)
	CreateUserSetting(userID uint, userSettingFormDTO *userDTOs.UserSettingFormDTO, tx ...utility.DatabaseTransaction) error
	UpdateUserSetting(userID uint, currentUserID uint, userSettingFormDTO *userDTOs.UserSettingFormDTO) error
	EnableMainNotification(userID uint) error
	DisableMainNotification(userID uint) error

	GetUserSecurityByUserID(userID uint) (*models.UserSecurity, error)
	CreateUserSecurity(userID uint, userSecurityFormDTO *userDTOs.UserSecurityFormDTO, tx ...utility.DatabaseTransaction) error
	UpdateUserSecurity(userSecurity *models.UserSecurity) error
	EnableMFA(userID uint) error
	DisableMFA(userID uint) error
	ResetPassword(userID uint) error

	GetNotificationSettings(userID uint) ([]userDTOs.NotificationSettingListResponseDTO, error)
	UpdateNotificationSetting(userID uint, notificationMasterID uint, enable bool) error
}

type UserSettingService struct {
	repo           organizationRepositories.UserSettingRepositoryInterface
	userSvc        UserServiceInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewUserSettingService(
	repo organizationRepositories.UserSettingRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface) UserSettingServiceInterface {
	return &UserSettingService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
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

func (uss *UserSettingService) CreateUserSetting(userID uint, userSettingFormDTO *userDTOs.UserSettingFormDTO, tx ...utility.DatabaseTransaction) error {
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

	return nil
}

func (uss *UserSettingService) UpdateUserSetting(userID uint, currentUserID uint, userSettingFormDTO *userDTOs.UserSettingFormDTO) error {
	existingUserSetting, err := uss.repo.GetUserSettingsByUserID(userID)

	if err != nil {
		return utility.NewResourceNotFoundError(err, nil)
	}
	// @review
	//do we need to validate timezone
	existingUserSetting.Language = userSettingFormDTO.Language
	existingUserSetting.TimezoneID = userSettingFormDTO.TimezoneID
	existingUserSetting.Timeout = userSettingFormDTO.Timeout
	existingUserSetting.UpdatedBy = &currentUserID

	if err = uss.repo.UpdateUserSetting(existingUserSetting); err != nil {
		return err
	}

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

func (uss *UserSettingService) EnableMFA(userID uint) error {
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

	return nil
}

func (uss *UserSettingService) DisableMFA(userID uint) error {
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

func (uss *UserSettingService) GetNotificationSettings(userID uint) ([]userDTOs.NotificationSettingListResponseDTO, error) {
	modules, err := uss.repo.GetNotificationSetting(userID)
	if err != nil {
		return nil, err
	}

	var response []userDTOs.NotificationSettingListResponseDTO
	for _, module := range modules {
		var notificationDTOs []userDTOs.NotificationListDTO
		// Check if NotificationMasters is nil or not before ranging over it
		if module.NotificationMasters != nil {
			for _, notification := range *module.NotificationMasters {
				enableStatus := false
				// Check if UserNotificationSetting is preloaded and has a value
				if notification.UserNotificationSettings != nil { // Assuming it might be a slice, handle accordingly
					for _, setting := range *notification.UserNotificationSettings {
						if setting.UserID == userID {
							enableStatus = setting.IsEnable
							break // Assuming only one setting per user per notification is possible
						}
					}
				}

				notificationDTO := userDTOs.NotificationListDTO{
					ID:         notification.ID,
					ModulesID:  notification.ModuleID,
					Identifier: notification.Identifier,
					TitleEn:    notification.TitleEn,
					TitleAr:    notification.TitleAr,
					IsEnable:   enableStatus,
				}
				notificationDTOs = append(notificationDTOs, notificationDTO)
			}
		}

		moduleDTO := userDTOs.NotificationSettingListResponseDTO{
			ModuleID:      module.ID,
			ModuleTitleEn: module.TitleEn,
			ModuleTitleAr: module.TitleAr,
			Notifications: &notificationDTOs,
		}
		response = append(response, moduleDTO)
	}

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
