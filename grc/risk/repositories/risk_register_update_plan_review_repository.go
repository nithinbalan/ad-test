package repositories

import (
	"errors"
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
	UpdateRiskEvaluationIsLatestByPlanIDAndType(planID uint, evaluationType string, tx ...utility.DatabaseTransaction) ([]*models.RiskEvaluation, error)
	GetRiskAnalysisAndControlAssociationNotApplicable(planID uint, tx ...utility.DatabaseTransaction) (bool, bool, error)
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

func (repo *RiskRegisterUpdatePlanReviewRepository) UpdateRiskEvaluationIsLatestByPlanIDAndType(planID uint, evaluationType string, tx ...utility.DatabaseTransaction) ([]*models.RiskEvaluation, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Demote all previously-latest evaluations for the risks covered by this plan
	if err := db.Model(&models.RiskEvaluation{}).
		Where("risk_register_id IN (SELECT risk_register_id FROM risk_evaluations WHERE risk_register_update_plan_id = ? AND type = ?) AND type = ? AND is_latest = ?", planID, evaluationType, evaluationType, true).
		Update("is_latest", false).Error; err != nil {
		return nil, err
	}

	// Promote all plan draft evaluations for this type
	if err := db.Model(&models.RiskEvaluation{}).
		Where("risk_register_update_plan_id = ? AND type = ?", planID, evaluationType).
		Update("is_latest", true).Error; err != nil {
		return nil, err
	}

	// Return the promoted evaluations so the risk register ratings can be updated (order by id so newest is applied last)
	var evaluations []*models.RiskEvaluation
	if err := db.Where("risk_register_update_plan_id = ? AND type = ? AND is_latest = ?", planID, evaluationType, true).Order("id asc").Find(&evaluations).Error; err != nil {
		return nil, err
	}
	return evaluations, nil
}

func (repo *RiskRegisterUpdatePlanReviewRepository) GetRiskAnalysisAndControlAssociationNotApplicable(planID uint, tx ...utility.DatabaseTransaction) (bool, bool, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskAnalysisResult struct {
		IsNotApplicable bool
	}
	riskAnalysisErr := db.Model(&models.RiskUpdatePlanTimeline{}).
		Joins("INNER JOIN risk_update_plan_timeline_types AS ruptt ON ruptt.id = risk_update_plan_timelines.risk_update_plan_timeline_type_id").
		Where("risk_update_plan_timelines.risk_update_plan_id = ? AND ruptt.type = ?", planID, "risk_analysis").
		Select("risk_update_plan_timelines.is_not_applicable").
		Scan(&riskAnalysisResult).Error

	if riskAnalysisErr != nil && !errors.Is(riskAnalysisErr, gorm.ErrRecordNotFound) {
		return false, false, riskAnalysisErr
	}

	var riskControlAssociationResult struct {
		IsNotApplicable bool
	}
	riskControlAssociationErr := db.Model(&models.RiskUpdatePlanTimeline{}).
		Joins("INNER JOIN risk_update_plan_timeline_types AS ruptt ON ruptt.id = risk_update_plan_timelines.risk_update_plan_timeline_type_id").
		Where("risk_update_plan_timelines.risk_update_plan_id = ? AND ruptt.type = ?", planID, "risk_control_association").
		Select("risk_update_plan_timelines.is_not_applicable").
		Scan(&riskControlAssociationResult).Error

	if riskControlAssociationErr != nil && !errors.Is(riskControlAssociationErr, gorm.ErrRecordNotFound) {
		return false, false, riskControlAssociationErr
	}

	return riskAnalysisResult.IsNotApplicable, riskControlAssociationResult.IsNotApplicable, nil
}
