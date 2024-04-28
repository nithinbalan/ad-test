package repositories

import (
	masterModels "jethur/grc/masters/models"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type UserSettingRepositoryInterface interface {
	GetNotificationSetting(userID uint) ([]masterModels.Module, error)
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

func (repo *UserSettingRepository) GetNotificationSetting(userID uint) ([]masterModels.Module, error) {
	var modules []masterModels.Module
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
