package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskUpdatePlanRiskAnalysisRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetRiskAnalysisList(planID uint, params *riskDTOs.RiskRegisterUpdatePlanListQueryParamsDTO) (*riskDTOs.RiskAnalysisListResponseDTO, error)

	// Bulk Inherent Evaluation
	DeletePlanInherentEvaluationsByRiskRegisterIDs(riskRegisterIDs []uint, planID uint, tx ...utility.DatabaseTransaction) error
	CreateBulkInherentEvaluations(evaluations []*models.RiskEvaluation, tx ...utility.DatabaseTransaction) error
	CreateBulkRiskImpactAnalyses(analyses []*models.RiskImpactAnalysis, tx ...utility.DatabaseTransaction) error
	CreateBulkRiskImpactAnalysisGuidelines(guidelines []*models.RiskImpactAnalysisGuideline, tx ...utility.DatabaseTransaction) error
}

type RiskUpdatePlanRiskAnalysisRepository struct {
	db *gorm.DB
}

func NewRiskUpdatePlanRiskAnalysisRepository(db *gorm.DB) *RiskUpdatePlanRiskAnalysisRepository {
	return &RiskUpdatePlanRiskAnalysisRepository{
		db: db,
	}
}

func (repo *RiskUpdatePlanRiskAnalysisRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskUpdatePlanRiskAnalysisRepository) GetRiskAnalysisList(planID uint, params *riskDTOs.RiskRegisterUpdatePlanListQueryParamsDTO) (*riskDTOs.RiskAnalysisListResponseDTO, error) {
	var risks []*riskDTOs.RiskAnalysisListDTO
	var totalCount int64

	// Build query to get risks from both sources, preferring plan evaluation inherent details
	unionQuery := `
		SELECT
			rr.id AS id,
			rr.reference_number AS reference_number,
			rr.title AS title,
			rr.description AS description,
			COALESCE(plan_eval.id, live_eval.id) AS risk_evaluation_id,
			COALESCE(plan_eval.risk_matrix_likelihood_id, live_eval.risk_matrix_likelihood_id, rr.inherent_likelihood_id) AS inherent_likelihood_id,
			COALESCE(plan_eval.risk_matrix_impact_id, live_eval.risk_matrix_impact_id, rr.inherent_impact_id) AS inherent_impact_id,
			ill.title_en AS inherent_likelihood_title_en,
			ill.title_ar AS inherent_likelihood_title_ar,
			iil.title_en AS inherent_impact_title_en,
			iil.title_ar AS inherent_impact_title_ar,
			irr.title_en AS inherent_risk_rating_title_en,
			irr.title_ar AS inherent_risk_rating_title_ar,
			irrv.value AS inherent_risk_rating_value,
			irr.color_code AS inherent_risk_rating_color_code
		FROM (
			SELECT rr.id
			FROM risk_register AS rr
			WHERE rr.risk_register_update_plan_id = ?

			UNION

			SELECT rr.id
			FROM risk_update_plan_risks AS rupr
			INNER JOIN risk_register AS rr ON rr.id = rupr.risk_register_id
			INNER JOIN risk_update_plan_timelines AS rupt ON rupt.id = rupr.risk_update_plan_timeline_id
			INNER JOIN risk_update_plan_timeline_types AS ruptt ON ruptt.id = rupt.risk_update_plan_timeline_type_id
			WHERE rupt.risk_update_plan_id = ? AND ruptt.type = 'erm_department_evaluation'
		) AS combined_risks
		INNER JOIN risk_register AS rr ON rr.id = combined_risks.id
		LEFT JOIN risk_evaluations AS plan_eval ON plan_eval.id = (SELECT MAX(id) FROM risk_evaluations WHERE risk_register_id = rr.id AND type = 'inherent' AND risk_register_update_plan_id = ?)
		LEFT JOIN risk_evaluations AS live_eval ON live_eval.risk_register_id = rr.id AND live_eval.type = 'inherent' AND live_eval.risk_register_update_plan_id IS NULL AND live_eval.is_latest = 1
		LEFT JOIN risk_matrix_likelihoods AS ill ON ill.id = COALESCE(plan_eval.risk_matrix_likelihood_id, live_eval.risk_matrix_likelihood_id, rr.inherent_likelihood_id)
		LEFT JOIN risk_matrix_impacts AS iil ON iil.id = COALESCE(plan_eval.risk_matrix_impact_id, live_eval.risk_matrix_impact_id, rr.inherent_impact_id)
		LEFT JOIN risk_matrix_risk_ratings AS irr ON irr.id = COALESCE(plan_eval.risk_matrix_risk_rating_id, live_eval.risk_matrix_risk_rating_id, rr.inherent_risk_rating_id)
		LEFT JOIN risk_matrix_risk_rating_values AS irrv ON irrv.id = COALESCE(plan_eval.risk_matrix_risk_rating_value_id, live_eval.risk_matrix_risk_rating_value_id, rr.inherent_risk_rating_value_id)
	`

	query := repo.db.Table("("+unionQuery+") AS combined_risks", planID, planID, planID).Select("*")

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(combined_risks.title LIKE ? OR
			combined_risks.description LIKE ? OR
			combined_risks.reference_number LIKE ? OR
			combined_risks.inherent_likelihood_title_en LIKE ? OR
			combined_risks.inherent_likelihood_title_ar LIKE ? OR
			combined_risks.inherent_impact_title_en LIKE ? OR
			combined_risks.inherent_impact_title_ar LIKE ? OR
			combined_risks.inherent_risk_rating_title_en LIKE ? OR
			combined_risks.inherent_risk_rating_title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Count total results before pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply sorting and pagination
	query = query.Order(params.SortBy + " " + params.Sort)

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&risks).Error; err != nil {
		return nil, err
	}

	// Fetch impact categories and guidelines for each risk
	for _, risk := range risks {
		var impactCategories []riskDTOs.RiskAnalysisImpactCategoryDTO

		var evaluationID uint
		evaluationQuery := repo.db.Table("risk_impact_analyses AS ria").
			Select("ria.risk_evaluation_id").
			Joins("INNER JOIN risk_evaluations AS re ON re.id = ria.risk_evaluation_id").
			Where("re.risk_register_id = ? AND re.type = ?", risk.ID, "inherent").
			Order(fmt.Sprintf("CASE WHEN re.risk_register_update_plan_id = %d THEN 0 ELSE 1 END, ria.risk_evaluation_id DESC", planID)).
			Limit(1)
		if err := evaluationQuery.Scan(&evaluationID).Error; err != nil {
			return nil, err
		}
		if evaluationID == 0 {
			risk.ImpactCategories = impactCategories
			continue
		}

		query := repo.db.Table("risk_impact_analyses AS ria").
			Select("ric.id, ric.title_en, ric.title_ar").
			Joins("INNER JOIN risk_impact_categories AS ric ON ric.id = ria.risk_impact_category_id").
			Where("ria.risk_evaluation_id = ? AND ric.status = ?", evaluationID, "active")

		err := query.Group("ric.id, ric.title_en, ric.title_ar").Find(&impactCategories).Error

		if err != nil {
			return nil, err
		}

		// Fetch guidelines for each impact category that are associated with this risk's impact analysis
		for i := range impactCategories {
			var guidelines []riskDTOs.RiskAnalysisGuidelineDTO

			guidelineQuery := repo.db.Table("risk_impact_analysis_guidelines AS riag").
				Select("rmig.id, rmig.guideline").
				Joins("INNER JOIN risk_impact_analyses AS ria ON ria.id = riag.risk_impact_analysis_id").
				Joins("INNER JOIN risk_matrix_impact_guidelines AS rmig ON rmig.id = riag.risk_matrix_impact_guideline_id").
				Where("ria.risk_evaluation_id = ? AND ria.risk_impact_category_id = ? AND rmig.status = ?", evaluationID, impactCategories[i].ID, "active")

			err := guidelineQuery.Group("rmig.id, rmig.guideline").Find(&guidelines).Error

			if err != nil {
				return nil, err
			}

			impactCategories[i].Guidelines = guidelines
		}

		risk.ImpactCategories = impactCategories
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &riskDTOs.RiskAnalysisListResponseDTO{
		Items: risks,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

// Bulk Inherent Evaluation Repository Methods

func (repo *RiskUpdatePlanRiskAnalysisRepository) DeletePlanInherentEvaluationsByRiskRegisterIDs(riskRegisterIDs []uint, planID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Get the evaluation IDs for the given risk register IDs that belong to this plan draft
	var evaluationIDs []uint
	if err := db.Model(&models.RiskEvaluation{}).
		Where("risk_register_id IN ? AND type = ? AND risk_register_update_plan_id = ?", riskRegisterIDs, "inherent", planID).
		Pluck("id", &evaluationIDs).Error; err != nil {
		return err
	}

	if len(evaluationIDs) == 0 {
		return nil
	}

	// Delete impact analysis guidelines first
	if err := db.Where("risk_impact_analysis_id IN (SELECT id FROM risk_impact_analyses WHERE risk_evaluation_id IN ?)", evaluationIDs).Delete(&models.RiskImpactAnalysisGuideline{}).Error; err != nil {
		return err
	}

	// Delete impact analyses
	if err := db.Where("risk_evaluation_id IN ?", evaluationIDs).Delete(&models.RiskImpactAnalysis{}).Error; err != nil {
		return err
	}

	// Delete risk appetite alignments if they exist
	if err := db.Where("risk_evaluation_id IN ?", evaluationIDs).Delete(&models.RiskAppetiteAlignment{}).Error; err != nil {
		return err
	}

	// Delete the plan draft evaluations
	return db.Where("risk_register_id IN ? AND type = ? AND risk_register_update_plan_id = ?", riskRegisterIDs, "inherent", planID).Delete(&models.RiskEvaluation{}).Error
}

func (repo *RiskUpdatePlanRiskAnalysisRepository) CreateBulkInherentEvaluations(evaluations []*models.RiskEvaluation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	for _, evaluation := range evaluations {
		if err := db.Create(evaluation).Error; err != nil {
			return err
		}
	}

	return nil
}

func (repo *RiskUpdatePlanRiskAnalysisRepository) CreateBulkRiskImpactAnalyses(analyses []*models.RiskImpactAnalysis, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	for _, analysis := range analyses {
		if err := db.Create(analysis).Error; err != nil {
			return err
		}
	}

	return nil
}

func (repo *RiskUpdatePlanRiskAnalysisRepository) CreateBulkRiskImpactAnalysisGuidelines(guidelines []*models.RiskImpactAnalysisGuideline, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	for _, guideline := range guidelines {
		if err := db.Create(guideline).Error; err != nil {
			return err
		}
	}

	return nil
}
