package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/identity-access-review/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type IdentityReviewCampaignRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.IdentityReviewCampaignListQueryParamsDTO) (*dtos.IdentityReviewCampaignListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaign, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaign, error)
	Create(campaign *models.IdentityReviewCampaign, tx ...utility.DatabaseTransaction) (uint, error)
	Update(campaign *models.IdentityReviewCampaign, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)
	CreateSystems(systems []*models.IdentityReviewCampaignSystem, tx ...utility.DatabaseTransaction) error
	DeleteSystemsByCampaignID(campaignID uint, tx ...utility.DatabaseTransaction) error
	CreateDocuments(documents []*models.IdentityReviewCampaignDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaignDocument, error)
	DeleteDocuments(documentIDs []uint, tx ...utility.DatabaseTransaction) error
}

type IdentityReviewCampaignRepository struct {
	db *gorm.DB
}

func NewIdentityReviewCampaignRepository(db *gorm.DB) *IdentityReviewCampaignRepository {
	return &IdentityReviewCampaignRepository{
		db: db,
	}
}

func (repo *IdentityReviewCampaignRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IdentityReviewCampaignRepository) GetAll(params *dtos.IdentityReviewCampaignListQueryParamsDTO) (*dtos.IdentityReviewCampaignListResponseDTO, error) {
	var campaigns []*dtos.IdentityReviewCampaignListDTO
	var totalCount int64

	query := repo.db.Table("identity_review_campaigns").
		Select(`
			identity_review_campaigns.id,
			identity_review_campaigns.reference_number,
			identity_review_campaigns.campaign_name,
			identity_review_campaigns.description,
			CONVERT(VARCHAR, identity_review_campaigns.start_date, 23) AS start_date,
			CONVERT(VARCHAR, identity_review_campaigns.due_date, 23) AS due_date,
			identity_review_campaigns.reminder_days_offset,
			identity_review_campaigns.completion_percentage,
			review_types.id AS review_type_id,
			review_types.title_en AS review_type_title_en,
			review_types.title_ar AS review_type_title_ar,
			review_campaign_statuses.id AS campaign_status_id,
			review_campaign_statuses.title_en AS campaign_status_title_en,
			review_campaign_statuses.title_ar AS campaign_status_title_ar,
			review_campaign_statuses.color_code AS campaign_status_color_code,
			reviewer.id AS reviewer_id,
			reviewer.first_name_en AS reviewer_first_name_en,
			reviewer.first_name_ar AS reviewer_first_name_ar,
			reviewer.last_name_en AS reviewer_last_name_en,
			reviewer.last_name_ar AS reviewer_last_name_ar,
			reviewer.initial AS reviewer_initial,
			approver.id AS approver_id,
			approver.first_name_en AS approver_first_name_en,
			approver.first_name_ar AS approver_first_name_ar,
			approver.last_name_en AS approver_last_name_en,
			approver.last_name_ar AS approver_last_name_ar,
			approver.initial AS approver_initial,
			identity_review_campaigns.created_at,
			identity_review_campaigns.updated_at
		`).
		Joins("LEFT JOIN review_types ON review_types.id = identity_review_campaigns.review_type_id").
		Joins("LEFT JOIN review_campaign_statuses ON review_campaign_statuses.id = identity_review_campaigns.campaign_status_id").
		Joins("LEFT JOIN users AS reviewer ON reviewer.id = identity_review_campaigns.reviewer_id").
		Joins("LEFT JOIN users AS approver ON approver.id = identity_review_campaigns.approver_id")

	for column, values := range params.Filters {
		if len(values) == 0 {
			continue
		}

		switch column {
		case "review_type_id":
			query = query.Where("identity_review_campaigns.review_type_id IN (?)", values)
		case "campaign_status_id":
			query = query.Where("identity_review_campaigns.campaign_status_id IN (?)", values)
		case "reviewer_id":
			query = query.Where("identity_review_campaigns.reviewer_id IN (?)", values)
		case "approver_id":
			query = query.Where("identity_review_campaigns.approver_id IN (?)", values)
		case "system_register_id":
			query = query.Where(`
				EXISTS (
					SELECT 1
					FROM identity_review_campaign_systems
					WHERE identity_review_campaign_systems.identity_review_campaign_id = identity_review_campaigns.id
					AND identity_review_campaign_systems.system_register_id IN (?)
				)`, values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			identity_review_campaigns.reference_number LIKE ? OR
			identity_review_campaigns.campaign_name LIKE ? OR
			identity_review_campaigns.description LIKE ? OR
			review_types.title_en LIKE ? OR
			review_types.title_ar LIKE ? OR
			review_campaign_statuses.title_en LIKE ? OR
			review_campaign_statuses.title_ar LIKE ? OR
			reviewer.first_name_en LIKE ? OR
			reviewer.first_name_ar LIKE ? OR
			reviewer.last_name_en LIKE ? OR
			reviewer.last_name_ar LIKE ? OR
			approver.first_name_en LIKE ? OR
			approver.first_name_ar LIKE ? OR
			approver.last_name_en LIKE ? OR
			approver.last_name_ar LIKE ?`,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword,
		)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&campaigns).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &dtos.IdentityReviewCampaignListResponseDTO{
		Items: campaigns,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *IdentityReviewCampaignRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaign, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var campaign models.IdentityReviewCampaign
	if err := db.First(&campaign, id).Error; err != nil {
		return nil, err
	}

	return &campaign, nil
}

func (repo *IdentityReviewCampaignRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaign, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var campaign models.IdentityReviewCampaign
	result := db.
		Preload("ReviewType").
		Preload("CampaignStatus").
		Preload("Reviewer").
		Preload("Approver").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Systems").
		Preload("Systems.SystemRegister").
		Preload("Documents").
		First(&campaign, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &campaign, nil
}

func (repo *IdentityReviewCampaignRepository) Create(campaign *models.IdentityReviewCampaign, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(campaign).Error; err != nil {
		return 0, err
	}

	return campaign.ID, nil
}

func (repo *IdentityReviewCampaignRepository) Update(campaign *models.IdentityReviewCampaign, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(campaign).Error
}

func (repo *IdentityReviewCampaignRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Delete(&models.IdentityReviewCampaign{}, id).Error
}

func (repo *IdentityReviewCampaignRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.IdentityReviewCampaign{}).
		Select("reference_number").
		Order("reference_number DESC").
		Limit(1).
		Pluck("reference_number", &lastReferenceNumber).
		Error
	if err != nil {
		return "", err
	}

	lastDigits := "0000"
	if len(lastReferenceNumber) >= 4 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-4:]
	}

	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("IRC-%d-%04d", time.Now().Year(), lastReferenceID+1), nil
}

func (repo *IdentityReviewCampaignRepository) CreateSystems(systems []*models.IdentityReviewCampaignSystem, tx ...utility.DatabaseTransaction) error {
	if len(systems) == 0 {
		return nil
	}

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(&systems).Error
}

func (repo *IdentityReviewCampaignRepository) DeleteSystemsByCampaignID(campaignID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("identity_review_campaign_id = ?", campaignID).
		Delete(&models.IdentityReviewCampaignSystem{}).Error
}

func (repo *IdentityReviewCampaignRepository) CreateDocuments(documents []*models.IdentityReviewCampaignDocument, tx ...utility.DatabaseTransaction) error {
	if len(documents) == 0 {
		return nil
	}

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(&documents).Error
}

func (repo *IdentityReviewCampaignRepository) GetDocumentByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaignDocument, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var document models.IdentityReviewCampaignDocument
	if err := db.First(&document, id).Error; err != nil {
		return nil, err
	}

	return &document, nil
}

func (repo *IdentityReviewCampaignRepository) DeleteDocuments(documentIDs []uint, tx ...utility.DatabaseTransaction) error {
	if len(documentIDs) == 0 {
		return nil
	}

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id IN (?)", documentIDs).
		Delete(&models.IdentityReviewCampaignDocument{}).Error
}
