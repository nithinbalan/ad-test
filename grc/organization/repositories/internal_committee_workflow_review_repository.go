package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type InternalCommitteeWorkflowReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateInternalCommitteeWorkflowReview(committeeReview *models.InternalCommitteeWorkflowReview, tx ...utility.DatabaseTransaction) error

	UpdateInternalCommitteeVersion(version *models.InternalCommitteeVersion, tx ...utility.DatabaseTransaction) error
}

type InternalCommitteeWorkflowReviewRepository struct {
	db *gorm.DB
}

func NewInternalCommitteeWorkflowReviewRepository(db *gorm.DB) *InternalCommitteeWorkflowReviewRepository {
	return &InternalCommitteeWorkflowReviewRepository{db: db}
}

// BeginTransaction starts a new transaction.
func (repo *InternalCommitteeWorkflowReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *InternalCommitteeWorkflowReviewRepository) CreateInternalCommitteeWorkflowReview(committeeReview *models.InternalCommitteeWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(committeeReview).Error; err != nil {
		return err
	}
	return nil
}

func (repo *InternalCommitteeWorkflowReviewRepository) UpdateInternalCommitteeVersion(version *models.InternalCommitteeVersion, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(version).Error
}
