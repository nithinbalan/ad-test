package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectDecisionLogReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateReview(review *models.ProjectDecisionLogWorkflowReview, tx ...utility.DatabaseTransaction) error
	CreateHistory(history *models.ProjectDecisionLogWorkflowHistory, tx ...utility.DatabaseTransaction) error
	UpdateDecisionLog(decisionLog *models.ProjectDecisionLog, tx ...utility.DatabaseTransaction) error
}

type ProjectDecisionLogReviewRepository struct {
	db *gorm.DB
}

func NewProjectDecisionLogReviewRepository(db *gorm.DB) *ProjectDecisionLogReviewRepository {
	return &ProjectDecisionLogReviewRepository{db: db}
}

func (pdlrr *ProjectDecisionLogReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(pdlrr.db)
}

func (pdlrr *ProjectDecisionLogReviewRepository) CreateReview(review *models.ProjectDecisionLogWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := pdlrr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(review).Error
}

func (pdlrr *ProjectDecisionLogReviewRepository) CreateHistory(history *models.ProjectDecisionLogWorkflowHistory, tx ...utility.DatabaseTransaction) error {
	db := pdlrr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(history).Error
}

func (pdlrr *ProjectDecisionLogReviewRepository) UpdateDecisionLog(decisionLog *models.ProjectDecisionLog, tx ...utility.DatabaseTransaction) error {
	db := pdlrr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(decisionLog).Error
}
