package repositories

import (
	"errors"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskInherentEvaluationRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(riskID uint, params *riskDTOs.RiskInherentEvaluationListQueryParamsDTO) ([]*riskDTOs.RiskInherentEvaluationListDTO, error)
	GetByID(id uint) (*models.RiskInherentEvaluation, error)
	Create(riskInherentEvaluation *models.RiskInherentEvaluation) (uint, error)
	Update(riskInherentEvaluation *models.RiskInherentEvaluation) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GetLatestWhyByRiskID(riskID uint) (uint, error)
	UpdateWhyValuesAfterDeletion(riskID uint, deletedWhy uint, tx ...utility.DatabaseTransaction) error
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

func (repo *RiskInherentEvaluationRepository) GetAll(riskID uint, params *riskDTOs.RiskInherentEvaluationListQueryParamsDTO) ([]*riskDTOs.RiskInherentEvaluationListDTO, error) {
	var riskInherentEvaluations []*riskDTOs.RiskInherentEvaluationListDTO

	query := repo.db.Table("risk_inherent_evaluation").
		Select(`
            risk_inherent_evaluation.*,
            risk_register.title AS risk_register_title,
            created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,   
            updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar          
        `).
		Joins(`
            LEFT JOIN risk_register ON risk_register.id = risk_inherent_evaluation.risk_register_id
            LEFT JOIN users AS created_user ON created_user.id = risk_inherent_evaluation.created_by
            LEFT JOIN users AS updated_user ON updated_user.id = risk_inherent_evaluation.updated_by
        `).
		Where(`risk_inherent_evaluation.risk_register_id`, riskID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
            (risk_inherent_evaluation.why LIKE ? OR
            risk_inherent_evaluation.question LIKE ? OR
            risk_inherent_evaluation.answer LIKE ?
            )`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "why" // Default sort by why if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Execute the query
	if err := query.Find(&riskInherentEvaluations).Error; err != nil {
		return nil, err
	}

	return riskInherentEvaluations, nil
}

func (repo *RiskInherentEvaluationRepository) GetByID(id uint) (*models.RiskInherentEvaluation, error) {
	var riskInherentEvaluation *models.RiskInherentEvaluation

	result := repo.db.First(&riskInherentEvaluation, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return riskInherentEvaluation, nil

}

func (repo *RiskInherentEvaluationRepository) Create(riskInherentEvaluation *models.RiskInherentEvaluation) (uint, error) {

	if err := repo.db.Create(riskInherentEvaluation).Error; err != nil {
		return 0, err
	}

	return riskInherentEvaluation.ID, nil
}

func (repo *RiskInherentEvaluationRepository) Update(riskInherentEvaluation *models.RiskInherentEvaluation) error {
	if err := repo.db.Save(riskInherentEvaluation).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskInherentEvaluationRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.RiskInherentEvaluation{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskInherentEvaluationRepository) GetLatestWhyByRiskID(riskID uint) (uint, error) {
	var latestWhy uint

	err := repo.db.Model(&models.RiskInherentEvaluation{}).Where("risk_register_id = ?", riskID).Order("why DESC").Limit(1).Pluck("why", &latestWhy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no records found, return 0
			return 0, nil
		}
		return 0, err
	}
	return latestWhy, nil
}

func (repo *RiskInherentEvaluationRepository) UpdateWhyValuesAfterDeletion(riskID uint, deletedWhy uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Decrement the "Why" values of subsequent entries after the deleted entry
	err := db.Model(&models.RiskInherentEvaluation{}).
		Where("risk_register_id = ? AND why > ?", riskID, deletedWhy).
		UpdateColumn("why", gorm.Expr("why - 1")).Error
	if err != nil {
		return err
	}
	return nil
}
