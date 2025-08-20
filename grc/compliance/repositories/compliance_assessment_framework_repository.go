package repositories

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ComplianceAssessmentFrameworkRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *complianceDTOs.ComplianceAssessmentFrameworkListQueryParamsDTO) ([]*complianceDTOs.ComplianceAssessmentFrameworkListDTO, error)
	GetByID(id uint) (*models.ComplianceAssessmentFramework, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessmentFramework, error)
	Create(complianceAssessmentFramework *models.ComplianceAssessmentFramework, tx ...utility.DatabaseTransaction) (uint, error)
	Update(complianceAssessmentFramework *models.ComplianceAssessmentFramework, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GetByType(status string) (*models.ComplianceAssessmentFramework, error)
	GetByFrameworkValueID(id uint) (*models.ComplianceAssessmentFrameworkValue, error)
}

// implementation
type ComplianceAssessmentFrameworkRepository struct {
	db *gorm.DB
}

func (repo *ComplianceAssessmentFrameworkRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func NewComplianceAssessmentFrameworkRepository(db *gorm.DB) ComplianceAssessmentFrameworkRepositoryInterface {
	return &ComplianceAssessmentFrameworkRepository{
		db: db,
	}
}

func (repo *ComplianceAssessmentFrameworkRepository) GetAll(params *complianceDTOs.ComplianceAssessmentFrameworkListQueryParamsDTO) ([]*complianceDTOs.ComplianceAssessmentFrameworkListDTO, error) {
	var complianceAssessmentFrameworks = []*complianceDTOs.ComplianceAssessmentFrameworkListDTO{}

	query := repo.db.Table("compliance_assessment_frameworks").
		Select(`
            compliance_assessment_frameworks.*,
			compliance_assessment_framework_types.id as compliance_assessment_framework_type_id,
			compliance_assessment_framework_types.title_en AS compliance_assessment_framework_type_title_en,
			compliance_assessment_framework_types.title_ar AS compliance_assessment_framework_type_title_ar,
			compliance_assessment_framework_types.description_en AS compliance_assessment_framework_type_description_en,
			compliance_assessment_framework_types.description_ar AS compliance_assessment_framework_type_description_ar,
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
			LEFT JOIN compliance_assessment_framework_types ON compliance_assessment_framework_types.id = compliance_assessment_frameworks.compliance_assessment_framework_type_id
			LEFT JOIN users AS created_user ON created_user.id = compliance_assessment_frameworks.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = compliance_assessment_frameworks.updated_by_id
        `)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(compliance_assessment_frameworks.title_en LIKE ? OR
			compliance_assessment_frameworks.title_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Apply sorting
	if params.SortBy != "" {
		sortOrder := "ASC"
		if params.Sort == "desc" {
			sortOrder = "DESC"
		}
		query = query.Order(params.SortBy + " " + sortOrder)
	}

	if err := query.Scan(&complianceAssessmentFrameworks).Error; err != nil {
		return nil, err
	}

	return complianceAssessmentFrameworks, nil

}

func (repo *ComplianceAssessmentFrameworkRepository) GetByID(id uint) (*models.ComplianceAssessmentFramework, error) {
	var complianceAssessmentFramework models.ComplianceAssessmentFramework
	if err := repo.db.First(&complianceAssessmentFramework, id).Error; err != nil {
		return nil, err
	}
	return &complianceAssessmentFramework, nil
}

func (repo *ComplianceAssessmentFrameworkRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessmentFramework, error) {
	var complianceAssessmentFramework models.ComplianceAssessmentFramework

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Preload("AssessmentFrameworkType").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&complianceAssessmentFramework, id).Error; err != nil {
		return nil, err
	}
	return &complianceAssessmentFramework, nil
}

func (repo *ComplianceAssessmentFrameworkRepository) Create(complianceAssessmentFramework *models.ComplianceAssessmentFramework, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(complianceAssessmentFramework).Error; err != nil {
		return 0, err
	}
	return complianceAssessmentFramework.ID, nil
}

func (repo *ComplianceAssessmentFrameworkRepository) Update(complianceAssessmentFramework *models.ComplianceAssessmentFramework, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(complianceAssessmentFramework).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceAssessmentFrameworkRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	var complianceAssessmentFramework *models.ComplianceAssessmentFramework
	if err := repo.db.First(&complianceAssessmentFramework, id).Error; err != nil {
		return err
	}

	if err := repo.db.Delete(&complianceAssessmentFramework).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceAssessmentFrameworkRepository) GetByType(status string) (*models.ComplianceAssessmentFramework, error) {
	var complianceAssessmentFramework models.ComplianceAssessmentFramework
	if err := repo.db.Where("status = ?", status).First(&complianceAssessmentFramework).Error; err != nil {
		return nil, err
	}
	return &complianceAssessmentFramework, nil
}

func (repo *ComplianceAssessmentFrameworkRepository) GetByFrameworkValueID(id uint) (*models.ComplianceAssessmentFrameworkValue, error) {
	var value models.ComplianceAssessmentFrameworkValue
	if err := repo.db.First(&value, id).Error; err != nil {
		return nil, err
	}
	return &value, nil
}

func (repo *ComplianceAssessmentFrameworkRepository) CreateComplianceAssessmentFrameworkValue(values []*models.ComplianceAssessmentFrameworkValue, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(&values).Error
}

func (repo *ComplianceAssessmentFrameworkRepository) DeleteComplianceAssessmentFrameworkValue(frameworkID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("compliance_assessment_framework_id = ?", frameworkID).Delete(&models.ComplianceAssessmentFrameworkValue{}).Error
}
