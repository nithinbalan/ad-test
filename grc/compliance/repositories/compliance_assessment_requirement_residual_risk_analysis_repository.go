package repositories

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceAssessmentRequirementResidualRiskAnalysisRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction

	// Residual Risk Analysis methods
	CreateResidualRiskAnalysis(riskAnalysis []*models.ComplianceAssessmentRequirementResidualRiskAnalysis, tx ...utility.DatabaseTransaction) error
	DeleteResidualRiskAnalysisByResponsibleUnitID(responsibleUnitID uint, tx ...utility.DatabaseTransaction) error
	GetResidualRiskAnalysisByResponsibleUnitID(responsibleUnitID uint) ([]*models.ComplianceAssessmentRequirementResidualRiskAnalysis, error)
	UpdateResidualRiskScore(responsibleUnitID uint, score float64, ratingID *uint, tx ...utility.DatabaseTransaction) error
	GetComplianceAssessmentRatings() ([]*models.ComplianceAssessmentRating, error)

	GetAllByResponsibleUnitID(responsibleUnitID uint, params *complianceDTOs.ComplianceAssessmentResidualControlListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentResidualControlListResponseDTO, error)
	Create(residualControls []*models.ComplianceAssessmentResidualControl, tx ...utility.DatabaseTransaction) error
	DeleteByResponsibleUnitID(responsibleUnitID uint, tx ...utility.DatabaseTransaction) error
	UpdateResidualControlByImpactType(responsibleUnitID uint, impactTypeID uint, tx ...utility.DatabaseTransaction) error

	AssessmentUpdateStatus(complianceAssessment *models.ComplianceAssessment, tx ...utility.DatabaseTransaction) error
}

type ComplianceAssessmentRequirementResidualRiskAnalysisRepository struct {
	db *gorm.DB
}

func NewComplianceAssessmentRequirementResidualRiskAnalysisRepository(db *gorm.DB) *ComplianceAssessmentRequirementResidualRiskAnalysisRepository {
	return &ComplianceAssessmentRequirementResidualRiskAnalysisRepository{
		db: db,
	}
}

