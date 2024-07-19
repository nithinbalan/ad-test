package repositories

import (
	"errors"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentRootCauseAnalysisRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(incidentID uint, params *incidentDTOs.IncidentRootCauseAnalysisListQueryParamsDTO) ([]*incidentDTOs.IncidentRootCauseAnalysisListDTO, error)
	GetByID(id uint) (*models.IncidentRootCauseAnalysis, error)
	GetByIncidentID(incidentID uint) ([]*models.IncidentRootCauseAnalysis, error)
	Create(incidentRootCauseAnalysis *models.IncidentRootCauseAnalysis, tx ...utility.DatabaseTransaction) (uint, error)
	Update(incidentRootCauseAnalysis *models.IncidentRootCauseAnalysis, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GetLatestWhyByIncidentID(incidentID uint) (uint, error)
	UpdateWhyValuesAfterDeletion(incidentID uint, deletedWhy uint, tx ...utility.DatabaseTransaction) error
}

type IncidentRootCauseAnalysisRepository struct {
	db *gorm.DB
}

func NewIncidentRootCauseAnalysisRepository(db *gorm.DB) *IncidentRootCauseAnalysisRepository {
	return &IncidentRootCauseAnalysisRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *IncidentRootCauseAnalysisRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IncidentRootCauseAnalysisRepository) GetAll(incidentID uint, params *incidentDTOs.IncidentRootCauseAnalysisListQueryParamsDTO) ([]*incidentDTOs.IncidentRootCauseAnalysisListDTO, error) {
	var incidentRootCauseAnalyses []*incidentDTOs.IncidentRootCauseAnalysisListDTO

	query := repo.db.Table("incident_root_cause_analyses").
		Select(`
            incident_root_cause_analyses.*,
            incident_register.title AS incident_register_title,
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
            LEFT JOIN incident_register ON incident_register.id = incident_root_cause_analyses.incident_register_id
            LEFT JOIN users AS created_user ON created_user.id = incident_root_cause_analyses.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = incident_root_cause_analyses.updated_by_id
        `).
		Where(`incident_root_cause_analyses.incident_register_id`, incidentID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
            (incident_root_cause_analyses.why LIKE ? OR
            incident_root_cause_analyses.question LIKE ? OR
            incident_root_cause_analyses.answer LIKE ?
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
	if err := query.Find(&incidentRootCauseAnalyses).Error; err != nil {
		return nil, err
	}

	return incidentRootCauseAnalyses, nil
}

func (repo *IncidentRootCauseAnalysisRepository) GetByID(id uint) (*models.IncidentRootCauseAnalysis, error) {
	var incidentRootCauseAnalysis *models.IncidentRootCauseAnalysis

	result := repo.db.First(&incidentRootCauseAnalysis, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return incidentRootCauseAnalysis, nil

}

func (repo *IncidentRootCauseAnalysisRepository) GetByIncidentID(incidentID uint) ([]*models.IncidentRootCauseAnalysis, error) {
	var analyses []*models.IncidentRootCauseAnalysis
	if err := repo.db.Where("incident_register_id = ?", incidentID).Find(&analyses).Error; err != nil {
		return nil, err
	}
	return analyses, nil
}

func (repo *IncidentRootCauseAnalysisRepository) Create(incidentRootCauseAnalysis *models.IncidentRootCauseAnalysis, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(incidentRootCauseAnalysis).Error; err != nil {
		return 0, err
	}

	return incidentRootCauseAnalysis.ID, nil
}

func (repo *IncidentRootCauseAnalysisRepository) Update(incidentRootCauseAnalysis *models.IncidentRootCauseAnalysis, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(incidentRootCauseAnalysis).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentRootCauseAnalysisRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.IncidentRootCauseAnalysis{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentRootCauseAnalysisRepository) UpdateWhyValuesAfterDeletion(incidentID uint, deletedWhy uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Decrement the "Why" values of subsequent entries after the deleted entry
	err := db.Model(&models.IncidentRootCauseAnalysis{}).
		Where("incident_register_id = ? AND why > ?", incidentID, deletedWhy).
		UpdateColumn("why", gorm.Expr("why - 1")).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *IncidentRootCauseAnalysisRepository) GetLatestWhyByIncidentID(incidentID uint) (uint, error) {
	var latestWhy uint

	err := repo.db.Model(&models.IncidentRootCauseAnalysis{}).Where("incident_register_id = ?", incidentID).Order("why DESC").Limit(1).Pluck("why", &latestWhy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no records found, return 0
			return 0, nil
		}
		return 0, err
	}
	return latestWhy, nil
}
