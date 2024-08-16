package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskKriRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(riskID uint, params *riskDTOs.RiskKriListQueryParamsDTO) (*riskDTOs.RiskKriListResponseDTO, error)
	GetByID(id uint) (*models.RiskKri, error)
	GetFullDetailsByID(id uint) (*models.RiskKri, error)
	Create(riskKris []*models.RiskKri, tx ...utility.DatabaseTransaction) ([]uint, error)
	Update(riskKri *models.RiskKri, tx ...utility.DatabaseTransaction) error
	Delete(riskKriID *uint, tx ...utility.DatabaseTransaction) error
	GetCountByRiskID(riskID uint, tx ...utility.DatabaseTransaction) (int64, error)

	DeleteRiskKRIs(kriIds []*uint, tx ...utility.DatabaseTransaction) error
	GetKrisByRiskID(riskID uint) ([]*models.RiskKri, error)
}

type RiskKriRepository struct {
	db *gorm.DB
}

func NewRiskKriRepository(db *gorm.DB) *RiskKriRepository {
	return &RiskKriRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *RiskKriRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskKriRepository) GetAll(riskID uint, params *riskDTOs.RiskKriListQueryParamsDTO) (*riskDTOs.RiskKriListResponseDTO, error) {
	var riskKris []*riskDTOs.RiskKriListDTO
	var totalCount int64

	query := repo.db.Table("risk_kris").
		Select(`
        risk_kris.*,
		kri_register.reference_number,
        kri_register.title,
		kri_register.description AS description,
		kri_register.calculation_type AS calculation_type,
		kri_register.low_threshold AS low_threshold,
		kri_register.medium_threshold AS medium_threshold,
		kri_register.high_threshold AS high_threshold,
		kri_register.current_value AS current_value,
		kri_register.last_reviewed_date AS last_reviewed_date,
		kri_register.next_review_date AS next_review_date,
		kri_register.last_review_note AS last_review_note,
		kri_categories.title_en AS kri_category_title_en,
		kri_categories.title_ar AS kri_category_title_ar,
		organization_units.title_en AS responsible_unit_title_en,
		organization_units.title_ar AS responsible_unit_title_ar,
		responsible_user.first_name_en AS responsible_user_first_name_en,
		responsible_user.first_name_ar AS responsible_user_first_name_ar,
		responsible_user.last_name_en AS responsible_user_last_name_en,
		responsible_user.last_name_ar AS responsible_user_last_name_ar,
		responsible_user.initial AS responsible_user_initial,
		responsible_user.color AS responsible_user_color,
		responsible_user.email AS responsible_user_email,
		responsible_user.profile_image_path AS responsible_user_profile_image_path,
		review_frequencies.title_en AS review_frequency_title_en,
		review_frequencies.title_ar AS review_frequency_title_ar,
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
		created_user.profile_image_path AS created_user_profile_image_path
    `).
		Joins(`
		LEFT JOIN kri_register ON kri_register.id = risk_kris.kri_register_id
		LEFT JOIN kri_categories ON kri_categories.id = kri_register.kri_category_id
		LEFT JOIN organization_units ON organization_units.id = kri_register.responsible_unit_id
		LEFT JOIN users AS responsible_user ON responsible_user.id = kri_register.responsible_user_id
		LEFT JOIN review_frequencies ON review_frequencies.id = kri_register.review_frequency_id
		LEFT JOIN risk_levels ON risk_levels.id = kri_register.risk_level_id
        LEFT JOIN users AS created_user ON created_user.id = risk_kris.created_by_id

    `).Where("risk_kris.risk_register_id = ?", riskID)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply LastReviewDate filters
	if params.LastReviewedDateFrom != "" {
		query = query.Where("kri_register.last_reviewed_date >= ?", params.LastReviewedDateFrom)
	}
	if params.LastReviewedDateTo != "" {
		query = query.Where("kri_register.last_reviewed_date <= ?", params.LastReviewedDateTo)
	}

	// Apply next_review_date filters
	if params.NextReviewDateFrom != "" {
		query = query.Where("kri_register.next_review_date >= ?", params.NextReviewDateFrom)
	}
	if params.NextReviewDateTo != "" {
		query = query.Where("kri_register.next_review_date <= ?", params.NextReviewDateTo)
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "high_risk_level":
			query = query.Where("risk_levels.type = ?", "high")
		case "medium_risk_level":
			query = query.Where("risk_levels.type = ?", "medium")
		case "low_risk_level":
			query = query.Where("risk_levels.type = ?", "low")
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
			kri_register.title LIKE ? OR
			kri_register.calculation_type LIKE ? OR
			kri_register.description LIKE ? OR
			kri_register.reference_number LIKE ? OR
			kri_register.low_threshold LIKE ? OR
			kri_register.medium_threshold LIKE ? OR
			kri_register.high_threshold LIKE ? OR
			kri_register.current_value LIKE ? OR
			kri_categories.title_en LIKE ? OR
			kri_categories.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			review_frequencies.title_en LIKE ? OR
			review_frequencies.title_ar LIKE ? OR
			risk_levels.title_en LIKE ? OR
			risk_levels.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ?
		`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "kri_register.title" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	err := query.Offset(offset).Limit(params.PageSize).Find(&riskKris).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &riskDTOs.RiskKriListResponseDTO{
		Items: riskKris,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *RiskKriRepository) GetByID(id uint) (*models.RiskKri, error) {
	var riskKri *models.RiskKri
	result := repo.db.First(&riskKri, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskKri, nil
}

func (repo *RiskKriRepository) GetFullDetailsByID(id uint) (*models.RiskKri, error) {
	var riskKri *models.RiskKri

	result := repo.db.
		Preload("RiskRegister").
		Preload("KriRegister").
		Preload("KriRegister.KRICategory").
		Preload("KriRegister.ResponsibleUnit").
		Preload("KriRegister.ResponsibleUser").
		Preload("KriRegister.ReviewFrequency").
		Preload("KriRegister.RiskLevel").
		Preload("KriRegister.ValueUpdates").
		First(&riskKri, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskKri, nil
}

func (repo *RiskKriRepository) Create(riskKris []*models.RiskKri, tx ...utility.DatabaseTransaction) ([]uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var ids []uint
	if err := db.Create(&riskKris).Error; err != nil {
		return nil, err
	}

	for _, riskKri := range riskKris {
		ids = append(ids, riskKri.ID)
	}

	return ids, nil

}

func (repo *RiskKriRepository) Update(riskKri *models.RiskKri, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(riskKri).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskKriRepository) Delete(riskKriID *uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id", riskKriID).Delete(&models.RiskKri{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskKriRepository) GetCountByRiskID(riskID uint, tx ...utility.DatabaseTransaction) (int64, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var count int64
	if err := db.Model(&models.RiskKri{}).
		Where("risk_register_id = ?", riskID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *RiskKriRepository) DeleteRiskKRIs(kriIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", kriIds).Delete(&models.RiskKri{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskKriRepository) GetKrisByRiskID(riskID uint) ([]*models.RiskKri, error) {
	var kris []*models.RiskKri
	if err := repo.db.Where("risk_register_id = ?", riskID).Find(&kris).Error; err != nil {
		return nil, err
	}
	return kris, nil
}
