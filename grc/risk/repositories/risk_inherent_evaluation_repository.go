package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskInherentEvaluationRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetByID(id uint) (*models.RiskEvaluation, error)
	Create(riskInherentEvaluation *models.RiskEvaluation, tx ...utility.DatabaseTransaction) (uint, error)
	Update(riskInherentEvaluation *models.RiskEvaluation, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	CreateRiskImpactEvaluations(evaluations []*models.RiskImpactEvaluation, tx ...utility.DatabaseTransaction) error
	CreateRiskAppetiteAlignments(alignments []*models.RiskAppetiteAlignment, tx ...utility.DatabaseTransaction) error
	DeleteRiskImpactEvaluationsByRiskEvaluationID(riskEvaluationID uint, tx ...utility.DatabaseTransaction) error
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

func (repo *RiskInherentEvaluationRepository) GetByID(id uint) (*models.RiskEvaluation, error) {
	var riskInherentEvaluation *models.RiskEvaluation

	result := repo.db.
		Preload("InherentRiskMatrixLikelihood").
		Preload("InherentRiskMatrixImpact").
		Preload("InherentRiskMatrixRiskRating").
		Preload("RiskImpactEvaluations").
		Preload("RiskAppetitieAlignments").
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

func (repo *RiskInherentEvaluationRepository) CreateRiskImpactEvaluations(evaluations []*models.RiskImpactEvaluation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&evaluations).Error; err != nil {
		return err
	}

	return nil

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

func (repo *RiskInherentEvaluationRepository) DeleteRiskImpactEvaluationsByRiskEvaluationID(riskEvaluationID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("risk_evaluation_id = ?", riskEvaluationID).Delete(&models.RiskImpactEvaluation{}).Error; err != nil {
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
