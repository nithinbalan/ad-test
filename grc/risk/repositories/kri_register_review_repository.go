package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type KRIRegisterReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateKRIRegisterReview(review *models.KRIRegisterWorkflowReview, tx ...utility.DatabaseTransaction) error
	UpdateKRIRegister(kriRegister *models.KRIRegister, tx ...utility.DatabaseTransaction) error
}

type KRIRegisterReviewRepository struct {
	db *gorm.DB
}

func NewKRIRegisterReviewRepository(db *gorm.DB) *KRIRegisterReviewRepository {
	return &KRIRegisterReviewRepository{db: db}
}

func (repo *KRIRegisterReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *KRIRegisterReviewRepository) CreateKRIRegisterReview(review *models.KRIRegisterWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(review).Error; err != nil {
		return err
	}
	return nil
}

func (repo *KRIRegisterReviewRepository) UpdateKRIRegister(kriRegister *models.KRIRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(kriRegister).Error
}
