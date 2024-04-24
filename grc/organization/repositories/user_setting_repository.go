package repositories

import (
	masterModels "jethur/grc/masters/models"
	organizationModels "jethur/grc/organization/models"

	"gorm.io/gorm"
)

type UserSettingRepositoryInterface interface {
	GetNotificationSetting(userID uint) ([]masterModels.Module, error)
	FindNotificationSetting(userID uint, notificationMasterID uint) (*organizationModels.UserNotificationSetting, error)
	CreateNotificationSetting(setting *organizationModels.UserNotificationSetting) error
	UpdateNotificationSetting(setting *organizationModels.UserNotificationSetting) error
}

type UserSettingRepository struct {
	db *gorm.DB
}

func NewUserSettingRepository(db *gorm.DB) *UserSettingRepository {
	return &UserSettingRepository{
		db: db,
	}
}

func (repo *UserSettingRepository) GetNotificationSetting(userID uint) ([]masterModels.Module, error) {
	var modules []masterModels.Module
	err := repo.db.Preload("NotificationMasters", func(db *gorm.DB) *gorm.DB {
		return db.Preload("UserNotificationSettings", "user_id = ?", userID)
	}).Find(&modules).Error
	return modules, err
}

func (repo *UserSettingRepository) FindNotificationSetting(userID, notificationMasterID uint) (*organizationModels.UserNotificationSetting, error) {
	var setting organizationModels.UserNotificationSetting
	err := repo.db.Where("user_id = ? AND notification_master_id = ?", userID, notificationMasterID).First(&setting).Error
	return &setting, err
}

func (repo *UserSettingRepository) CreateNotificationSetting(setting *organizationModels.UserNotificationSetting) error {
	return repo.db.Create(setting).Error
}

func (repo *UserSettingRepository) UpdateNotificationSetting(setting *organizationModels.UserNotificationSetting) error {
	return repo.db.Save(setting).Error
}
