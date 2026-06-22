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

type AwarenessPostCampaignSurveyRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *awarenessDTOs.AwarenessPostCampaignSurveyListQueryParamsDTO) (*awarenessDTOs.AwarenessPostCampaignSurveyListResponseDTO, error)
	GetByID(id uint) (*models.AwarenessPostCampaignSurvey, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AwarenessPostCampaignSurvey, error)
	Create(survey *models.AwarenessPostCampaignSurvey, tx ...utility.DatabaseTransaction) (uint, error)
	Update(survey *models.AwarenessPostCampaignSurvey, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	CreateSurveyThemes(themes []*models.AwarenessPostCampaignSurveyTheme, tx ...utility.DatabaseTransaction) error
	DeleteSurveyThemes(surveyID uint, tx ...utility.DatabaseTransaction) error
}

type AwarenessPostCampaignSurveyRepository struct {
	db *gorm.DB
}

func NewAwarenessPostCampaignSurveyRepository(db *gorm.DB) *AwarenessPostCampaignSurveyRepository {
	return &AwarenessPostCampaignSurveyRepository{
		db: db,
	}
}

func (repo *AwarenessPostCampaignSurveyRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AwarenessPostCampaignSurveyRepository) GetAll(params *awarenessDTOs.AwarenessPostCampaignSurveyListQueryParamsDTO) (*awarenessDTOs.AwarenessPostCampaignSurveyListResponseDTO, error) {
	var surveys []*awarenessDTOs.AwarenessPostCampaignSurveyListDTO
	var totalCount int64

	query := repo.db.Table("awareness_post_campaign_surveys").
		Select(`
			awareness_post_campaign_surveys.*,
			awareness_campaigns.campaign_title AS linked_campaign_title,
			awareness_maturity_levels.title_en AS maturity_level_title_en,
			awareness_maturity_levels.title_ar AS maturity_level_title_ar,
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
			LEFT JOIN awareness_campaigns ON awareness_campaigns.id = awareness_post_campaign_surveys.linked_campaign_id
			LEFT JOIN awareness_maturity_levels ON awareness_maturity_levels.id = awareness_post_campaign_surveys.maturity_level_id
			LEFT JOIN users AS created_user ON created_user.id = awareness_post_campaign_surveys.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = awareness_post_campaign_surveys.updated_by_id
		`)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("awareness_post_campaign_surveys.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(awareness_post_campaign_surveys.reference_number LIKE ? OR
			awareness_post_campaign_surveys.target_audience LIKE ? OR
			awareness_post_campaign_surveys.executing_party LIKE ? OR
			awareness_campaigns.campaign_title LIKE ? OR
			awareness_maturity_levels.title_en LIKE ? OR
			awareness_maturity_levels.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword)
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&surveys).Error; err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &awarenessDTOs.AwarenessPostCampaignSurveyListResponseDTO{
		Items: surveys,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AwarenessPostCampaignSurveyRepository) GetByID(id uint) (*models.AwarenessPostCampaignSurvey, error) {
	var survey *models.AwarenessPostCampaignSurvey

	result := repo.db.First(&survey, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return survey, nil
}

func (repo *AwarenessPostCampaignSurveyRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AwarenessPostCampaignSurvey, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var survey *models.AwarenessPostCampaignSurvey

	result := db.
		Preload("LinkedCampaign").
		Preload("LinkedCampaign.AwarenessCampaignType").
		Preload("LinkedCampaign.AwarenessCampaignStatus").
		Preload("LinkedCampaign.Owner").
		Preload("MaturityLevel").
		Preload("SurveyThemes.AwarenessAssessmentTheme").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&survey, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return survey, nil
}

func (repo *AwarenessPostCampaignSurveyRepository) Create(survey *models.AwarenessPostCampaignSurvey, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(survey).Error; err != nil {
		return 0, err
	}

	return survey.ID, nil
}

func (repo *AwarenessPostCampaignSurveyRepository) Update(survey *models.AwarenessPostCampaignSurvey, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(survey).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AwarenessPostCampaignSurveyRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Delete(&models.AwarenessPostCampaignSurvey{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AwarenessPostCampaignSurveyRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.AwarenessPostCampaignSurvey{}).
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
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "APCS", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *AwarenessPostCampaignSurveyRepository) CreateSurveyThemes(themes []*models.AwarenessPostCampaignSurveyTheme, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&themes).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AwarenessPostCampaignSurveyRepository) DeleteSurveyThemes(surveyID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("awareness_post_campaign_survey_id = ?", surveyID).Delete(&models.AwarenessPostCampaignSurveyTheme{}).Error; err != nil {
		return err
	}
	return nil
}