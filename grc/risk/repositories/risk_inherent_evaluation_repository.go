package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskInherentEvaluationRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetFullDetailsByID(riskID uint) (*models.RiskEvaluation, error)
	GetByID(id uint) (*models.RiskEvaluation, error)
	Create(riskInherentEvaluation *models.RiskEvaluation, tx ...utility.DatabaseTransaction) (uint, error)
	Update(riskInherentEvaluation *models.RiskEvaluation, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	CreateRiskImpactAnalysis(riskImpactAnalysis *models.RiskImpactAnalysis, tx ...utility.DatabaseTransaction) error
	CreateRiskImpactAnalysisGuideline(riskImpactAnalysisGuideline *models.RiskImpactAnalysisGuideline, tx ...utility.DatabaseTransaction) error
	CreateRiskAppetiteAlignments(alignments []*models.RiskAppetiteAlignment, tx ...utility.DatabaseTransaction) error
	DeleteRiskImpactAnalysesByRiskEvaluationID(riskEvaluationID uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskAppetiteAlignmentsByRiskEvaluationID(riskEvaluationID uint, tx ...utility.DatabaseTransaction) error
}

type RiskInherentEvaluationRepository struct {
	db *gorm.DB
}

func NewRiskInherentEvaluationRepository(db *gorm.DB) *RiskInherentEvaluationRepository {
	return &RiskInherentEvaluationRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *RiskInherentEvaluationRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskInherentEvaluationRepository) GetFullDetailsByID(riskID uint) (*models.RiskEvaluation, error) {
	var riskInherentEvaluation models.RiskEvaluation

	// Preload all necessary relationships
	result := repo.db.
		Preload("RiskMatrixLikelihood").
		Preload("RiskMatrixImpact").
		Preload("RiskMatrixRiskRating").
		Preload("RiskMatrixRiskRatingValue").
		Preload("RiskImpactAnalyses.RiskImpactCategory").
		Preload("RiskImpactAnalyses.RiskMatrixImpact").
		Preload("RiskImpactAnalyses.RiskImpactAnalysisGuidelines.RiskMatrixImpactGuideline").
		Preload("RiskAppetiteAlignments.RiskAppetiteStatement").
		Preload("RiskAppetiteAlignments.RiskAppetiteToleranceLevel").
		Where("risk_register_id = ? AND type = ?", riskID, "inherent").
		First(&riskInherentEvaluation)
	if result.Error != nil {
		return nil, result.Error
	}

	return &riskInherentEvaluation, nil
}

func (repo *RiskInherentEvaluationRepository) GetByID(id uint) (*models.RiskEvaluation, error) {
	var riskInherentEvaluation *models.RiskEvaluation

	result := repo.db.
		First(&riskInherentEvaluation, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return riskInherentEvaluation, nil

}

func (repo *RiskInherentEvaluationRepository) Create(riskInherentEvaluation *models.RiskEvaluation, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(riskInherentEvaluation).Error; err != nil {
		return 0, err
	}

	return riskInherentEvaluation.ID, nil
}

func (repo *RiskInherentEvaluationRepository) Update(riskInherentEvaluation *models.RiskEvaluation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(riskInherentEvaluation).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskInherentEvaluationRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.RiskEvaluation{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskInherentEvaluationRepository) CreateRiskImpactAnalysis(riskImpactAnalysis *models.RiskImpactAnalysis, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(riskImpactAnalysis).Error
}

func (repo *RiskInherentEvaluationRepository) CreateRiskImpactAnalysisGuideline(riskImpactAnalysisGuideline *models.RiskImpactAnalysisGuideline, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(riskImpactAnalysisGuideline).Error
}

func (repo *RiskInherentEvaluationRepository) CreateRiskAppetiteAlignments(alignments []*models.RiskAppetiteAlignment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&alignments).Error; err != nil {
		return err
	}

	return nil

}

func (repo *RiskInherentEvaluationRepository) DeleteRiskImpactAnalysesByRiskEvaluationID(riskEvaluationID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Find all RiskImpactAnalyses associated with the given RiskEvaluationID
	var riskImpactAnalyses []*models.RiskImpactAnalysis
	if err := db.Where("risk_evaluation_id = ?", riskEvaluationID).Find(&riskImpactAnalyses).Error; err != nil {
		return err
	}

	// Delete related RiskImpactAnalysisGuidelines for each RiskImpactAnalysis
	for _, riskImpactAnalysis := range riskImpactAnalyses {
		if err := repo.DeleteRiskImpactAnalysisGuidelinesByRiskImpactAnalysisID(riskImpactAnalysis.ID, tx...); err != nil {
			return err
		}
	}

	// Delete the RiskImpactAnalyses
	if err := db.Where("risk_evaluation_id = ?", riskEvaluationID).Delete(&models.RiskImpactAnalysis{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskInherentEvaluationRepository) DeleteRiskImpactAnalysisGuidelinesByRiskImpactAnalysisID(riskImpactAnalysisID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("risk_impact_analysis_id = ?", riskImpactAnalysisID).Delete(&models.RiskImpactAnalysisGuideline{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskInherentEvaluationRepository) DeleteRiskAppetiteAlignmentsByRiskEvaluationID(riskEvaluationID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("risk_evaluation_id = ?", riskEvaluationID).Delete(&models.RiskAppetiteAlignment{}).Error; err != nil {
		return err
	}
	return nil
}
