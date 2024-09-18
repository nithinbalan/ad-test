package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type NotificationRepositoryInteface interface {
	BeginTransaction() utility.DatabaseTransaction
	Create(notification *models.UserNotification, tx ...utility.DatabaseTransaction) error
	GetNotificationMasterByIdentifier(identifier string) (*models.NotificationMaster, error)
	GetUserNotificationSetting(userID uint, notificationMasterID uint) (*models.UserNotificationSetting, error)
	CreateEmailLog(emailLog *models.EmailLog, tx ...utility.DatabaseTransaction) error
}

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (repo *NotificationRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *NotificationRepository) Create(notification *models.UserNotification, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(&notification).Error
}

func (repo *NotificationRepository) GetNotificationMasterByIdentifier(identifier string) (*models.NotificationMaster, error) {
	var notificationMaster models.NotificationMaster
	if err := repo.db.Where("identifier = ?", identifier).First(&notificationMaster).Error; err != nil {
		return nil, err
	}
	return &notificationMaster, nil
}

func (repo *NotificationRepository) GetUserNotificationSetting(userID uint, notificationMasterID uint) (*models.UserNotificationSetting, error) {
	var setting models.UserNotificationSetting
	if err := repo.db.Where("user_id = ? AND notification_master_id = ?", userID, notificationMasterID).First(&setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &setting, nil
}

func (repo *NotificationRepository) CreateEmailLog(emailLog *models.EmailLog, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(&emailLog).Error
}
