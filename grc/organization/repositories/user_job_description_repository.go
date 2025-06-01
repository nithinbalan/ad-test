package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type UserJobDescriptionRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateUserJobDescriptions(tx utility.DatabaseTransaction, userJDs *[]models.UserJobDescription) error
	CreateUserJobDescriptionCompetencies(tx utility.DatabaseTransaction, userJDCs *[]models.UserJobDescriptionCompetency) error
	CreateUserJobDescriptionCertificates(tx utility.DatabaseTransaction, userJDCerts *[]models.UserJobDescriptionCertificate) error

	GetByUserAndJD(userID, jdID uint) (*models.UserJobDescription, error)
	DeleteUserJobDescriptionCertificatesByUserJDID(tx utility.DatabaseTransaction, userJDID uint) error
	DeleteUserJobDescriptionCompetenciesByUserJDID(tx utility.DatabaseTransaction, userJDID uint) error
	DeleteUserJobDescription(tx utility.DatabaseTransaction, userJDID uint) error
}

type UserJobDescriptionRepository struct {
	db *gorm.DB
}

func NewUserJobDescriptionRepository(db *gorm.DB) *UserJobDescriptionRepository {
	return &UserJobDescriptionRepository{
		db: db,
	}
}

func (repo *UserJobDescriptionRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *UserJobDescriptionRepository) CreateUserJobDescriptions(tx utility.DatabaseTransaction, userJDs *[]models.UserJobDescription) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	return db.Create(userJDs).Error
}

func (repo *UserJobDescriptionRepository) CreateUserJobDescriptionCompetencies(tx utility.DatabaseTransaction, userJDCs *[]models.UserJobDescriptionCompetency) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	return db.Create(userJDCs).Error
}

func (repo *UserJobDescriptionRepository) CreateUserJobDescriptionCertificates(tx utility.DatabaseTransaction, userJDCerts *[]models.UserJobDescriptionCertificate) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	return db.Create(userJDCerts).Error
}
func (repo *UserJobDescriptionRepository) GetByUserAndJD(userID, jdID uint) (*models.UserJobDescription, error) {
	var userJD models.UserJobDescription
	err := repo.db.Where("user_id = ? AND job_description_id = ?", userID, jdID).First(&userJD).Error
	if err != nil {
		return nil, err
	}
	return &userJD, nil
}

func (repo *UserJobDescriptionRepository) DeleteUserJobDescriptionCertificatesByUserJDID(tx utility.DatabaseTransaction, userJDID uint) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	return db.Where("user_job_description_competency_id IN (?)",
		db.Model(&models.UserJobDescriptionCompetency{}).Select("id").Where("user_job_description_id = ?", userJDID),
	).Delete(&models.UserJobDescriptionCertificate{}).Error
}

func (repo *UserJobDescriptionRepository) DeleteUserJobDescriptionCompetenciesByUserJDID(tx utility.DatabaseTransaction, userJDID uint) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	return db.Where("user_job_description_id = ?", userJDID).
		Delete(&models.UserJobDescriptionCompetency{}).Error
}

func (repo *UserJobDescriptionRepository) DeleteUserJobDescription(tx utility.DatabaseTransaction, userJDID uint) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	return db.Delete(&models.UserJobDescription{}, userJDID).Error
}
