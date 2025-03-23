package repositories

import (
	assetDTOs "jethur/grc/asset/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetKriRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(assetID uint, params *assetDTOs.AssetKriListQueryParamsDTO) (*assetDTOs.AssetKriListResponseDTO, error)
	GetByID(id uint) (*models.AssetKri, error)
	GetFullDetailsByID(id uint) (*models.AssetKri, error)
	Create(assetKris []*models.AssetKri, tx ...utility.DatabaseTransaction) ([]uint, error)
	Update(assetKri *models.AssetKri, tx ...utility.DatabaseTransaction) error
	Delete(assetKriID *uint, tx ...utility.DatabaseTransaction) error
	DeleteAssetKRIs(kriIds []*uint, tx ...utility.DatabaseTransaction) error
	GetKrisByAssetID(assetID uint) ([]*models.AssetKri, error)
	GetKRIByID(id uint) (*models.KRIRegister, error)
}

type AssetKriRepository struct {
	db *gorm.DB
}

func NewAssetKriRepository(db *gorm.DB) *AssetKriRepository {
	return &AssetKriRepository{
		db: db,
	}
}

func (repo *AssetKriRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AssetKriRepository) GetAll(assetID uint, params *assetDTOs.AssetKriListQueryParamsDTO) (*assetDTOs.AssetKriListResponseDTO, error) {
	var assetKris []*assetDTOs.AssetKriListDTO
	var totalCount int64

	query := repo.db.Table("asset_kris").
		Select(`
			asset_kris.*,
			kri_register.reference_number,
			kri_register.title,
			kri_register.description,
			kri_register.calculation_type,
			kri_register.low_threshold,
			kri_register.medium_threshold,
			kri_register.high_threshold,
			kri_register.current_value,
			kri_register.last_reviewed_date,
			kri_register.next_review_date,
			kri_register.last_review_note,
			kri_categories.title_en AS kri_category_title_en,
			kri_categories.title_ar AS kri_category_title_ar,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.id AS responsible_user_id,
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
			created_user.first_name_en,
			created_user.first_name_ar,
			created_user.last_name_en,
			created_user.last_name_ar,
			created_user.initial,
			created_user.color,
			created_user.email,
			created_user.profile_image_path
		`).
		Joins(`
			LEFT JOIN kri_register ON kri_register.id = asset_kris.kri_register_id
			LEFT JOIN kri_categories ON kri_categories.id = kri_register.kri_category_id
			LEFT JOIN organization_units ON organization_units.id = kri_register.responsible_unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = kri_register.responsible_user_id
			LEFT JOIN review_frequencies ON review_frequencies.id = kri_register.review_frequency_id
			LEFT JOIN risk_levels ON risk_levels.id = kri_register.risk_level_id
			LEFT JOIN users AS created_user ON created_user.id = asset_kris.created_by_id
		`).
		Where("asset_kris.asset_register_id = ?", assetID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.LastReviewedDateFrom != "" {
		query = query.Where("kri_register.last_reviewed_date >= ?", params.LastReviewedDateFrom)
	}
	if params.LastReviewedDateTo != "" {
		query = query.Where("kri_register.last_reviewed_date <= ?", params.LastReviewedDateTo)
	}
	if params.NextReviewDateFrom != "" {
		query = query.Where("kri_register.next_review_date >= ?", params.NextReviewDateFrom)
	}
	if params.NextReviewDateTo != "" {
		query = query.Where("kri_register.next_review_date <= ?", params.NextReviewDateTo)
	}

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

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			kri_register.title LIKE ? OR
			kri_register.calculation_type LIKE ? OR
			kri_register.description LIKE ? OR
			kri_register.reference_number LIKE ? OR
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
			responsible_user.last_name_ar LIKE ?
		`, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "kri_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	if err := query.Offset(offset).Limit(params.PageSize).Find(&assetKris).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &assetDTOs.AssetKriListResponseDTO{
		Items: assetKris,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *AssetKriRepository) GetByID(id uint) (*models.AssetKri, error) {
	var kri *models.AssetKri
	if err := repo.db.First(&kri, id).Error; err != nil {
		return nil, err
	}
	return kri, nil
}

func (repo *AssetKriRepository) GetFullDetailsByID(id uint) (*models.AssetKri, error) {
	var kri *models.AssetKri
	err := repo.db.
		Preload("AssetRegister").
		Preload("KriRegister").
		Preload("KriRegister.KRICategory").
		Preload("KriRegister.ResponsibleUnit").
		Preload("KriRegister.ResponsibleUser").
		Preload("KriRegister.ReviewFrequency").
		Preload("KriRegister.RiskLevel").
		Preload("KriRegister.ValueUpdates").
		Preload("CreatedByUser").
		First(&kri, id).Error

	if err != nil {
		return nil, err
	}
	return kri, nil
}

func (repo *AssetKriRepository) Create(assetKris []*models.AssetKri, tx ...utility.DatabaseTransaction) ([]uint, error) {
	db := repo.db
	if len(tx) > 0 {
		db = tx[0].DB()
	}

	if err := db.Create(&assetKris).Error; err != nil {
		return nil, err
	}

	var ids []uint
	for _, kri := range assetKris {
		ids = append(ids, kri.ID)
	}
	return ids, nil
}

func (repo *AssetKriRepository) Update(assetKri *models.AssetKri, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 {
		db = tx[0].DB()
	}
	return db.Save(assetKri).Error
}

func (repo *AssetKriRepository) Delete(assetKriID *uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 {
		db = tx[0].DB()
	}
	return db.Where("id = ?", assetKriID).Delete(&models.AssetKri{}).Error
}

func (repo *AssetKriRepository) DeleteAssetKRIs(kriIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 {
		db = tx[0].DB()
	}
	return db.Where("id IN (?)", kriIds).Delete(&models.AssetKri{}).Error
}

func (repo *AssetKriRepository) GetKrisByAssetID(assetID uint) ([]*models.AssetKri, error) {
	var kris []*models.AssetKri
	if err := repo.db.Where("asset_register_id = ?", assetID).Find(&kris).Error; err != nil {
		return nil, err
	}
	return kris, nil
}

func (repo *AssetKriRepository) GetKRIByID(id uint) (*models.KRIRegister, error) {
	var kriRegister *models.KRIRegister

	result := repo.db.First(&kriRegister, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return kriRegister, nil

}
