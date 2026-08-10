package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IdentityReviewCampaignReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateReview(review *models.IdentityReviewCampaignWorkflowReview, tx ...utility.DatabaseTransaction) error
	CreateHistory(history *models.IdentityReviewCampaignWorkflowHistory, tx ...utility.DatabaseTransaction) error
	UpdateCampaign(campaign *models.IdentityReviewCampaign, tx ...utility.DatabaseTransaction) error
}

type IdentityReviewCampaignReviewRepository struct {
	db *gorm.DB
}

func NewIdentityReviewCampaignReviewRepository(db *gorm.DB) IdentityReviewCampaignReviewRepositoryInterface {
	return &IdentityReviewCampaignReviewRepository{
		db: db,
	}
}

func (repo *IdentityReviewCampaignReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IdentityReviewCampaignReviewRepository) CreateReview(review *models.IdentityReviewCampaignWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(review).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IdentityReviewCampaignReviewRepository) CreateHistory(history *models.IdentityReviewCampaignWorkflowHistory, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(history).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IdentityReviewCampaignReviewRepository) UpdateCampaign(campaign *models.IdentityReviewCampaign, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(campaign).Error
}
