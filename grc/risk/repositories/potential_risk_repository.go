package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type PotentialRiskRepositoryInterface interface {
	GetAll(params *riskDTOs.PotentialRiskListQueryParamsDTO) (*riskDTOs.PotentialRiskListResponseDTO, error)
	GetByID(id uint) (*models.PotentialRisk, error)
	GetFullDetailsByID(id uint) (*models.PotentialRisk, error)
	Create(potentialRisk *models.PotentialRisk) (uint, error)
	Update(potentialRisk *models.PotentialRisk) error
	Delete(id uint) error
	ExistsByTitle(Title string, RiskCategoryID uint) (bool, error)
	ExistsByTitleExcludingID(id uint, Title string, RiskCategoryID uint) (bool, error)
}

type PotentialRiskRepository struct {
	db *gorm.DB
}

func NewPotentialRiskRepository(db *gorm.DB) *PotentialRiskRepository {
	return &PotentialRiskRepository{
		db: db,
	}
}

func (repo *PotentialRiskRepository) GetAll(params *riskDTOs.PotentialRiskListQueryParamsDTO) (*riskDTOs.PotentialRiskListResponseDTO, error) {

	var potentialRisks []*riskDTOs.PotentialRiskListDTO
	var totalCount int64

	query := repo.db.Table("potential_risks").
		Select(`
            potential_risks.*,
			risk_categories.title_en AS risk_category_title_en,
			risk_categories.title_ar AS risk_category_title_ar,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			potential_risk_statuses.title_en AS potential_risk_status_title_en,
			potential_risk_statuses.title_ar AS potential_risk_status_title_ar,
			potential_risk_statuses.color_code AS potential_risk_status_color_code,
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
		    LEFT JOIN risk_categories ON risk_categories.id = potential_risks.risk_category_id
		    LEFT JOIN severity_levels ON severity_levels.id = potential_risks.severity_level_id
		    LEFT JOIN potential_risk_statuses ON potential_risk_statuses.id = potential_risks.potential_risk_status_id
			LEFT JOIN users AS created_user ON created_user.id = potential_risks.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = potential_risks.updated_by_id
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "confirmed":
			query = query.Where("potential_risk_statuses.type = ?", "confirmed")
		case "dismissed":
			query = query.Where("potential_risk_statuses.type = ?", "dismissed")
		case "new":
			query = query.Where("potential_risk_statuses.type = ?", "new")

		}
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(potential_risks.title LIKE ? OR
			potential_risks.risk_source LIKE ? OR
			risk_categories.title_en LIKE ? OR
			risk_categories.title_ar LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
			potential_risk_statuses.title_en LIKE ? OR
			potential_risk_statuses.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "title" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the queryc
	err := query.Offset(offset).Limit(params.PageSize).Find(&potentialRisks).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &riskDTOs.PotentialRiskListResponseDTO{
		Items: potentialRisks,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *PotentialRiskRepository) GetByID(id uint) (*models.PotentialRisk, error) {
	var potentialRisk *models.PotentialRisk

	result := repo.db.First(&potentialRisk, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return potentialRisk, nil

}

func (repo *PotentialRiskRepository) GetFullDetailsByID(id uint) (*models.PotentialRisk, error) {
	var potentialRisk *models.PotentialRisk // Declare the variable inside the function

	result := repo.db.
		Preload("RiskCategory").
		Preload("PotentialRiskStatus").
		Preload("SeverityLevel").
		// Preload("CreatedByUser").
		// Preload("UpdatedByUser").
		First(&potentialRisk, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return potentialRisk, nil
}

func (repo *PotentialRiskRepository) Create(potentialRisk *models.PotentialRisk) (uint, error) {

	if err := repo.db.Create(potentialRisk).Error; err != nil {
		return 0, err
	}

	return potentialRisk.ID, nil
}

func (repo *PotentialRiskRepository) Update(potentialRisk *models.PotentialRisk) error {

	if err := repo.db.Save(potentialRisk).Error; err != nil {
		return err
	}

	return nil
}

func (repo *PotentialRiskRepository) Delete(id uint) error {

	var potentialRisk *models.PotentialRisk
	if err := repo.db.First(&potentialRisk, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.PotentialRisk{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *PotentialRiskRepository) ExistsByTitle(Title string, RiskCategoryID uint) (bool, error) {
	var count int64

	// Check for the existence of Title or QuestionAr in a single query
	err := repo.db.Model(&models.PotentialRisk{}).
		Where("title = ? AND risk_category_id = ?", Title, RiskCategoryID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *PotentialRiskRepository) ExistsByTitleExcludingID(id uint, Title string, RiskCategoryID uint) (bool, error) {
	var count int64
	err := repo.db.Model(&models.PotentialRisk{}).
		Where("id <> ? AND title = ? AND risk_category_id = ?", id, Title, RiskCategoryID).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
