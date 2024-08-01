package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type KRIRegisterRepositoryInterface interface {
	GetAll(params *riskDTOs.KRIRegisterListQueryParamsDTO) (*riskDTOs.KRIRegisterListResponseDTO, error)
	GetByID(id uint) (*models.KRIRegister, error)
	GetFullDetailsByID(id uint) (*models.KRIRegister, error)
	Create(kriRegister *models.KRIRegister) (uint, error)
	Update(kriRegister *models.KRIRegister, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error
	ExistsByTitle(Title string, KRICategoryID uint) (bool, error)
	ExistsByTitleExcludingID(id uint, Title string, KRICategoryID uint) (bool, error)
	GenerateReferenceNumber() (string, error)
}

type KRIRegisterRepository struct {
	db *gorm.DB
}

func NewKRIRegisterRepository(db *gorm.DB) *KRIRegisterRepository {
	return &KRIRegisterRepository{
		db: db,
	}
}

func (repo *KRIRegisterRepository) GetAll(params *riskDTOs.KRIRegisterListQueryParamsDTO) (*riskDTOs.KRIRegisterListResponseDTO, error) {
	var kriRegisters []*riskDTOs.KRIRegisterListDTO
	var totalCount int64

	query := repo.db.Table("kri_register").
		Select(`
			kri_register.*,
			CONVERT(VARCHAR, kri_register.last_reviewed_date, 23) AS last_reviewed_date,
			CONVERT(VARCHAR, kri_register.next_review_date, 23) AS next_review_date,
			kri_categories.title_en AS kri_category_title_en,
			kri_categories.title_ar AS kri_category_title_ar,
			organization_units.title_en AS owner_unit_title_en,
			organization_units.title_ar AS owner_unit_title_ar,
			owner.first_name_en AS owner_first_name_en,
			owner.first_name_ar AS owner_first_name_ar,
			owner.last_name_en AS owner_last_name_en,
			owner.last_name_ar AS owner_last_name_ar,
			owner.initial AS owner_initial,
			owner.color AS owner_color,
			owner.email AS owner_email,
			owner.profile_image_path AS owner_profile_image_path,
			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,
			risk_levels.title_en AS risk_level_title_en,
			risk_levels.title_ar AS risk_level_title_ar,
			risk_levels.color_code AS risk_level_color_code,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
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
			LEFT JOIN kri_categories ON kri_categories.id = kri_register.kri_category_id
			LEFT JOIN organization_units ON organization_units.id = kri_register.owner_unit_id
			LEFT JOIN users AS owner ON owner.id = kri_register.owner_id
			LEFT JOIN review_frequencies ON review_frequencies.id = kri_register.review_frequency_id
			LEFT JOIN risk_levels ON risk_levels.id = kri_register.risk_level_id
			LEFT JOIN users AS created_user ON created_user.id = kri_register.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = kri_register.updated_by_id
		`)

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
			owner.first_name_en LIKE ? OR
			owner.first_name_ar LIKE ? OR
			owner.last_name_en LIKE ? OR
			owner.last_name_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?
		`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
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

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&kriRegisters).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &riskDTOs.KRIRegisterListResponseDTO{
		Items: kriRegisters,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}


func (repo *KRIRegisterRepository) GetByID(id uint) (*models.KRIRegister, error) {
	var kriRegister *models.KRIRegister

	result := repo.db.First(&kriRegister, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return kriRegister, nil

}

func (repo *KRIRegisterRepository) GetFullDetailsByID(id uint) (*models.KRIRegister, error) {
	var kriRegister *models.KRIRegister // Declare the variable inside the function

	result := repo.db.
		Preload("KRICategory").
		Preload("OwnerUnit").
		Preload("Owner").
		Preload("ReviewFrequency").
		Preload("RiskLevel").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&kriRegister, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return kriRegister, nil
}

func (repo *KRIRegisterRepository) Create(kriRegister *models.KRIRegister) (uint, error) {
	if err := repo.db.Create(kriRegister).Error; err != nil {
		return 0, err
	}
	return kriRegister.ID, nil
}

func (repo *KRIRegisterRepository) Update(kriRegister *models.KRIRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(kriRegister).Error; err != nil {
		return err
	}
	return nil
}

func (repo *KRIRegisterRepository) Delete(id uint) error {

	var kriRegister *models.KRIRegister
	if err := repo.db.First(&kriRegister, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.KRIRegister{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *KRIRegisterRepository) ExistsByTitle(Title string, KRICategoryID uint) (bool, error) {
	var count int64

	// Check for the existence of Title or QuestionAr in a single query
	err := repo.db.Model(&models.KRIRegister{}).
		Where("title = ? AND kri_category_id = ? ", Title, KRICategoryID).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil

}

func (repo *KRIRegisterRepository) ExistsByTitleExcludingID(id uint, Title string, KRICategoryID uint) (bool, error) {
	var count int64
	err := repo.db.Model(&models.KRIRegister{}).
		Where("id <> ? AND title = ? AND kri_category_id = ?", id, Title, KRICategoryID).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *KRIRegisterRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.KRIRegister{}).
		Select("reference_number").
		Order("reference_number desc").
		Limit(1).
		Pluck("reference_number", &lastReferenceNumber).
		Error

	if err != nil {
		return "", err
	}

	lastDigits := "0000"
	if lastReferenceNumber != "" && len(lastReferenceNumber) >= 4 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-4:]
	}

	currentYear := time.Now().Year()
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "KRI", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}