func (repo *ComplianceAssessmentRequirementResidualRiskAnalysisRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

// CreateResidualRiskAnalysis creates residual risk analysis records
func (repo *ComplianceAssessmentRequirementResidualRiskAnalysisRepository) CreateResidualRiskAnalysis(riskAnalysis []*models.ComplianceAssessmentRequirementResidualRiskAnalysis, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(riskAnalysis) == 0 {
		return nil
	}

	return db.Create(&riskAnalysis).Error
}

// UpdateResidualRiskScore updates the residual risk score and rating for a responsible unit
func (repo *ComplianceAssessmentRequirementResidualRiskAnalysisRepository) UpdateResidualRiskScore(responsibleUnitID uint, score float64, ratingID *uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	updates := map[string]interface{}{
		"residual_risk_score": score,
	}

	if ratingID != nil {
		updates["residual_risk_rating_id"] = *ratingID
	}

	return db.Model(&models.ComplianceAssessmentRequirementResponsibleUnit{}).
		Where("id = ?", responsibleUnitID).
		Updates(updates).Error
}

// DeleteResidualRiskAnalysisByResponsibleUnitID deletes all residual risk analysis for a responsible unit
func (repo *ComplianceAssessmentRequirementResidualRiskAnalysisRepository) DeleteResidualRiskAnalysisByResponsibleUnitID(responsibleUnitID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("compliance_assessment_requirement_responsible_unit_id = ?", responsibleUnitID).
		Delete(&models.ComplianceAssessmentRequirementResidualRiskAnalysis{}).Error
}

// GetResidualRiskAnalysisByResponsibleUnitID retrieves all residual risk analysis for a responsible unit
func (repo *ComplianceAssessmentRequirementResidualRiskAnalysisRepository) GetResidualRiskAnalysisByResponsibleUnitID(responsibleUnitID uint) ([]*models.ComplianceAssessmentRequirementResidualRiskAnalysis, error) {
	var riskAnalysisList []*models.ComplianceAssessmentRequirementResidualRiskAnalysis

	err := repo.db.Where("compliance_assessment_requirement_responsible_unit_id = ?", responsibleUnitID).
		Find(&riskAnalysisList).Error

	if err != nil {
		return nil, err
	}

	return riskAnalysisList, nil
}

// GetComplianceAssessmentRatings retrieves all compliance assessment ratings
func (repo *ComplianceAssessmentRequirementResidualRiskAnalysisRepository) GetComplianceAssessmentRatings() ([]*models.ComplianceAssessmentRating, error) {
	var ratings []*models.ComplianceAssessmentRating

	err := repo.db.Where("status = ?", "active").
		Order("score_from ASC").
		Find(&ratings).Error

	if err != nil {
		return nil, err
	}

	return ratings, nil
}

// GetAllByResponsibleUnitID retrieves all residual controls for a responsible unit
func (repo *ComplianceAssessmentRequirementResidualRiskAnalysisRepository) GetAllByResponsibleUnitID(responsibleUnitID uint, params *complianceDTOs.ComplianceAssessmentResidualControlListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentResidualControlListResponseDTO, error) {
	var residualControls []complianceDTOs.ComplianceAssessmentResidualControlListDTO

	// Base query with joins
	query := repo.db.Table("compliance_assessment_residual_controls carc").
		Select(`
			carc.id,
			carc.compliance_assessment_impact_type_id,
			cait.title_en as impact_type_title_en,
			cait.title_ar as impact_type_title_ar,
			CAST(crira.score AS CHAR) as inherent_risk_score,
			carc.control_id,
			c.reference_number as control_reference_number,
			c.title_en as control_title_en,
			c.title_ar as control_title_ar,
			c.description as control_description,
			carc.link_control,
			carc.is_not_applicable
		`).
		Joins("LEFT JOIN compliance_assessment_impact_types cait ON carc.compliance_assessment_impact_type_id = cait.id").
		Joins("LEFT JOIN controls c ON carc.control_id = c.id").
		Joins("LEFT JOIN compliance_requirement_inherent_risk_analysis crira ON carc.compliance_assessment_impact_type_id = crira.compliance_assessment_impact_type_id").
		Where("carc.compliance_assessment_requirement_responsible_unit_id = ?", responsibleUnitID)

	// Apply keyword search
	if params.Keyword != "" {
		searchTerm := "%" + params.Keyword + "%"
		query = query.Where("(cait.title_en LIKE ? OR cait.title_ar LIKE ? OR c.title_en LIKE ? OR c.title_ar LIKE ? OR c.reference_number LIKE ?)",
			searchTerm, searchTerm, searchTerm, searchTerm, searchTerm)
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "link_control":
			query = query.Where("carc.link_control = ?", true)
		case "not_applicable":
			query = query.Where("carc.is_not_applicable = ?", true)
		}
	}

	// Apply sorting
	if params.SortBy != "" && params.Sort != "" {
		orderClause := params.SortBy + " " + params.Sort
		query = query.Order(orderClause)
	} else {
		query = query.Order("cait.title_en ASC")
	}

	// Execute query
	if err := query.Find(&residualControls).Error; err != nil {
		return nil, err
	}

	return &complianceDTOs.ComplianceAssessmentResidualControlListResponseDTO{
		Items: residualControls,
	}, nil
}

// Create creates residual control records
func (repo *ComplianceAssessmentRequirementResidualRiskAnalysisRepository) Create(residualControls []*models.ComplianceAssessmentResidualControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(residualControls) == 0 {
		return nil
	}

	return db.Create(&residualControls).Error
}

// DeleteByResponsibleUnitID deletes all residual controls for a responsible unit
func (repo *ComplianceAssessmentRequirementResidualRiskAnalysisRepository) DeleteByResponsibleUnitID(responsibleUnitID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("compliance_assessment_requirement_responsible_unit_id = ?", responsibleUnitID).
		Delete(&models.ComplianceAssessmentResidualControl{}).Error
}

// UpdateResidualControlByImpactType updates IsNotApplicable flag for residual control by impact type
func (repo *ComplianceAssessmentRequirementResidualRiskAnalysisRepository) UpdateResidualControlByImpactType(responsibleUnitID uint, impactTypeID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	updates := map[string]interface{}{
		"is_not_applicable": true,
	}

	return db.Model(&models.ComplianceAssessmentResidualControl{}).
		Where("compliance_assessment_requirement_responsible_unit_id = ? AND compliance_assessment_impact_type_id = ?", responsibleUnitID, impactTypeID).
		Updates(updates).Error
}

func (repo *ComplianceAssessmentRequirementResidualRiskAnalysisRepository) AssessmentUpdateStatus(complianceAssessment *models.ComplianceAssessment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(complianceAssessment).Error; err != nil {
		return err
	}
	return nil
}
