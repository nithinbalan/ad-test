package repositories

import (
	"errors"

	"jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type NotificationRepositoryInteface interface {
	BeginTransaction() utility.DatabaseTransaction
	Create(notifications []*models.UserNotification, tx ...utility.DatabaseTransaction) error
	GetNotificationMasterByIdentifier(identifier string) (*models.NotificationMaster, error)
	GetUserNotificationSetting(userID uint, notificationMasterID uint) (*models.UserNotificationSetting, error)
	CreateEmailLog(emailLog *models.EmailLog, tx ...utility.DatabaseTransaction) error
	EnqueueEmail(email *models.Email, recipients []*models.EmailRecipient, attachments []*models.EmailAttachment, tx ...utility.DatabaseTransaction) error
	InsertEmailAttachments(attachments []*models.EmailAttachment, tx ...utility.DatabaseTransaction) error

	GetSocialMediaLinks() (*dtos.SocialMediaLinksDTO, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserSettingByUserID(userID uint) (*models.UserSetting, error)
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

func (repo *NotificationRepository) Create(notifications []*models.UserNotification, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(&notifications).Error
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

// EnqueueEmail inserts the email job row, all recipient rows, and all attachment rows
// in a single transaction so the queue entry is always consistent.
func (repo *NotificationRepository) EnqueueEmail(email *models.Email, recipients []*models.EmailRecipient, attachments []*models.EmailAttachment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Insert parent email row first to get the auto-generated ID
	if err := db.Create(email).Error; err != nil {
		return err
	}

	// Assign generated ID to all recipients then bulk insert
	for _, r := range recipients {
		r.EmailID = email.ID
	}
	if len(recipients) > 0 {
		if err := db.Create(&recipients).Error; err != nil {
			return err
		}
	}

	// Assign generated ID to all attachments then bulk insert
	for _, a := range attachments {
		a.EmailID = email.ID
	}
	if len(attachments) > 0 {
		if err := db.Create(&attachments).Error; err != nil {
			return err
		}
	}

	return nil
}

func (repo *NotificationRepository) InsertEmailAttachments(attachments []*models.EmailAttachment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if len(attachments) == 0 {
		return nil
	}
	return db.Create(&attachments).Error
}

func (repo *NotificationRepository) GetSocialMediaLinks() (*dtos.SocialMediaLinksDTO, error) {
	var links dtos.SocialMediaLinksDTO
	result := repo.db.Model(&models.StrategicProfile{}).
		Select("linkedin", "x", "facebook").
		First(&links)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &links, nil
}

func (repo *NotificationRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := repo.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repo *NotificationRepository) GetUserSettingByUserID(userID uint) (*models.UserSetting, error) {
	var userSetting models.UserSetting
	if err := repo.db.Where("user_id = ?", userID).First(&userSetting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &userSetting, nil
}
