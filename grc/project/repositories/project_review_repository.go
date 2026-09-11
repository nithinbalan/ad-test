package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProjectReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateReview(review *models.ProjectWorkflowReview, tx ...utility.DatabaseTransaction) error
	CreateHistory(history *models.ProjectWorkflowHistory, tx ...utility.DatabaseTransaction) error
	UpdateProject(project *models.Project, tx ...utility.DatabaseTransaction) error
}

type ProjectReviewRepository struct {
	db *gorm.DB
}

func NewProjectReviewRepository(db *gorm.DB) *ProjectReviewRepository {
	return &ProjectReviewRepository{db: db}
}

func (prr *ProjectReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(prr.db)
}

func (prr *ProjectReviewRepository) CreateReview(review *models.ProjectWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := prr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(review).Error
}

func (prr *ProjectReviewRepository) CreateHistory(history *models.ProjectWorkflowHistory, tx ...utility.DatabaseTransaction) error {
	db := prr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(history).Error
}

func (prr *ProjectReviewRepository) UpdateProject(project *models.Project, tx ...utility.DatabaseTransaction) error {
	db := prr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Omit(clause.Associations).Save(project).Error
}
