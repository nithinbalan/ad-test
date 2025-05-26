package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type CompetencyRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Competency, error)
	GetByID(id uint) (*models.Competency, error)
	Create(competency *models.Competency) (*models.Competency, error)
	Update(id uint, updatedCompetency *models.Competency) (*models.Competency, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// Implementation
type CompetencyRepository struct {
	db *gorm.DB
}

func NewCompetencyRepository(db *gorm.DB) *CompetencyRepository {
	return &CompetencyRepository{db: db}
}

func (repo *CompetencyRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Competency, error) {
	var competencies []*models.Competency
	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
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

	query = query.Order(sortBy + " " + sort)
	err := query.Find(&competencies).Error
	return competencies, err
}

func (repo *CompetencyRepository) GetByID(id uint) (*models.Competency, error) {
	var competency *models.Competency
	err := repo.db.First(&competency, id).Error
	if err != nil {
		return nil, err
	}
	return competency, nil
}

func (repo *CompetencyRepository) Create(competency *models.Competency) (*models.Competency, error) {
	err := repo.db.Create(&competency).Error
	return competency, err
}

func (repo *CompetencyRepository) Update(id uint, updatedCompetency *models.Competency) (*models.Competency, error) {
	var existing *models.Competency
	if err := repo.db.First(&existing, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&existing).Updates(updatedCompetency).Error; err != nil {
		return nil, err
	}

	return existing, nil
}

func (repo *CompetencyRepository) Delete(id uint) error {
	var competency *models.Competency
	if err := repo.db.First(&competency, id).Error; err != nil {
		return err
	}
	return repo.db.Delete(&models.Competency{}, id).Error
}

func (repo *CompetencyRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	err := repo.db.Model(&models.Competency{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error
	if err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		err := repo.db.Model(&models.Competency{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error
		if err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *CompetencyRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	err := repo.db.Model(&models.Competency{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error
	if err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		err := repo.db.Model(&models.Competency{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error
		if err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *CompetencyRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.Competency{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *CompetencyRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.Competency{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
