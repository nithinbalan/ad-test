package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskMatrixImpactGuidelineRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *riskDTOs.RiskMatrixImpactGuidelineListQueryParamsDTO) (*riskDTOs.RiskMatrixImpactGuidelineListResponseDTO, error)
	GetByID(id uint) (*models.RiskMatrixImpactGuideline, error)
	GetFullDetailsByID(id uint) (*models.RiskMatrixImpactGuideline, error)
	CreateBatch(guidelines []*models.RiskMatrixImpactGuideline, tx ...utility.DatabaseTransaction) ([]uint, error)
	Update(riskMatrixImpactGuideline *models.RiskMatrixImpactGuideline) error
	Delete(id uint) error
}

type RiskMatrixImpactGuidelineRepository struct {
	db *gorm.DB
}

func NewRiskMatrixImpactGuidelineRepository(db *gorm.DB) *RiskMatrixImpactGuidelineRepository {
	return &RiskMatrixImpactGuidelineRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *RiskMatrixImpactGuidelineRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskMatrixImpactGuidelineRepository) GetAll(params *riskDTOs.RiskMatrixImpactGuidelineListQueryParamsDTO) (*riskDTOs.RiskMatrixImpactGuidelineListResponseDTO, error) {

	var riskMatrixImpactGuidelines []*riskDTOs.RiskMatrixImpactGuidelineListDTO
	var totalCount int64

	query := repo.db.Table("risk_matrix_impact_guidelines").
	Select(`
        risk_matrix_impact_guidelines.*,
        risk_impact_categories.title_en AS risk_impact_category_title_en,
        risk_impact_categories.title_ar AS risk_impact_category_title_ar,
        risk_matrix_impacts.title_en AS risk_matrix_impact_title_en,
        risk_matrix_impacts.title_ar AS risk_matrix_impact_title_ar,
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
		LEFT JOIN risk_impact_categories ON risk_impact_categories.id = risk_matrix_impact_guidelines.risk_impact_category_id
        LEFT JOIN risk_matrix_impacts ON risk_matrix_impacts.id = risk_matrix_impact_guidelines.risk_matrix_impact_id
        LEFT JOIN users AS created_user ON created_user.id = risk_matrix_impact_guidelines.created_by
        LEFT JOIN users AS updated_user ON updated_user.id = risk_matrix_impact_guidelines.updated_by
    `)



	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(risk_impact_categories.title_en LIKE ? OR
			risk_impact_categories.title_ar LIKE ? OR
			risk_matrix_impacts.title_en LIKE ? OR
			risk_matrix_impacts.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}


	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "risk_matrix_impact_guidelines.guideline" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the queryc
	err := query.Offset(offset).Limit(params.PageSize).Find(&riskMatrixImpactGuidelines).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &riskDTOs.RiskMatrixImpactGuidelineListResponseDTO{
		Items: riskMatrixImpactGuidelines,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *RiskMatrixImpactGuidelineRepository) GetByID(id uint) (*models.RiskMatrixImpactGuideline, error) {
	var riskMatrixImpactGuideline *models.RiskMatrixImpactGuideline

	result := repo.db.First(&riskMatrixImpactGuideline, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return riskMatrixImpactGuideline, nil

}

func (repo *RiskMatrixImpactGuidelineRepository) GetFullDetailsByID(id uint) (*models.RiskMatrixImpactGuideline, error) {
	var riskMatrixImpactGuideline *models.RiskMatrixImpactGuideline // Declare the variable inside the function

	result := repo.db.
		Preload("RiskImpactCategory").
		Preload("RiskMatrixImpact").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&riskMatrixImpactGuideline, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskMatrixImpactGuideline, nil
}

func (repo *RiskMatrixImpactGuidelineRepository) CreateBatch(guidelines []*models.RiskMatrixImpactGuideline, tx ...utility.DatabaseTransaction) ([]uint, error) {
    db := repo.db
    if len(tx) > 0 && tx[0] != nil {
        db = tx[0].DB()
    }
    
    var ids []uint
    if err := db.Create(&guidelines).Error; err != nil {
        return nil, err
    }
    
    for _, guideline := range guidelines {
        ids = append(ids, guideline.ID)
    }
    
    return ids, nil
}


func (repo *RiskMatrixImpactGuidelineRepository) Update(riskMatrixImpactGuideline *models.RiskMatrixImpactGuideline) error {

	if err := repo.db.Save(riskMatrixImpactGuideline).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskMatrixImpactGuidelineRepository) Delete(id uint) error {

	var riskMatrixImpactGuideline *models.RiskMatrixImpactGuideline
	if err := repo.db.First(&riskMatrixImpactGuideline, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.RiskMatrixImpactGuideline{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
