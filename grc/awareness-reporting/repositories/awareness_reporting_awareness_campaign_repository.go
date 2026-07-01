package repositories

import (
	"fmt"
	awarenessDTOs "jethur/grc/awareness-reporting/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type AwarenessCampaignRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *awarenessDTOs.AwarenessCampaignListQueryParamsDTO) (*awarenessDTOs.AwarenessCampaignListResponseDTO, error)
	GetByID(id uint) (*models.AwarenessCampaign, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AwarenessCampaign, error)
	Create(campaign *models.AwarenessCampaign, tx ...utility.DatabaseTransaction) (uint, error)
	Update(campaign *models.AwarenessCampaign, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	// Detail
	CreateDetail(detail *models.AwarenessCampaignDetail, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateDetail(detail *models.AwarenessCampaignDetail, tx ...utility.DatabaseTransaction) error
	GetDetailByAwarenessCampaignID(campaignID uint, tx ...utility.DatabaseTransaction) (*models.AwarenessCampaignDetail, error)
	DeleteDetail(detailID uint, tx ...utility.DatabaseTransaction) error

	CreateCampaignChannels(channels []*models.AwarenessCampaignChannel, tx ...utility.DatabaseTransaction) error
	DeleteCampaignChannels(detailID uint, tx ...utility.DatabaseTransaction) error

	// Document
	CreateDocuments(docs []*models.AwarenessCampaignDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.AwarenessCampaignDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type AwarenessCampaignRepository struct {
	db *gorm.DB
}

func NewAwarenessCampaignRepository(db *gorm.DB) *AwarenessCampaignRepository {
	return &AwarenessCampaignRepository{
		db: db,
	}
}

func (repo *AwarenessCampaignRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AwarenessCampaignRepository) GetAll(params *awarenessDTOs.AwarenessCampaignListQueryParamsDTO) (*awarenessDTOs.AwarenessCampaignListResponseDTO, error) {
	var campaigns []*awarenessDTOs.AwarenessCampaignListDTO
	var totalCount int64

	query := repo.db.Table("awareness_campaigns").
		Select(`
			awareness_campaigns.*,
			awareness_campaign_types.title_en AS awareness_campaign_type_title_en,
			awareness_campaign_types.title_ar AS awareness_campaign_type_title_ar,
			awareness_priorities.title_en AS awareness_priority_title_en,
			awareness_priorities.title_ar AS awareness_priority_title_ar,
			awareness_campaign_statuses.title_en AS awareness_campaign_status_title_en,
			awareness_campaign_statuses.title_ar AS awareness_campaign_status_title_ar,
			awareness_campaign_statuses.color_code AS awareness_campaign_status_color_code,
			organization_units.title_en AS organization_unit_title_en,
			organization_units.title_ar AS organization_unit_title_ar,
			owner_user.id AS owner_id,
			owner_user.first_name_en AS owner_first_name_en,
			owner_user.first_name_ar AS owner_first_name_ar,
			owner_user.last_name_en AS owner_last_name_en,
			owner_user.last_name_ar AS owner_last_name_ar,
			owner_user.initial AS owner_initial,
			owner_user.color AS owner_color,
			owner_user.email AS owner_email,
			owner_user.profile_image_path AS owner_profile_image_path,
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
			LEFT JOIN awareness_campaign_types ON awareness_campaign_types.id = awareness_campaigns.awareness_campaign_type_id
			LEFT JOIN awareness_priorities ON awareness_priorities.id = awareness_campaigns.awareness_priority_id
			LEFT JOIN awareness_campaign_statuses ON awareness_campaign_statuses.id = awareness_campaigns.awareness_campaign_status_id
			LEFT JOIN organization_units ON organization_units.id = awareness_campaigns.organization_unit_id
			LEFT JOIN users AS owner_user ON owner_user.id = awareness_campaigns.owner_id
			LEFT JOIN users AS created_user ON created_user.id = awareness_campaigns.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = awareness_campaigns.updated_by_id
		`)

	// Apply dynamic filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply date filters
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("awareness_campaigns.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(awareness_campaigns.campaign_title LIKE ? OR
			awareness_campaigns.reference_number LIKE ? OR
			awareness_campaigns.objective LIKE ? OR
			awareness_campaigns.target_audience LIKE ? OR
			awareness_campaign_types.title_en LIKE ? OR
			awareness_campaign_types.title_ar LIKE ? OR
			awareness_priorities.title_en LIKE ? OR
			awareness_priorities.title_ar LIKE ? OR
			awareness_campaign_statuses.title_en LIKE ? OR
			awareness_campaign_statuses.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			owner_user.first_name_en LIKE ? OR
			owner_user.first_name_ar LIKE ? OR
			owner_user.last_name_en LIKE ? OR
			owner_user.last_name_ar LIKE ? OR
			owner_user.email LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&campaigns).Error; err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &awarenessDTOs.AwarenessCampaignListResponseDTO{
		Items: campaigns,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AwarenessCampaignRepository) GetByID(id uint) (*models.AwarenessCampaign, error) {
	var campaign *models.AwarenessCampaign

	result := repo.db.First(&campaign, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return campaign, nil
}

func (repo *AwarenessCampaignRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AwarenessCampaign, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var campaign *models.AwarenessCampaign

	result := db.
		Preload("AwarenessCampaignType").
		Preload("AwarenessPriority").
		Preload("AwarenessCampaignStatus").
		Preload("Owner").
		Preload("ResponsibleUnit").
		Preload("AwarenessCampaignDetails").
		Preload("AwarenessCampaignDetails.Documents").
		Preload("AwarenessCampaignDetails.AwarenessEvaluationRating").
		Preload("AwarenessCampaignDetails.RiskClassification").
		Preload("AwarenessCampaignDetails.Channels").
		Preload("AwarenessCampaignDetails.Channels.AwarenessChannel").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&campaign, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return campaign, nil
}

func (repo *AwarenessCampaignRepository) Create(campaign *models.AwarenessCampaign, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(campaign).Error; err != nil {
		return 0, err
	}

	return campaign.ID, nil
}

func (repo *AwarenessCampaignRepository) Update(campaign *models.AwarenessCampaign, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(campaign).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AwarenessCampaignRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Delete(&models.AwarenessCampaign{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AwarenessCampaignRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.AwarenessCampaign{}).
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
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "AC", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *AwarenessCampaignRepository) CreateDetail(detail *models.AwarenessCampaignDetail, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(detail).Error; err != nil {
		return 0, err
	}

	return detail.ID, nil
}

func (repo *AwarenessCampaignRepository) UpdateDetail(detail *models.AwarenessCampaignDetail, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(detail).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AwarenessCampaignRepository) GetDetailByAwarenessCampaignID(campaignID uint, tx ...utility.DatabaseTransaction) (*models.AwarenessCampaignDetail, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var detail models.AwarenessCampaignDetail
	if err := db.Where("awareness_campaign_id = ?", campaignID).First(&detail).Error; err != nil {
		return nil, err
	}
	return &detail, nil
}

func (repo *AwarenessCampaignRepository) DeleteDetail(detailID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Delete(&models.AwarenessCampaignDetail{}, detailID).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AwarenessCampaignRepository) CreateCampaignChannels(channels []*models.AwarenessCampaignChannel, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&channels).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AwarenessCampaignRepository) DeleteCampaignChannels(detailID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("awareness_campaign_details_id = ?", detailID).Delete(&models.AwarenessCampaignChannel{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AwarenessCampaignRepository) GetDocumentByID(id uint) (*models.AwarenessCampaignDocument, error) {
	var document models.AwarenessCampaignDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *AwarenessCampaignRepository) CreateDocuments(docs []*models.AwarenessCampaignDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AwarenessCampaignRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("id IN (?)", docIds).Delete(&models.AwarenessCampaignDocument{}).Error; err != nil {
		return err
	}
	return nil
}
