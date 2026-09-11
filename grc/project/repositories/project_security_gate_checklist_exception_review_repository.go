package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectSecurityGateChecklistExceptionReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateReview(review *models.ProjectSecurityGateChecklistExceptionWorkflowReview, tx ...utility.DatabaseTransaction) error
	CreateHistory(history *models.ProjectSecurityGateChecklistExceptionWorkflowHistory, tx ...utility.DatabaseTransaction) error
	CreateReviewDocuments(documents []*models.ProjectSecurityGateChecklistExceptionWorkflowReviewDocument, tx ...utility.DatabaseTransaction) error
	UpdateChecklist(checklist *models.ProjectSecurityGateChecklist, tx ...utility.DatabaseTransaction) error
}

type ProjectSecurityGateChecklistExceptionReviewRepository struct {
	db *gorm.DB
}

func NewProjectSecurityGateChecklistExceptionReviewRepository(db *gorm.DB) ProjectSecurityGateChecklistExceptionReviewRepositoryInterface {
	return &ProjectSecurityGateChecklistExceptionReviewRepository{db: db}
}

func (repo *ProjectSecurityGateChecklistExceptionReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ProjectSecurityGateChecklistExceptionReviewRepository) CreateReview(review *models.ProjectSecurityGateChecklistExceptionWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(review).Error
}

func (repo *ProjectSecurityGateChecklistExceptionReviewRepository) CreateHistory(history *models.ProjectSecurityGateChecklistExceptionWorkflowHistory, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(history).Error
}

func (repo *ProjectSecurityGateChecklistExceptionReviewRepository) CreateReviewDocuments(documents []*models.ProjectSecurityGateChecklistExceptionWorkflowReviewDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(documents).Error
}

func (repo *ProjectSecurityGateChecklistExceptionReviewRepository) UpdateChecklist(checklist *models.ProjectSecurityGateChecklist, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(checklist).Error
}
