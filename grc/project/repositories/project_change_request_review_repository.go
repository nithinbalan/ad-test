package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectChangeRequestReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateReview(review *models.ProjectChangeRequestWorkflowReview, tx ...utility.DatabaseTransaction) error
	CreateHistory(history *models.ProjectChangeRequestWorkflowHistory, tx ...utility.DatabaseTransaction) error
	UpdateChangeRequest(projectChangeRequest *models.ProjectChangeRequest, tx ...utility.DatabaseTransaction) error
}

type ProjectChangeRequestReviewRepository struct {
	db *gorm.DB
}

func NewProjectChangeRequestReviewRepository(db *gorm.DB) *ProjectChangeRequestReviewRepository {
	return &ProjectChangeRequestReviewRepository{db: db}
}

func (pcrrr *ProjectChangeRequestReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(pcrrr.db)
}

func (pcrrr *ProjectChangeRequestReviewRepository) CreateReview(review *models.ProjectChangeRequestWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := pcrrr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(review).Error
}

func (pcrrr *ProjectChangeRequestReviewRepository) CreateHistory(history *models.ProjectChangeRequestWorkflowHistory, tx ...utility.DatabaseTransaction) error {
	db := pcrrr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(history).Error
}

func (pcrrr *ProjectChangeRequestReviewRepository) UpdateChangeRequest(projectChangeRequest *models.ProjectChangeRequest, tx ...utility.DatabaseTransaction) error {
	db := pcrrr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(projectChangeRequest).Error
}
