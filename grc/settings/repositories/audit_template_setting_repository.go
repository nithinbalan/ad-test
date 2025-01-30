package repositories

import (
	"jethur/grc/models"
	settingsDTOs "jethur/grc/settings/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditTemplateSettingRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *settingsDTOs.AuditTemplateSettingListQueryParamsDTO) ([]*settingsDTOs.AuditTemplateSettingListDTO, error)
	GetByID(id uint) (*models.AuditTemplateSetting, error)
	GetFirst() (*models.AuditTemplateSetting, error)
	Update(auditInvestigation *models.AuditTemplateSetting, tx ...utility.DatabaseTransaction) error
}

type AuditTemplateSettingRepository struct {
	db *gorm.DB
}

func NewAuditTemplateSettingRepository(db *gorm.DB) *AuditTemplateSettingRepository {
	return &AuditTemplateSettingRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditTemplateSettingRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditTemplateSettingRepository) GetAll(params *settingsDTOs.AuditTemplateSettingListQueryParamsDTO) ([]*settingsDTOs.AuditTemplateSettingListDTO, error) {
	var auditTemplateSettings []*settingsDTOs.AuditTemplateSettingListDTO

	query := repo.db.Table("audit_template_settings").
		Select(`
			audit_template_settings.*,
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
			LEFT JOIN users AS created_user ON created_user.id = audit_template_settings.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = audit_template_settings.updated_by_id
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
			(audit_template_settings.logo_image_path LIKE ? OR
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
	err := query.Find(&auditTemplateSettings).Error
	if err != nil {
		return nil, err
	}

	return auditTemplateSettings, nil
}

func (repo *AuditTemplateSettingRepository) GetByID(id uint) (*models.AuditTemplateSetting, error) {
	var auditTemplateSetting *models.AuditTemplateSetting
	result := repo.db.First(&auditTemplateSetting, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditTemplateSetting, nil
}

func (repo *AuditTemplateSettingRepository) GetFirst() (*models.AuditTemplateSetting, error) {
	var auditTemplateSetting models.AuditTemplateSetting
	result := repo.db.First(&auditTemplateSetting) // No id required here

	if result.Error != nil {
		return nil, result.Error
	}

	return &auditTemplateSetting, nil
}

func (repo *AuditTemplateSettingRepository) Update(auditTemplateSetting *models.AuditTemplateSetting, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditTemplateSetting).Error; err != nil {
		return err
	}

	return nil
}
