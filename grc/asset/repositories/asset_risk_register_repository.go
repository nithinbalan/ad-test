package repositories

import (
	"errors"
	"fmt"
	assetRiskDTOs "jethur/grc/asset/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"math"
	"strings"

	"strconv"
	"time"

	"gorm.io/gorm"
)

type AssetRiskRegisterRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error)
	GetAll(params *assetRiskDTOs.AssetRiskRegisterListQueryParamsDTO) (*assetRiskDTOs.AssetRiskRegisterListResponseDTO, error)

	GenerateReferenceNumber(tx ...utility.DatabaseTransaction) (string, error)
	Create(riskRegister *models.RiskRegister, tx ...utility.DatabaseTransaction) (uint, error)
	Update(riskRegister *models.RiskRegister, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	CreateRiskTimeline(timelines []models.RiskTimeline, tx ...utility.DatabaseTransaction) error
	DeleteRiskTimeline(riskID uint, tx ...utility.DatabaseTransaction) error

	CreateRiskAsset(asset *models.RiskAsset, tx ...utility.DatabaseTransaction) error
	UpdateRiskAssetByRiskId(asset *models.RiskAsset, tx ...utility.DatabaseTransaction) error
	DeleteRiskAssetByRiskId(riskID uint, tx ...utility.DatabaseTransaction) error

	CreateRiskVulnerability(vulnerability *models.RiskVulnerability, tx ...utility.DatabaseTransaction) error
	UpdateRiskVulnerabilityByRiskId(vulnerability *models.RiskVulnerability, tx ...utility.DatabaseTransaction) error
	DeleteRiskVulnerabilityByRiskId(riskID uint, tx ...utility.DatabaseTransaction) error

	// Masters
	GetRiskSourceByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskSource, error)
	GetRiskCategoryByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskCategory, error)
	GetAllTimelineTypes(tx ...utility.DatabaseTransaction) ([]models.RiskTimelineType, error)
	GetTimelineTypeByType(timelineType string, tx ...utility.DatabaseTransaction) (*models.RiskTimelineType, error)
	GetTimelineStatusByType(status string, tx ...utility.DatabaseTransaction) (*models.TimelineStatus, error)
	GetRiskTimelineByRiskAndType(riskID uint, timelineTypeID uint, tx ...utility.DatabaseTransaction) (*models.RiskTimeline, error)
	UpdateRiskTimeline(riskTimeline *models.RiskTimeline, tx ...utility.DatabaseTransaction) error
	GetRiskStatusByType(status string, tx ...utility.DatabaseTransaction) (*models.RiskStatus, error)
	GetImpactByScore(score int, tx ...utility.DatabaseTransaction) (*models.RiskMatrixImpact, error)
	GetLikelihoodByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskMatrixLikelihood, error)
	GetRiskRatingValueByValue(value int, tx ...utility.DatabaseTransaction) (*models.RiskMatrixRiskRatingValue, error)
	GetRiskRatingByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskMatrixRiskRating, error)
}

type AssetRiskRegisterRepository struct {
	db *gorm.DB
}

func NewAssetRiskRegisterRepository(db *gorm.DB) AssetRiskRegisterRepositoryInterface {
	return &AssetRiskRegisterRepository{
		db: db,
	}
}

