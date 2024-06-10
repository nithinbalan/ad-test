package repositories

import (
	"jethur/grc/models"
	userDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type UserSettingRepositoryInterface interface {
	GetUserSettingsByUserID(userId uint) (*models.UserSetting, error)
	GetUserSettingsWithMfaByUserID(userId uint) (*userDTOs.UserSettingDTO, error)
	CreateUserSetting(userSetting *models.UserSetting, tx ...utility.DatabaseTransaction) error
	UpdateUserSetting(userSetting *models.UserSetting) error

	GetUserSecurityByUserID(userId uint) (*models.UserSecurity, error)
	CreateUserSecurity(userSecurity *models.UserSecurity, tx ...utility.DatabaseTransaction) error
	UpdateUserSecurity(userSecurity *models.UserSecurity) error

	GetNotificationSetting(userID uint) ([]models.Module, error)
	FindNotificationSetting(userID uint, notificationMasterID uint) (*models.UserNotificationSetting, error)
	CreateNotificationSetting(setting *models.UserNotificationSetting) error
	UpdateNotificationSetting(setting *models.UserNotificationSetting) error
}

type UserSettingRepository struct {
	db *gorm.DB
}

func NewUserSettingRepository(db *gorm.DB) *UserSettingRepository {
	return &UserSettingRepository{
		db: db,
	}
}

func (repo *UserSettingRepository) GetUserSettingsByUserID(userID uint) (*models.UserSetting, error) {
	var userSetting *models.UserSetting
	result := repo.db.Where("user_id = ?", userID).First(&userSetting)
	if result.Error != nil {
		return nil, result.Error
	}
	return userSetting, nil
}

func (repo *UserSettingRepository) GetUserSettingsWithMfaByUserID(userID uint) (*userDTOs.UserSettingDTO, error) {
	var userSettingDTO userDTOs.UserSettingDTO
	query := `SELECT us.user_id, us.timezone_id, us.language, us.timeout, us.is_notification, 
	             uss.is_mfa 
	          FROM user_settings us 
	          JOIN user_securities uss ON us.user_id = uss.user_id 
	          WHERE us.user_id = ?`

	result := repo.db.Raw(query, userID).Scan(&userSettingDTO)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userSettingDTO, nil
}

func (repo *UserSettingRepository) CreateUserSetting(userSetting *models.UserSetting, tx ...utility.DatabaseTransaction) error {
	db := repo.db

	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&userSetting).Error; err != nil {
		return err
	}

	return nil
}

func (repo *UserSettingRepository) UpdateUserSetting(userSetting *models.UserSetting) error {
	if err := repo.db.Save(userSetting).Error; err != nil {
		return err
	}
	return nil
}

func (repo *UserSettingRepository) GetUserSecurityByUserID(userId uint) (*models.UserSecurity, error) {
	var userSecurity *models.UserSecurity
	result := repo.db.Where("user_id = ?", userId).First(&userSecurity)

	if result.Error != nil {
		return nil, result.Error
	}

	return userSecurity, nil

}

func (repo *UserSettingRepository) CreateUserSecurity(userSecurity *models.UserSecurity, tx ...utility.DatabaseTransaction) error {
	db := repo.db

	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&userSecurity).Error; err != nil {
		return err
	}

	return nil
}

func (repo *UserSettingRepository) UpdateUserSecurity(userSecurity *models.UserSecurity) error {

	if err := repo.db.Save(userSecurity).Error; err != nil {
		return err
	}

	return nil
}

func (repo *UserSettingRepository) GetNotificationSetting(userID uint) ([]models.Module, error) {
	var modules []models.Module
	err := repo.db.Preload("NotificationMasters", func(db *gorm.DB) *gorm.DB {
		return db.Preload("UserNotificationSettings", "user_id = ?", userID)
	}).Find(&modules).Error
	return modules, err
}

func (repo *UserSettingRepository) FindNotificationSetting(userID, notificationMasterID uint) (*models.UserNotificationSetting, error) {
	var setting models.UserNotificationSetting

	err := repo.db.Where("user_id = ? AND notification_master_id = ?", userID, notificationMasterID).First(&setting).Error

	return &setting, err
}

func (repo *UserSettingRepository) CreateNotificationSetting(setting *models.UserNotificationSetting) error {
	return repo.db.Create(setting).Error
}

func (repo *UserSettingRepository) UpdateNotificationSetting(setting *models.UserNotificationSetting) error {
	return repo.db.Save(setting).Error
}
