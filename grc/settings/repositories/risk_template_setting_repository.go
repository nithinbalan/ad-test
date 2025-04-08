package repositories

import (
	"jethur/grc/models"
	settingsDTOs "jethur/grc/settings/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskTemplateSettingRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *settingsDTOs.RiskTemplateSettingListQueryParamsDTO) ([]*settingsDTOs.RiskTemplateSettingListDTO, error)
	GetByID(id uint) (*models.RiskTemplateSetting, error)
	GetFirst() (*models.RiskTemplateSetting, error)
	Update(riskInvestigation *models.RiskTemplateSetting, tx ...utility.DatabaseTransaction) error
}

type RiskTemplateSettingRepository struct {
	db *gorm.DB
}

func NewRiskTemplateSettingRepository(db *gorm.DB) *RiskTemplateSettingRepository {
	return &RiskTemplateSettingRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *RiskTemplateSettingRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskTemplateSettingRepository) GetAll(params *settingsDTOs.RiskTemplateSettingListQueryParamsDTO) ([]*settingsDTOs.RiskTemplateSettingListDTO, error) {
	var riskTemplateSettings []*settingsDTOs.RiskTemplateSettingListDTO

	query := repo.db.Table("risk_template_settings").
		Select(`
			risk_template_settings.*,
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
			LEFT JOIN users AS created_user ON created_user.id = risk_template_settings.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = risk_template_settings.updated_by_id
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(risk_template_settings.logo_image_path LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ?)`,

			keyword, keyword, keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "created_at" // Default sort by created_at if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Execute the query
	err := query.Find(&riskTemplateSettings).Error
	if err != nil {
		return nil, err
	}

	return riskTemplateSettings, nil
}

func (repo *RiskTemplateSettingRepository) GetByID(id uint) (*models.RiskTemplateSetting, error) {
	var riskTemplateSetting *models.RiskTemplateSetting
	result := repo.db.First(&riskTemplateSetting, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskTemplateSetting, nil
}

func (repo *RiskTemplateSettingRepository) GetFirst() (*models.RiskTemplateSetting, error) {
	var riskTemplateSetting models.RiskTemplateSetting
	result := repo.db.First(&riskTemplateSetting) // No id required here

	if result.Error != nil {
		return nil, result.Error
	}

	return &riskTemplateSetting, nil
}

func (repo *RiskTemplateSettingRepository) Update(riskTemplateSetting *models.RiskTemplateSetting, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(riskTemplateSetting).Error; err != nil {
		return err
	}

	return nil
}