func (repo *AssetRiskRegisterRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AssetRiskRegisterRepository) GetAll(params *assetRiskDTOs.AssetRiskRegisterListQueryParamsDTO) (*assetRiskDTOs.AssetRiskRegisterListResponseDTO, error) {
	var totalCount int64
	riskRegisters := []*assetRiskDTOs.AssetRiskRegisterListDTO{}

	query := repo.db.Table("risk_register").
		Select(`
			risk_register.id,
			risk_register.reference_number,
			risk_vulnerabilities.vulnerability_id,
			risk_assets.asset_register_id AS asset_id,

			risk_register.risk_code,
			risk_register.title,
			CONVERT(VARCHAR, risk_register.registration_date, 23) AS registration_date,
			risk_register.justification,
			risk_register.expected_cause,
			risk_register.confidentiality_score,
			risk_register.integrity_score,
			risk_register.availability_score,

			risk_categories.title_en AS risk_category_title_en,
			risk_categories.title_ar AS risk_category_title_ar,

			division.title_en AS division_title_en,
			division.title_ar AS division_title_ar,

			section.title_en AS section_title_en,
			section.title_ar AS section_title_ar,

			current_likelihood.title_en AS current_likelihood_title_en,
			current_likelihood.title_ar AS current_likelihood_title_ar,

			current_impact.title_en AS current_impact_title_en,
			current_impact.title_ar AS current_impact_title_ar,

			current_risk_rating.title_en AS current_risk_rating_title_en,
			current_risk_rating.title_ar AS current_risk_rating_title_ar,
			current_risk_rating.color_code AS current_risk_rating_color_code,

			risk_register.is_strategic,
			risk_register.is_operational,
			risk_register.is_internal,
			risk_register.is_external,

			CONVERT(VARCHAR, risk_register.created_at, 23) AS created_at,
			CONVERT(VARCHAR, risk_register.updated_at, 23) AS updated_at,

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
		Joins("LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id").
		Joins("LEFT JOIN organization_units AS division ON division.id = risk_register.division_id").
		Joins("LEFT JOIN organization_units AS section ON section.id = risk_register.section_id").
		Joins("LEFT JOIN risk_matrix_likelihoods AS current_likelihood ON current_likelihood.id = risk_register.current_likelihood_id").
		Joins("LEFT JOIN risk_matrix_impacts AS current_impact ON current_impact.id = risk_register.current_impact_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings AS current_risk_rating ON current_risk_rating.id = risk_register.current_risk_rating_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = risk_register.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = risk_register.updated_by_id").
		Joins("LEFT JOIN risk_assets ON risk_assets.risk_register_id = risk_register.id").
		Joins("LEFT JOIN risk_vulnerabilities ON risk_vulnerabilities.risk_register_id = risk_register.id").
		Joins("LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Where("risk_statuses.type NOT IN (?)", models.PendingRiskStatusTypes).
		Where("risk_register.is_cyber_risk = ?", true)

	// Search
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"

		query = query.Where(`(
			risk_register.risk_code LIKE ? OR
			risk_register.title LIKE ? OR
			risk_categories.title_en LIKE ? OR
			risk_categories.title_ar LIKE ? OR
			division.title_en LIKE ? OR
			division.title_ar LIKE ? OR
			section.title_en LIKE ? OR
			section.title_ar LIKE ?)`,

			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
		)
	}

	// Dynamic Filter
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Quick filter
	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			query = query.Where("risk_statuses.type = ? ", "open")
		case "identified":
			query = query.Where("risk_statuses.type = ? ", "identified")
		}
	}

	params.Sort = strings.ToLower(params.Sort)

	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}

	// Order
	query = query.Order(params.SortBy + " " + params.Sort)

	// Count
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Export
	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Scan(&riskRegisters).Error; err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(params.PageSize)))

	return &assetRiskDTOs.AssetRiskRegisterListResponseDTO{
		Items: riskRegisters,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *AssetRiskRegisterRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskRegister models.RiskRegister
	if err := db.Where("id = ?", id).First(&riskRegister).Error; err != nil {
		return nil, err
	}

	return &riskRegister, nil
}

func (repo *AssetRiskRegisterRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskRegister models.RiskRegister

	err := db.
		Preload("Vulnerabilities.Vulnerability").
		Preload("Assets.AssetRegister").
		Preload("Controls").
		Preload("Controls.Control").
		Preload("RiskCategory").
		Preload("RiskStatus").
		Preload("Section").
		Preload("Division").
		Preload("ResponsibleUser").
		Preload("CurrentLikelihood").
		Preload("CurrentImpact").
		Preload("CurrentRiskRating").
		Preload("CurrentRiskRatingValue").
		Preload("InherentLikelihood").
		Preload("InherentImpact").
		Preload("InherentRiskRating").
		Preload("InherentRiskRatingValue").
		Preload("ResidualLikelihood").
		Preload("ResidualImpact").
		Preload("ResidualRiskRating").
		Preload("ResidualRiskRatingValue").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Where("is_cyber_risk", true).
		First(&riskRegister, id).Error

	if err != nil {
		return nil, err
	}

	return &riskRegister, nil
}

func (repo *AssetRiskRegisterRepository) GenerateReferenceNumber(tx ...utility.DatabaseTransaction) (string, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var lastReferenceNumber string
	err := db.Model(&models.RiskRegister{}).
		Select("reference_number").
		Order("reference_number desc").
		Where("reference_number LIKE ?", "IR-%").
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
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "IR", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *AssetRiskRegisterRepository) Create(riskRegister *models.RiskRegister, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(riskRegister).Error; err != nil {
		return 0, err
	}

	return riskRegister.ID, nil
}

func (repo *AssetRiskRegisterRepository) Update(riskRegister *models.RiskRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(riskRegister).Error
}

func (repo *AssetRiskRegisterRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Delete(&models.RiskRegister{}, id).Error
}

func (repo *AssetRiskRegisterRepository) CreateRiskAsset(asset *models.RiskAsset, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(asset).Error
}

func (repo *AssetRiskRegisterRepository) UpdateRiskAssetByRiskId(asset *models.RiskAsset, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Model(&models.RiskAsset{}).
		Where("risk_register_id = ?", asset.RiskRegisterID).
		Updates(asset).Error
}

func (repo *AssetRiskRegisterRepository) DeleteRiskAssetByRiskId(riskID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.
		Where("risk_register_id = ?", riskID).
		Delete(&models.RiskAsset{}).Error
}

func (repo *AssetRiskRegisterRepository) CreateRiskVulnerability(vulnerability *models.RiskVulnerability, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(vulnerability).Error
}

func (repo *AssetRiskRegisterRepository) UpdateRiskVulnerabilityByRiskId(vulnerability *models.RiskVulnerability, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Model(&models.RiskVulnerability{}).
		Where("risk_register_id = ?", vulnerability.RiskRegisterID).
		Updates(vulnerability).Error
}

func (repo *AssetRiskRegisterRepository) DeleteRiskVulnerabilityByRiskId(riskID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.
		Where("risk_register_id = ?", riskID).
		Delete(&models.RiskVulnerability{}).Error
}

func (repo *AssetRiskRegisterRepository) CreateRiskTimeline(timelines []models.RiskTimeline, tx ...utility.DatabaseTransaction) error {
	if len(timelines) == 0 {
		return nil
	}

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(&timelines).Error
}

func (repo *AssetRiskRegisterRepository) DeleteRiskTimeline(riskID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.
		Where("risk_register_id = ?", riskID).
		Delete(&models.RiskTimeline{}).Error
}

func (repo *AssetRiskRegisterRepository) GetRiskSourceByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskSource, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var source *models.RiskSource
	result := db.First(&source, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return source, nil
}

func (repo *AssetRiskRegisterRepository) GetRiskCategoryByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskCategory, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var category *models.RiskCategory
	result := db.First(&category, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return category, nil
}

func (repo *AssetRiskRegisterRepository) GetAllTimelineTypes(tx ...utility.DatabaseTransaction) ([]models.RiskTimelineType, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var timelineTypes []models.RiskTimelineType
	if err := db.Find(&timelineTypes).Error; err != nil {
		return nil, err
	}
	return timelineTypes, nil
}

func (repo *AssetRiskRegisterRepository) GetTimelineStatusByType(status string, tx ...utility.DatabaseTransaction) (*models.TimelineStatus, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var timelineStatus *models.TimelineStatus
	result := db.Where("type = ?", status).First(&timelineStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Timeline Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return timelineStatus, nil
}

func (repo *AssetRiskRegisterRepository) GetRiskStatusByType(status string, tx ...utility.DatabaseTransaction) (*models.RiskStatus, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskStatus *models.RiskStatus
	result := db.Where("type = ?", status).First(&riskStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Risk Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return riskStatus, nil
}

func (repo *AssetRiskRegisterRepository) GetImpactByScore(score int, tx ...utility.DatabaseTransaction) (*models.RiskMatrixImpact, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var impact *models.RiskMatrixImpact

	riskMatrixType := "isms"
	result := db.
		Joins("JOIN risk_matrix ON risk_matrix.id = risk_matrix_impacts.risk_matrix_id").
		Where("risk_matrix.type = ? AND risk_matrix_impacts.score = ?", riskMatrixType, score).
		First(&impact)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, fmt.Sprintf("Risk Matrix Impact with score %d", score))
	} else if result.Error != nil {
		return nil, result.Error
	}

	return impact, nil
}

func (repo *AssetRiskRegisterRepository) GetLikelihoodByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskMatrixLikelihood, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var likelihood *models.RiskMatrixLikelihood

	riskMatrixType := "isms"
	result := db.
		Joins("JOIN risk_matrix ON risk_matrix.id = risk_matrix_likelihoods.risk_matrix_id").
		Where("risk_matrix.type = ? AND risk_matrix_likelihoods.id = ?", riskMatrixType, id).
		First(&likelihood)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, fmt.Sprintf("Risk Matrix Likelihood with ID %d", id))
	} else if result.Error != nil {
		return nil, result.Error
	}

	return likelihood, nil
}

func (repo *AssetRiskRegisterRepository) GetRiskRatingValueByValue(value int, tx ...utility.DatabaseTransaction) (*models.RiskMatrixRiskRatingValue, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskRatingValue *models.RiskMatrixRiskRatingValue

	riskMatrixType := "isms"
	result := db.
		Joins("JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_matrix_risk_rating_values.risk_matrix_risk_rating_id").
		Joins("JOIN risk_matrix ON risk_matrix.id = risk_matrix_risk_ratings.risk_matrix_id").
		Where("risk_matrix.type = ? AND risk_matrix_risk_rating_values.value = ?", riskMatrixType, value).
		First(&riskRatingValue)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, fmt.Sprintf("Risk Rating Value with value %d", value))
	} else if result.Error != nil {
		return nil, result.Error
	}

	return riskRatingValue, nil
}

func (repo *AssetRiskRegisterRepository) GetRiskRatingByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskMatrixRiskRating, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskRating *models.RiskMatrixRiskRating

	riskMatrixType := "isms"
	result := db.
		Joins("JOIN risk_matrix ON risk_matrix.id = risk_matrix_risk_ratings.risk_matrix_id").
		Where("risk_matrix.type = ? AND risk_matrix_risk_ratings.id = ?", riskMatrixType, id).
		First(&riskRating)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, fmt.Sprintf("Risk Matrix Risk Rating with ID %d", id))
	} else if result.Error != nil {
		return nil, result.Error
	}

	return riskRating, nil
}

func (repo *AssetRiskRegisterRepository) GetTimelineTypeByType(timelineType string, tx ...utility.DatabaseTransaction) (*models.RiskTimelineType, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskTimelineType *models.RiskTimelineType
	result := db.Where("type = ?", timelineType).First(&riskTimelineType)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Risk Timeline Type with type "+timelineType)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return riskTimelineType, nil
}

func (repo *AssetRiskRegisterRepository) GetRiskTimelineByRiskAndType(riskID uint, timelineTypeID uint, tx ...utility.DatabaseTransaction) (*models.RiskTimeline, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskTimeline models.RiskTimeline
	result := db.Where("risk_register_id = ? AND risk_timeline_type_id = ?", riskID, timelineTypeID).First(&riskTimeline)
	if result.Error != nil {
		return nil, result.Error
	}

	return &riskTimeline, nil
}

func (repo *AssetRiskRegisterRepository) UpdateRiskTimeline(riskTimeline *models.RiskTimeline, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(riskTimeline).Error
}
