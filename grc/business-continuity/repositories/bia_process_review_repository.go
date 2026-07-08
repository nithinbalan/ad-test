package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BiaProcessReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateBiaProcessReview(review *models.BiaProcessWorkflowReview, tx ...utility.DatabaseTransaction) error
	UpdateBiaProcess(biaProcess *models.BiaProcess, tx ...utility.DatabaseTransaction) error
}

type BiaProcessReviewRepository struct {
	db *gorm.DB
}

func NewBiaProcessReviewRepository(db *gorm.DB) BiaProcessReviewRepositoryInterface {
	return &BiaProcessReviewRepository{db: db}
}

func (repo *BiaProcessReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *BiaProcessReviewRepository) CreateBiaProcessReview(review *models.BiaProcessWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(review).Error
}

func (repo *BiaProcessReviewRepository) UpdateBiaProcess(biaProcess *models.BiaProcess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(biaProcess).Error
}
