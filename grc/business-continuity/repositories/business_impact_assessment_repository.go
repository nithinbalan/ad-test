package repositories

import (
	"fmt"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BusinessImpactAssessmentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *businessContinuityDTOs.BusinessImpactAssessmentListQueryParamsDTO) (*businessContinuityDTOs.BusinessImpactAssessmentListResponseDTO, error)
	Create(businessImpactAssessment *models.BusinessImpactAssessment, tx ...utility.DatabaseTransaction) error
	DeleteBusinessImpactAssessmentsByProcessID(processID uint, tx ...utility.DatabaseTransaction) error
	DeleteBusinessImpactAssessmentsByAssetRegisterID(AssetRegisterID uint, tx ...utility.DatabaseTransaction) error
	GetBusinessImpactAssessmentsByProcessID(processID uint) ([]*models.BusinessImpactAssessment, error)
	GetBusinessImpactAssessmentsByAssetRegisterID(AssetRegisterID uint) ([]*models.BusinessImpactAssessment, error)
}

type BusinessImpactAssessmentRepository struct {
	db *gorm.DB
}

func NewBusinessImpactAssessmentRepository(db *gorm.DB) *BusinessImpactAssessmentRepository {
	return &BusinessImpactAssessmentRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *BusinessImpactAssessmentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *BusinessImpactAssessmentRepository) GetAll(params *businessContinuityDTOs.BusinessImpactAssessmentListQueryParamsDTO) (*businessContinuityDTOs.BusinessImpactAssessmentListResponseDTO, error) {
	var businessImpactAssessments []*businessContinuityDTOs.BusinessImpactAssessmentListDTO
	var totalCount int64

	query := repo.db.Table("business_impact_assessments").
		Select(`
            business_impact_assessments.*,
			processes.title_en AS process_title_en,
			processes.title_ar AS process_title_ar,
			asset_register.title_en AS asset_register_title_en,
			asset_register.title_ar AS asset_register_title_ar,
			impact_categories.title_en AS impact_category_title_en,
			impact_categories.title_ar AS impact_category_title_ar,
			risk_levels.title_en AS risk_level_title_en,
			risk_levels.title_ar AS risk_level_title_ar,
			risk_levels.color_code AS risk_level_color_code,
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
		    LEFT JOIN processes ON processes.id = business_impact_assessments.process_id
		    LEFT JOIN asset_register ON asset_register.id = business_impact_assessments.asset_register_id
			LEFT JOIN impact_categories ON impact_categories.id = business_impact_assessments.impact_category_id
		    LEFT JOIN risk_levels ON risk_levels.id = business_impact_assessments.risk_level_id
			LEFT JOIN users AS created_user ON created_user.id = business_impact_assessments.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = business_impact_assessments.updated_by_id
        `)
	// Apply dynamic filters to the query
	fmt.Println(params.Filters)
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(processes.title_en LIKE ? OR
			processes.title_ar LIKE ? OR
			asset_register.title_en LIKE ? OR
			asset_register.title_ar LIKE ? OR
			impact_categories.title_en LIKE ? OR
			impact_categories.title_ar LIKE ? OR
			risk_levels.title_en LIKE ? OR
			risk_levels.title_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&businessImpactAssessments).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &businessContinuityDTOs.BusinessImpactAssessmentListResponseDTO{
		Items: businessImpactAssessments,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *BusinessImpactAssessmentRepository) Create(businessImpactAssessment *models.BusinessImpactAssessment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(businessImpactAssessment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BusinessImpactAssessmentRepository) DeleteBusinessImpactAssessmentsByProcessID(processID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("process_id = ?", processID).Delete(&models.BusinessImpactAssessment{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *BusinessImpactAssessmentRepository) DeleteBusinessImpactAssessmentsByAssetRegisterID(AssetRegisterID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("asset_register_id = ?", AssetRegisterID).Delete(&models.BusinessImpactAssessment{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *BusinessImpactAssessmentRepository) GetBusinessImpactAssessmentsByProcessID(processID uint) ([]*models.BusinessImpactAssessment, error) {
	var businessImpactAssessments []*models.BusinessImpactAssessment
	err := repo.db.Table("business_impact_assessments").
		Select("business_impact_assessments.*").
		Where("process_id = ?", processID).
		Order("id").
		Find(&businessImpactAssessments).Error

	if err != nil {
		return nil, err
	}
	return businessImpactAssessments, nil
}

func (repo *BusinessImpactAssessmentRepository) GetBusinessImpactAssessmentsByAssetRegisterID(AssetRegisterID uint) ([]*models.BusinessImpactAssessment, error) {
	var businessImpactAssessments []*models.BusinessImpactAssessment
	err := repo.db.Table("business_impact_assessments").
		Select("business_impact_assessments.*").
		Where("asset_register_id = ?", AssetRegisterID).
		Order("id").
		Find(&businessImpactAssessments).Error

	if err != nil {
		return nil, err
	}
	return businessImpactAssessments, nil
}
