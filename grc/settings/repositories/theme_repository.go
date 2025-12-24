package repositories

import (
	"jethur/grc/models"
	"jethur/grc/settings/dtos"

	"gorm.io/gorm"
)

type ThemeRepositoryInterface interface {
	GetAll(Keyword string, SortBy string, Sort string) ([]*dtos.ThemeListDTO, error)
	GetFullDetailsByID(id uint) (*models.Theme, error)
	GetByID(id uint) (*models.Theme, error)
	Create(theme *models.Theme) (*models.Theme, error)
	Update(theme *models.Theme) error
	Delete(id uint) error
}

type ThemeRepository struct {
	db *gorm.DB
}

func NewThemeRepository(db *gorm.DB) *ThemeRepository {
	return &ThemeRepository{db: db}
}

func (repo *ThemeRepository) GetAll(Keyword string, SortBy string, Sort string) ([]*dtos.ThemeListDTO, error) {
	var themes = []*dtos.ThemeListDTO{}

	query := repo.db.
		Table("themes").
		Select(`
			themes.*,
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
			LEFT JOIN users AS created_user ON created_user.id = themes.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = themes.updated_by_id
		`)

	if Keyword != "" {
		keyword := "%" + Keyword + "%"
		query = query.Where(`
			(
				themes.title_en LIKE ? OR
				themes.title_ar LIKE ? OR
				themes.description_en LIKE ? OR
			    themes.description_ar LIKE ? OR
				created_user.first_name_en LIKE ? OR
				created_user.first_name_ar LIKE ? OR
				created_user.last_name_en LIKE ? OR
				created_user.last_name_ar LIKE ?
			)
		`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
		)
	}
	if SortBy == "" {
		SortBy = "created_at"
	}
	if Sort != "asc" && Sort != "desc" {
		Sort = "asc"
	}

	query = query.Order(SortBy + " " + Sort)

	err := query.Find(&themes).Error
	if err != nil {
		return nil, err
	}

	return themes, nil
}

func (repo *ThemeRepository) GetFullDetailsByID(id uint) (*models.Theme, error) {
	var theme models.Theme

	result := repo.db.
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&theme, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &theme, nil
}

func (repo *ThemeRepository) GetByID(id uint) (*models.Theme, error) {
	var theme models.Theme

	result := repo.db.First(&theme, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &theme, nil
}

func (repo *ThemeRepository) Create(theme *models.Theme) (*models.Theme, error) {
	err := repo.db.Create(theme).Error
	if err != nil {
		return nil, err
	}

	return theme, nil
}

func (repo *ThemeRepository) Update(theme *models.Theme) error {
	if err := repo.db.Save(theme).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ThemeRepository) Delete(id uint) error {
	return repo.db.Delete(&models.Theme{}, id).Error
}
