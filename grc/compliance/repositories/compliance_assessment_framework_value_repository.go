package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/compliance/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type ComplianceAssessmentFrameworkValueRepositoryInterface interface {
	GetAll(complianceAssessmentFrameworkID uint, params *dtos.ComplianceAssessmentFrameworkValueListQueryParamsDTO) (*dtos.ComplianceAssessmentFrameworkValueListResponseDTO, error)
	GetByID(id uint) (*models.ComplianceAssessmentFrameworkValue, error)
	GetFullDetailsByID(id uint) (*models.ComplianceAssessmentFrameworkValue, error)
	Create(value *models.ComplianceAssessmentFrameworkValue) (uint, error)
	Update(value *models.ComplianceAssessmentFrameworkValue) error
	Delete(id uint) error
	GetByFrameworkID(complianceAssessmentFrameworkID uint) ([]*models.ComplianceAssessmentFrameworkValue, error)
}

type ComplianceAssessmentFrameworkValueRepository struct {
	db *gorm.DB
}

func NewComplianceAssessmentFrameworkValueRepository(db *gorm.DB) *ComplianceAssessmentFrameworkValueRepository {
	return &ComplianceAssessmentFrameworkValueRepository{db: db}
}

func (repo *ComplianceAssessmentFrameworkValueRepository) GetAll(complianceAssessmentFrameworkID uint, params *dtos.ComplianceAssessmentFrameworkValueListQueryParamsDTO) (*dtos.ComplianceAssessmentFrameworkValueListResponseDTO, error) {
	var values []*dtos.ComplianceAssessmentFrameworkValueListDTO
	var totalCount int64

	query := repo.db.Table("compliance_assessment_framework_values").
		Select(`
			compliance_assessment_framework_values.id,
			compliance_assessment_framework_values.title_en,
			compliance_assessment_framework_values.title_ar,
			compliance_assessment_framework_values.description,
			compliance_assessment_framework_values.compliance_assessment_framework_id,
			frameworks.title_en as compliance_assessment_framework_title_en,
			frameworks.title_ar as compliance_assessment_framework_title_ar,
			compliance_assessment_framework_values.compliance_status_id,
			status.title_en as compliance_status_title_en,
			status.title_ar as compliance_status_title_ar,
			status.type as compliance_status_type,
			status.color_code as compliance_status_color_code,
			compliance_assessment_framework_values.value,
			compliance_assessment_framework_values.score_min,
			compliance_assessment_framework_values.score_max,
			compliance_assessment_framework_values.is_applicable,
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
		Joins(`LEFT JOIN compliance_assessment_frameworks AS frameworks ON compliance_assessment_framework_values.compliance_assessment_framework_id = frameworks.id`).
		Joins(`LEFT JOIN compliance_statuses AS status ON compliance_assessment_framework_values.compliance_status_id = status.id`).
		Joins(`LEFT JOIN users AS created_user ON compliance_assessment_framework_values.created_by_id = created_user.id`).
		Joins(`LEFT JOIN users AS updated_user ON compliance_assessment_framework_values.updated_by_id = updated_user.id`).
		Where("compliance_assessment_framework_values.compliance_assessment_framework_id = ?", complianceAssessmentFrameworkID)

	// Apply filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(compliance_assessment_framework_values.title_en LIKE ? OR
			compliance_assessment_framework_values.title_ar LIKE ? OR
			compliance_assessment_framework_values.description LIKE ? OR
			compliance_assessment_framework_values.value LIKE ?)
		`, keyword, keyword, keyword, keyword)
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Find(&values).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &dtos.ComplianceAssessmentFrameworkValueListResponseDTO{
		Items: values,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ComplianceAssessmentFrameworkValueRepository) GetByID(id uint) (*models.ComplianceAssessmentFrameworkValue, error) {
	var value models.ComplianceAssessmentFrameworkValue
	result := repo.db.First(&value, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &value, nil
}

func (repo *ComplianceAssessmentFrameworkValueRepository) GetFullDetailsByID(id uint) (*models.ComplianceAssessmentFrameworkValue, error) {
	var value models.ComplianceAssessmentFrameworkValue
	result := repo.db.
		Preload("ComplianceAssessmentFramework").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("ComplianceStatus").
		First(&value, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &value, nil
}

func (repo *ComplianceAssessmentFrameworkValueRepository) Create(value *models.ComplianceAssessmentFrameworkValue) (uint, error) {
	if err := repo.db.Create(value).Error; err != nil {
		return 0, err
	}
	return value.ID, nil
}

func (repo *ComplianceAssessmentFrameworkValueRepository) Update(value *models.ComplianceAssessmentFrameworkValue) error {
	if err := repo.db.Save(value).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceAssessmentFrameworkValueRepository) Delete(id uint) error {
	var value *models.ComplianceAssessmentFrameworkValue
	if err := repo.db.First(&value, id).Error; err != nil {
		return err
	}

	if err := repo.db.Delete(&models.ComplianceAssessmentFrameworkValue{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceAssessmentFrameworkValueRepository) GetByFrameworkID(complianceAssessmentFrameworkID uint) ([]*models.ComplianceAssessmentFrameworkValue, error) {
	var values []*models.ComplianceAssessmentFrameworkValue
	if err := repo.db.Where("compliance_assessment_framework_id = ?", complianceAssessmentFrameworkID).Find(&values).Error; err != nil {
		return nil, err
	}
	return values, nil
}
