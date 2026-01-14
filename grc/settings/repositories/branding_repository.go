package repositories

import (
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type BrandingRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.BrandingListQueryParamsDTO) ([]*dtos.BrandingListDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Branding, error)
	GetFirst() (*models.Branding, error)
	Update(branding *models.Branding, tx ...utility.DatabaseTransaction) error
}

type BrandingRepository struct {
	db *gorm.DB
}

func NewBrandingRepository(db *gorm.DB) BrandingRepositoryInterface {
	return &BrandingRepository{
		db: db,
	}
}

func (repo *BrandingRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *BrandingRepository) GetAll(params *dtos.BrandingListQueryParamsDTO) ([]*dtos.BrandingListDTO, error) {

	var branding []*dtos.BrandingListDTO
	query := repo.db.Table("branding").
		Select(` 
		branding.id,
		branding.login_page_title_en,
		branding.login_page_title_ar,
		branding.login_page_subtitle_en,
		branding.login_page_subtitle_ar,
		branding.login_page_bg_image_path,
		branding.org_logo_image_path, 
		branding.favicon_image_path, 
		branding.loader_gif_path,
		branding.created_at, 
		branding.updated_at, 
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
		LEFT JOIN users AS created_user ON created_user.id = branding.created_by_id 
		LEFT JOIN users AS updated_user ON updated_user.id = branding.updated_by_id
		`)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(branding.login_page_title_en LIKE ? OR
			branding.login_page_title_ar LIKE ? OR
			branding.login_page_subtitle_en LIKE ? OR
			branding.login_page_subtitle_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword,
			keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "created_at"
	}

	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	err := query.Find(&branding).Error
	if err != nil {
		return nil, err
	}

	return branding, nil
}

func (repo *BrandingRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Branding, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var branding models.Branding

	if err := db.First(&branding, id).Error; err != nil {
		return nil, err
	}

	return &branding, nil
}

func (repo *BrandingRepository) GetFirst() (*models.Branding, error) {
	var branding models.Branding

	result := repo.db.First(&branding)
	if result.Error != nil {
		return nil, result.Error
	}

	return &branding, nil
}

func (repo *BrandingRepository) Update(branding *models.Branding, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(branding).Error; err != nil {
		return err
	}

	return nil
}
