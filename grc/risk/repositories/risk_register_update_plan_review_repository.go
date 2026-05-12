package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskRegisterUpdatePlanReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	UpdateRiskRegisterUpdatePlan(riskRegisterUpdatePlan *models.RiskRegisterUpdatePlan, tx ...utility.DatabaseTransaction) error
	CreateRiskRegisterUpdatePlanReview(riskRegisterUpdatePlanReview *models.RiskRegisterUpdatePlanWorkflowReview, tx ...utility.DatabaseTransaction) error
	CreateRiskRegisterUpdatePlanWorkflowHistory(historyReview *models.RiskRegisterUpdatePlanWorkflowHistory, tx ...utility.DatabaseTransaction) error
	UpdateRiskRegisterUpdatePlanWorkflowParticipant(participant *models.RiskRegisterUpdatePlanWorkflowParticipant, tx ...utility.DatabaseTransaction) error
}

type RiskRegisterUpdatePlanReviewRepository struct {
	db *gorm.DB
}

func NewRiskRegisterUpdatePlanReviewRepository(db *gorm.DB) *RiskRegisterUpdatePlanReviewRepository {
	return &RiskRegisterUpdatePlanReviewRepository{db: db}
}

func (repo *RiskRegisterUpdatePlanReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskRegisterUpdatePlanReviewRepository) UpdateRiskRegisterUpdatePlan(riskRegisterUpdatePlan *models.RiskRegisterUpdatePlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(riskRegisterUpdatePlan).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanReviewRepository) CreateRiskRegisterUpdatePlanReview(riskRegisterUpdatePlanReview *models.RiskRegisterUpdatePlanWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(riskRegisterUpdatePlanReview).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanReviewRepository) CreateRiskRegisterUpdatePlanWorkflowHistory(historyReview *models.RiskRegisterUpdatePlanWorkflowHistory, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(historyReview).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanReviewRepository) UpdateRiskRegisterUpdatePlanWorkflowParticipant(participant *models.RiskRegisterUpdatePlanWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(participant).Error; err != nil {
		return err
	}

	return nil
}
