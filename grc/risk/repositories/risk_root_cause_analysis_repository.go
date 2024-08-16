package repositories

import (
	"errors"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskRootCauseAnalysisRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(riskID uint, params *riskDTOs.RiskRootCauseAnalysisListQueryParamsDTO) ([]*riskDTOs.RiskRootCauseAnalysisListDTO, error)
	GetByID(id uint) (*models.RiskRootCauseAnalysis, error)
	Create(riskRootCauseAnalysis *models.RiskRootCauseAnalysis, tx ...utility.DatabaseTransaction) (uint, error)
	Update(riskRootCauseAnalysis *models.RiskRootCauseAnalysis) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GetLatestWhyByRiskID(riskID uint) (uint, error)
	UpdateWhyValuesAfterDeletion(riskID uint, deletedWhy uint, tx ...utility.DatabaseTransaction) error
	GetCountByRiskID(riskID uint, tx ...utility.DatabaseTransaction) (int64, error)
}

type RiskRootCauseAnalysisRepository struct {
	db *gorm.DB
}

func NewRiskRootCauseAnalysisRepository(db *gorm.DB) *RiskRootCauseAnalysisRepository {
	return &RiskRootCauseAnalysisRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *RiskRootCauseAnalysisRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskRootCauseAnalysisRepository) GetAll(riskID uint, params *riskDTOs.RiskRootCauseAnalysisListQueryParamsDTO) ([]*riskDTOs.RiskRootCauseAnalysisListDTO, error) {
	var riskRootCauseAnalyses []*riskDTOs.RiskRootCauseAnalysisListDTO

	query := repo.db.Table("risk_root_cause_analysis").
		Select(`
            risk_root_cause_analysis.*,
            risk_register.title AS risk_register_title,
            created_user.id AS created_user_id,
            created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,   
            updated_user.id AS updated_user_id,
    		updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar          
        `).
		Joins(`
            LEFT JOIN risk_register ON risk_register.id = risk_root_cause_analysis.risk_register_id
            LEFT JOIN users AS created_user ON created_user.id = risk_root_cause_analysis.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = risk_root_cause_analysis.updated_by_id
        `).
		Where(`risk_root_cause_analysis.risk_register_id`, riskID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
            (risk_root_cause_analysis.why LIKE ? OR
            risk_root_cause_analysis.question LIKE ? OR
            risk_root_cause_analysis.answer LIKE ?
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
	if err := query.Find(&riskRootCauseAnalyses).Error; err != nil {
		return nil, err
	}

	return riskRootCauseAnalyses, nil
}

func (repo *RiskRootCauseAnalysisRepository) GetByID(id uint) (*models.RiskRootCauseAnalysis, error) {
	var riskRootCauseAnalysis *models.RiskRootCauseAnalysis

	result := repo.db.First(&riskRootCauseAnalysis, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return riskRootCauseAnalysis, nil

}

func (repo *RiskRootCauseAnalysisRepository) Create(riskRootCauseAnalysis *models.RiskRootCauseAnalysis, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(riskRootCauseAnalysis).Error; err != nil {
		return 0, err
	}

	return riskRootCauseAnalysis.ID, nil
}

func (repo *RiskRootCauseAnalysisRepository) Update(riskRootCauseAnalysis *models.RiskRootCauseAnalysis) error {
	if err := repo.db.Save(riskRootCauseAnalysis).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRootCauseAnalysisRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.RiskRootCauseAnalysis{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRootCauseAnalysisRepository) GetLatestWhyByRiskID(riskID uint) (uint, error) {
	var latestWhy uint

	err := repo.db.Model(&models.RiskRootCauseAnalysis{}).Where("risk_register_id = ?", riskID).Order("why DESC").Limit(1).Pluck("why", &latestWhy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no records found, return 0
			return 0, nil
		}
		return 0, err
	}
	return latestWhy, nil
}

func (repo *RiskRootCauseAnalysisRepository) UpdateWhyValuesAfterDeletion(riskID uint, deletedWhy uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Decrement the "Why" values of subsequent entries after the deleted entry
	err := db.Model(&models.RiskRootCauseAnalysis{}).
		Where("risk_register_id = ? AND why > ?", riskID, deletedWhy).
		UpdateColumn("why", gorm.Expr("why - 1")).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *RiskRootCauseAnalysisRepository) GetCountByRiskID(riskID uint, tx ...utility.DatabaseTransaction) (int64, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var count int64
	if err := db.Model(&models.RiskRootCauseAnalysis{}).
		Where("risk_register_id = ?", riskID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
