package repositories

import (
	"errors"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditFindingRootCauseAnalysisRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(findingID uint, params *auditDTOs.AuditFindingRootCauseAnalysisListQueryParamsDTO) ([]*auditDTOs.AuditFindingRootCauseAnalysisListDTO, error)
	GetByID(id uint) (*models.AuditFindingRootCauseAnalysis, error)
	GetByFindingID(findingID uint) ([]*models.AuditFindingRootCauseAnalysis, error)
	Create(auditFindingRootCauseAnalysis *models.AuditFindingRootCauseAnalysis, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditFindingRootCauseAnalysis *models.AuditFindingRootCauseAnalysis, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GetLatestWhyByFindingID(findingID uint) (uint, error)
	UpdateWhyValuesAfterDeletion(findingID uint, deletedWhy uint, tx ...utility.DatabaseTransaction) error
	GetCountByFindingID(findingID uint, tx ...utility.DatabaseTransaction) (int64, error)
}

type AuditFindingRootCauseAnalysisRepository struct {
	db *gorm.DB
}

func NewAuditFindingRootCauseAnalysisRepository(db *gorm.DB) *AuditFindingRootCauseAnalysisRepository {
	return &AuditFindingRootCauseAnalysisRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditFindingRootCauseAnalysisRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditFindingRootCauseAnalysisRepository) GetAll(findingID uint, params *auditDTOs.AuditFindingRootCauseAnalysisListQueryParamsDTO) ([]*auditDTOs.AuditFindingRootCauseAnalysisListDTO, error) {
	var auditFindingRootCauseAnalyses []*auditDTOs.AuditFindingRootCauseAnalysisListDTO

	query := repo.db.Table("audit_finding_root_cause_analyses").
		Select(`
            audit_finding_root_cause_analyses.*,
            created_user.id AS created_user_id,
        	created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar, 
			created_user.initial AS created_user_initial,
        	created_user.color AS created_user_color,
        	created_user.email AS created_user_email,
        	created_user.profile_image_path AS created_user_profile_image_path,  
            updated_user.id AS updated_user_id,
    		updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
        	updated_user.color AS updated_user_color,
        	updated_user.email AS updated_user_email,
        	updated_user.profile_image_path AS updated_user_profile_image_path          
        `).
		Joins(`
            LEFT JOIN users AS created_user ON created_user.id = audit_finding_root_cause_analyses.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = audit_finding_root_cause_analyses.updated_by_id
        `).
		Where(`audit_finding_root_cause_analyses.audit_finding_id`, findingID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
            (audit_finding_root_cause_analyses.why LIKE ? OR
            audit_finding_root_cause_analyses.question LIKE ? OR
            audit_finding_root_cause_analyses.answer LIKE ?
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
	if err := query.Find(&auditFindingRootCauseAnalyses).Error; err != nil {
		return nil, err
	}

	return auditFindingRootCauseAnalyses, nil
}

func (repo *AuditFindingRootCauseAnalysisRepository) GetByID(id uint) (*models.AuditFindingRootCauseAnalysis, error) {
	var auditFindingRootCauseAnalysis *models.AuditFindingRootCauseAnalysis

	result := repo.db.First(&auditFindingRootCauseAnalysis, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditFindingRootCauseAnalysis, nil

}

func (repo *AuditFindingRootCauseAnalysisRepository) GetByFindingID(findingID uint) ([]*models.AuditFindingRootCauseAnalysis, error) {
	var analyses []*models.AuditFindingRootCauseAnalysis
	if err := repo.db.Where("audit_finding_id = ?", findingID).Find(&analyses).Error; err != nil {
		return nil, err
	}
	return analyses, nil
}

func (repo *AuditFindingRootCauseAnalysisRepository) Create(auditFindingRootCauseAnalysis *models.AuditFindingRootCauseAnalysis, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditFindingRootCauseAnalysis).Error; err != nil {
		return 0, err
	}

	return auditFindingRootCauseAnalysis.ID, nil
}

func (repo *AuditFindingRootCauseAnalysisRepository) Update(auditFindingRootCauseAnalysis *models.AuditFindingRootCauseAnalysis, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditFindingRootCauseAnalysis).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditFindingRootCauseAnalysisRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.AuditFindingRootCauseAnalysis{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditFindingRootCauseAnalysisRepository) UpdateWhyValuesAfterDeletion(findingID uint, deletedWhy uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Decrement the "Why" values of subsequent entries after the deleted entry
	err := db.Model(&models.AuditFindingRootCauseAnalysis{}).
		Where("audit_finding_id = ? AND why > ?", findingID, deletedWhy).
		UpdateColumn("why", gorm.Expr("why - 1")).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *AuditFindingRootCauseAnalysisRepository) GetLatestWhyByFindingID(findingID uint) (uint, error) {
	var latestWhy uint

	err := repo.db.Model(&models.AuditFindingRootCauseAnalysis{}).Where("audit_finding_id = ?", findingID).Order("why DESC").Limit(1).Pluck("why", &latestWhy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no records found, return 0
			return 0, nil
		}
		return 0, err
	}
	return latestWhy, nil
}

func (repo *AuditFindingRootCauseAnalysisRepository) GetCountByFindingID(findingID uint, tx ...utility.DatabaseTransaction) (int64, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var count int64
	if err := db.Model(&models.AuditFindingRootCauseAnalysis{}).
		Where("audit_finding_id = ?", findingID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
