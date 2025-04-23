package repositories

import (
	"jethur/grc/models"
	settingsDTOs "jethur/grc/settings/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ReportTemplateSettingRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *settingsDTOs.ReportTemplateSettingListQueryParamsDTO) ([]*settingsDTOs.ReportTemplateSettingListDTO, error)
	GetByID(id uint) (*models.ReportTemplateSetting, error)
	GetFirst() (*models.ReportTemplateSetting, error)
	Update(reportInvestigation *models.ReportTemplateSetting, tx ...utility.DatabaseTransaction) error
}

type ReportTemplateSettingRepository struct {
	db *gorm.DB
}

func NewReportTemplateSettingRepository(db *gorm.DB) *ReportTemplateSettingRepository {
	return &ReportTemplateSettingRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ReportTemplateSettingRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ReportTemplateSettingRepository) GetAll(params *settingsDTOs.ReportTemplateSettingListQueryParamsDTO) ([]*settingsDTOs.ReportTemplateSettingListDTO, error) {
	var reportTemplateSettings []*settingsDTOs.ReportTemplateSettingListDTO

	query := repo.db.Table("report_template_settings").
		Select(`
			report_template_settings.*,
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
			LEFT JOIN users AS created_user ON created_user.id = report_template_settings.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = report_template_settings.updated_by_id
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
			(report_template_settings.logo_image_path LIKE ? OR
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
	err := query.Find(&reportTemplateSettings).Error
	if err != nil {
		return nil, err
	}

	return reportTemplateSettings, nil
}

func (repo *ReportTemplateSettingRepository) GetByID(id uint) (*models.ReportTemplateSetting, error) {
	var reportTemplateSetting *models.ReportTemplateSetting
	result := repo.db.First(&reportTemplateSetting, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return reportTemplateSetting, nil
}

func (repo *ReportTemplateSettingRepository) GetFirst() (*models.ReportTemplateSetting, error) {
	var reportTemplateSetting models.ReportTemplateSetting
	result := repo.db.First(&reportTemplateSetting) // No id required here

	if result.Error != nil {
		return nil, result.Error
	}

	return &reportTemplateSetting, nil
}

func (repo *ReportTemplateSettingRepository) Update(reportTemplateSetting *models.ReportTemplateSetting, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(reportTemplateSetting).Error; err != nil {
		return err
	}

	return nil
}
