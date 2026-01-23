package repositories

import (
	"errors"
	"jethur/grc/models"
	violationReportDTOs "jethur/grc/violation-report/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ViolationReportRootCauseAnalysisRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(violationReportID uint, params *violationReportDTOs.ViolationReportRootCauseAnalysisListQueryParamsDTO) ([]*violationReportDTOs.ViolationReportRootCauseAnalysisListDTO, error)
	GetByID(id uint) (*models.ViolationReportRootCauseAnalysis, error)
	GetByViolationReportID(violationReportID uint) ([]*models.ViolationReportRootCauseAnalysis, error)
	GetAllByViolationReportID(violationReportID uint, tx ...utility.DatabaseTransaction) ([]*models.ViolationReportRootCauseAnalysis, error)
	Create(violationReportRootCauseAnalysis *models.ViolationReportRootCauseAnalysis, tx ...utility.DatabaseTransaction) (uint, error)
	Update(violationReportRootCauseAnalysis *models.ViolationReportRootCauseAnalysis, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GetLatestWhyByViolationReportID(violationReportID uint) (uint, error)
	UpdateWhyValuesAfterDeletion(violationReportID uint, deletedWhy uint, tx ...utility.DatabaseTransaction) error
	GetCountByRiskID(violationReportID uint, tx ...utility.DatabaseTransaction) (int64, error)
}

type ViolationReportRootCauseAnalysisRepository struct {
	db *gorm.DB
}

func NewViolationReportRootCauseAnalysisRepository(db *gorm.DB) *ViolationReportRootCauseAnalysisRepository {
	return &ViolationReportRootCauseAnalysisRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ViolationReportRootCauseAnalysisRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ViolationReportRootCauseAnalysisRepository) GetAll(violationReportID uint, params *violationReportDTOs.ViolationReportRootCauseAnalysisListQueryParamsDTO) ([]*violationReportDTOs.ViolationReportRootCauseAnalysisListDTO, error) {
	var violationReportRootCauseAnalyses []*violationReportDTOs.ViolationReportRootCauseAnalysisListDTO

	query := repo.db.Table("violation_report_root_cause_analyses").
		Select(`
            violation_report_root_cause_analyses.*,
            violation_reports.title_en AS violation_report_title_en,
            violation_reports.title_ar AS violation_report_title_ar,
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
            LEFT JOIN violation_reports ON violation_reports.id = violation_report_root_cause_analyses.violation_report_id
            LEFT JOIN users AS created_user ON created_user.id = violation_report_root_cause_analyses.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = violation_report_root_cause_analyses.updated_by_id
        `).
		Where(`violation_report_root_cause_analyses.violation_report_id`, violationReportID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
            (violation_report_root_cause_analyses.why LIKE ? OR
            violation_report_root_cause_analyses.question LIKE ? OR
            violation_report_root_cause_analyses.answer LIKE ?
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
	if err := query.Find(&violationReportRootCauseAnalyses).Error; err != nil {
		return nil, err
	}

	return violationReportRootCauseAnalyses, nil
}

func (repo *ViolationReportRootCauseAnalysisRepository) GetByID(id uint) (*models.ViolationReportRootCauseAnalysis, error) {
	var violationReportRootCauseAnalysis *models.ViolationReportRootCauseAnalysis

	result := repo.db.First(&violationReportRootCauseAnalysis, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return violationReportRootCauseAnalysis, nil

}

func (repo *ViolationReportRootCauseAnalysisRepository) GetByViolationReportID(violationReportID uint) ([]*models.ViolationReportRootCauseAnalysis, error) {
	var analyses []*models.ViolationReportRootCauseAnalysis
	if err := repo.db.Where("violation_report_id = ?", violationReportID).Find(&analyses).Error; err != nil {
		return nil, err
	}
	return analyses, nil
}

func (repo *ViolationReportRootCauseAnalysisRepository) GetAllByViolationReportID(violationReportID uint, tx ...utility.DatabaseTransaction) ([]*models.ViolationReportRootCauseAnalysis, error) {

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var rootCauseAnalyses []*models.ViolationReportRootCauseAnalysis

	err := db.Model(&models.ViolationReportRootCauseAnalysis{}).
		Where("violation_report_id = ?", violationReportID).
		Order("why ASC").
		Find(&rootCauseAnalyses).Error

	if err != nil {
		return nil, err
	}

	return rootCauseAnalyses, nil
}

func (repo *ViolationReportRootCauseAnalysisRepository) Create(violationReportRootCauseAnalysis *models.ViolationReportRootCauseAnalysis, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(violationReportRootCauseAnalysis).Error; err != nil {
		return 0, err
	}

	return violationReportRootCauseAnalysis.ID, nil
}

func (repo *ViolationReportRootCauseAnalysisRepository) Update(violationReportRootCauseAnalysis *models.ViolationReportRootCauseAnalysis, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(violationReportRootCauseAnalysis).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ViolationReportRootCauseAnalysisRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.ViolationReportRootCauseAnalysis{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ViolationReportRootCauseAnalysisRepository) UpdateWhyValuesAfterDeletion(violationReportID uint, deletedWhy uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Decrement the "Why" values of subsequent entries after the deleted entry
	err := db.Model(&models.ViolationReportRootCauseAnalysis{}).
		Where("violation_report_id = ? AND why > ?", violationReportID, deletedWhy).
		UpdateColumn("why", gorm.Expr("why - 1")).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *ViolationReportRootCauseAnalysisRepository) GetLatestWhyByViolationReportID(violationReportID uint) (uint, error) {
	var latestWhy uint

	err := repo.db.Model(&models.ViolationReportRootCauseAnalysis{}).Where("violation_report_id = ?", violationReportID).Order("why DESC").Limit(1).Pluck("why", &latestWhy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no records found, return 0
			return 0, nil
		}
		return 0, err
	}
	return latestWhy, nil
}

func (repo *ViolationReportRootCauseAnalysisRepository) GetCountByRiskID(violationReportID uint, tx ...utility.DatabaseTransaction) (int64, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var count int64
	if err := db.Model(&models.ViolationReportRootCauseAnalysis{}).
		Where("violation_report_id = ?", violationReportID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
