package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentClosureReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	UpdateIncidentClosure(incidentClosure *models.IncidentClosure, tx ...utility.DatabaseTransaction) error
	CreateIncidentClosureReview(incidentClosureReview *models.IncidentClosureWorkflowReview, tx ...utility.DatabaseTransaction) error
	CreateIncidentClosureWorkflowHistory(historyReview *models.IncidentClosureWorkflowHistory, tx ...utility.DatabaseTransaction) error
	UpdateIncidentClosureWorkflowParticipant(participant *models.IncidentClosureWorkflowParticipant, tx ...utility.DatabaseTransaction) error
}

type IncidentClosureReviewRepository struct {
	db *gorm.DB
}

func NewIncidentClosureReviewRepository(db *gorm.DB) *IncidentClosureReviewRepository {
	return &IncidentClosureReviewRepository{db: db}
}

func (repo *IncidentClosureReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IncidentClosureReviewRepository) UpdateIncidentClosure(incidentClosure *models.IncidentClosure, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(incidentClosure).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentClosureReviewRepository) CreateIncidentClosureReview(incidentClosureReview *models.IncidentClosureWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(incidentClosureReview).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentClosureReviewRepository) CreateIncidentClosureWorkflowHistory(historyReview *models.IncidentClosureWorkflowHistory, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(historyReview).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentClosureReviewRepository) UpdateIncidentClosureWorkflowParticipant(participant *models.IncidentClosureWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(participant).Error; err != nil {
		return err
	}

	return nil
}
