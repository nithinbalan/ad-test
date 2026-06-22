package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AwarenessAssessmentThemeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessAssessmentTheme, error)
	GetByID(id uint) (*models.AwarenessAssessmentTheme, error)
	Create(theme *models.AwarenessAssessmentTheme) (*models.AwarenessAssessmentTheme, error)
	Update(id uint, theme *models.AwarenessAssessmentTheme) (*models.AwarenessAssessmentTheme, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn string, titleAr string) (bool, error)
}

type AwarenessAssessmentThemeRepository struct {
	db *gorm.DB
}

func NewAwarenessAssessmentThemeRepository(db *gorm.DB) *AwarenessAssessmentThemeRepository {
	return &AwarenessAssessmentThemeRepository{
		db: db,
	}
}

func (repo *AwarenessAssessmentThemeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessAssessmentTheme, error) {
	var themes []*models.AwarenessAssessmentTheme
	query := repo.db

	if keyword != "" {
		query = query.Where(
			"(title_en LIKE ? OR title_ar LIKE ?)",
			"%"+keyword+"%",
			"%"+keyword+"%",
		)
	}

	if sortBy == "" {
		sortBy = "title_en"
	}

	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	result := query.Order(sortBy + " " + sort).Find(&themes)

	return themes, result.Error
}

func (repo *AwarenessAssessmentThemeRepository) GetByID(id uint) (*models.AwarenessAssessmentTheme, error) {
	var theme models.AwarenessAssessmentTheme

	if err := repo.db.First(&theme, id).Error; err != nil {
		return nil, err
	}

	return &theme, nil
}

func (repo *AwarenessAssessmentThemeRepository) Create(theme *models.AwarenessAssessmentTheme) (*models.AwarenessAssessmentTheme, error) {
	if err := repo.db.Create(theme).Error; err != nil {
		return nil, err
	}

	return theme, nil
}

func (repo *AwarenessAssessmentThemeRepository) Update(id uint, updatedTheme *models.AwarenessAssessmentTheme) (*models.AwarenessAssessmentTheme, error) {
	var theme models.AwarenessAssessmentTheme

	if err := repo.db.First(&theme, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&theme).Updates(updatedTheme).Error; err != nil {
		return nil, err
	}

	return &theme, nil
}

func (repo *AwarenessAssessmentThemeRepository) Delete(id uint) error {
	var theme models.AwarenessAssessmentTheme

	if err := repo.db.First(&theme, id).Error; err != nil {
		return err
	}

	return repo.db.Delete(&theme).Error
}

func (repo *AwarenessAssessmentThemeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.AwarenessAssessmentTheme{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AwarenessAssessmentTheme{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *AwarenessAssessmentThemeRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.AwarenessAssessmentTheme{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AwarenessAssessmentTheme{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *AwarenessAssessmentThemeRepository) ExistsByTitle(titleEn string, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.AwarenessAssessmentTheme{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *AwarenessAssessmentThemeRepository) ExistsByTitleExcludingID(id uint, titleEn string, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.AwarenessAssessmentTheme{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}